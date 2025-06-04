package main

import (
	"github.com/Siravitt/kafka-mongo-tracker/config"
	"github.com/Siravitt/kafka-mongo-tracker/database"
	"github.com/Siravitt/kafka-mongo-tracker/logger"
	"github.com/Siravitt/kafka-mongo-tracker/server"
)

func main() {
	cfg := config.C(config.Env)
	logger.New(logger.GCPKeyReplacer)

	mongodb, dbClose := database.NewMongoDB(cfg.Database.MongoURL)
	defer dbClose()

	_ = mongodb

	app := server.NewRouter(cfg)
	app.StartServer()
}
