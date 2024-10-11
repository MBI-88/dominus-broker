package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoConfig struct{}

func (*mongoConfig) CreateClient(dsn string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(dsn)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		panic(err)
	}

	return client
}

type MongoConfigInt interface {
	CreateClient(dsn string) *mongo.Client 
}


func NewMongoConfig() MongoConfigInt {
	return new(mongoConfig)
}