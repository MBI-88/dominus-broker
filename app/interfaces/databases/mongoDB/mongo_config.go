package mongoDB

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoConfig struct {
	isActive bool
}

func (m mongoConfig) CreateClient(dsn string) *mongo.Client {
	if m.isActive {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
	return nil
}

type MongoConfigInt interface {
	CreateClient(dsn string) *mongo.Client
}

func NewMongoConfig(isActive bool) MongoConfigInt {
	return &mongoConfig{
		isActive: isActive,
	}
}
