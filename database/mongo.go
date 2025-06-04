package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewMongoDB(url string) (*mongo.Client, func()) {
	options := options.Client().
		ApplyURI(url).
		SetConnectTimeout(mongoConnectTimeout).
		SetMaxConnIdleTime(connMaxIdleTime).
		SetMaxConnecting(0). // use default = 2
		SetSocketTimeout(mongoSocketTimeout).
		SetTimeout(mongoTransactionTimeout).
		SetMaxPoolSize(mongoMaxPoolSize)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect to mongo
	client, err := mongo.Connect(ctx, options)
	if err != nil {
		log.Panic("error while creating connection to the database!!", err)
	}

	// check the connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Panic("could not ping database", err)
	}

	// clean up func
	fnDBClose := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Panic("error while closing database", err)
			return
		}
	}

	return client, fnDBClose
}
