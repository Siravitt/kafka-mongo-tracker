package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Server        Server
	Database      Database
	KafKa         KafKa
	AccessControl AccessControl
	Header        Header
}

type Server struct {
	HostName string `env:"HOSTNAME"`
	Port     string `env:"PORT,notEmpty"`
}

type Database struct {
	MongoURL string `env:"MONGO_URL"`
}

type KafKa struct {
	ProducerURL []string `env:"KAFKA_PRODUCER_URL" envSeparator:","`
	ConsumerURL []string `env:"KAFKA_CONSUMER_URL" envSeparator:","`
}

type AccessControl struct {
	AllowOrigin string `env:"ACCESS_CONTROL_ALLOW_ORIGIN"`
}

type Header struct {
	RefIDHeaderKey string `env:"REF_ID_HEADER_KEY,notEmpty"`
}

var once sync.Once
var config Config

func prefix(e string) string {
	if e == "" {
		return ""
	}

	return fmt.Sprintf("%s_", e)
}

func parseEnv[T any](opts env.Options) (T, error) {
	var t T

	if err := env.Parse(&t); err != nil {
		return t, err
	}

	env.ParseWithOptions(&t, opts)

	return t, nil
}

func C(envPrefix string) Config {
	once.Do(func() {
		opts := env.Options{
			Prefix: prefix(envPrefix),
		}

		var err error
		config, err = parseEnv[Config](opts)
		if err != nil {
			log.Fatal(err)
		}
	})

	return config
}
