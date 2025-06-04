package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Siravitt/kafka-mongo-tracker/config"
	"github.com/gin-gonic/gin"
)

const (
	gracefulShutdownDuration = 10 * time.Second
	serverReadHeaderTimeout  = 5 * time.Second
	serverReadTimeout        = 5 * time.Second
	serverWriteTimeout       = 10 * time.Second // request hangup after this durations
	handlerTimeout           = serverWriteTimeout - (time.Millisecond * 100)
)

var headers = []string{
	"Content-Type",
	"Content-Length",
	"Accept-Encoding",
	"X-CSRF-Token",
	"Authorization",
	"accept",
	"origin",
	"Cache-Control",
	"X-Requested-With",
	os.Getenv("REF_ID_HEADER_KEY"),
}

type HTTPServer struct {
	r             *gin.Engine
	config        config.Config
	routeHandlers []RouteHandler
	middlewares   []gin.HandlerFunc
}

type RouteHandler interface {
	RegisterRoutes(r *gin.Engine)
}

func NewRouter(cfg config.Config) *HTTPServer {
	r := gin.New()

	if config.IsLocalEnv() {
		r.Use(gin.Logger())
	}

	srv := &HTTPServer{
		r:      r,
		config: cfg,
	}

	return srv
}

func (s *HTTPServer) StartServer() {
	// readiness, metrics, liveness
	s.readiness()

	// middlewares
	s.r.Use(s.middlewares...)
	s.r.Use(
		securityHeaders,
		accessControl,
		handlerTimeoutMiddleware,
	)

	// routes
	for _, handler := range s.routeHandlers {
		handler.RegisterRoutes(s.r)
	}

	srv := &http.Server{
		Addr:              ":" + s.config.Server.Port,
		Handler:           s.r,
		ReadHeaderTimeout: serverReadHeaderTimeout,
		ReadTimeout:       serverReadTimeout,
		WriteTimeout:      serverWriteTimeout,
	}

	// running server
	go func() {
		slog.Info("HTTP server started on :" + s.config.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server ListenAndServe: " + err.Error())
		}
	}()

	shutdownServer(srv)
}

func (s *HTTPServer) AddHandlers(plugin ...RouteHandler) {
	s.routeHandlers = append(s.routeHandlers, plugin...)
}

func (s *HTTPServer) AddMiddlewares(plugin ...gin.HandlerFunc) {
	s.middlewares = append(s.middlewares, plugin...)
}

func (s *HTTPServer) readiness() {
	s.r.GET("/readiness", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
}

func securityHeaders(c *gin.Context) {
	c.Header("X-Frame-Options", "DENY")
	c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
	c.Header("X-XSS-Protection", "1; mode=block")
	c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	c.Header("Referrer-Policy", "strict-origin")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
	c.Next()
}

func accessControl(c *gin.Context) {
	cfg := config.C(config.Env)
	c.Writer.Header().Set("Access-Control-Allow-Origin", cfg.AccessControl.AllowOrigin)
	c.Writer.Header().Set("Access-Control-Request-Method", "POST, GET, PUT, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
	c.Next()
}

func handlerTimeoutMiddleware(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerTimeout)
	defer cancel()
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}

func shutdownServer(server *http.Server, cleanupFns ...func()) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Wait for termination signal
	<-ctx.Done()

	slog.Info("Shutdown signal received...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Error shutting down server:", err)
	} else {
		slog.Info("Server gracefully stopped")
	}

	// Execute optional cleanup functions
	for _, fn := range cleanupFns {
		fn()
	}
	slog.Info("Cleanup tasks completed.")
}
