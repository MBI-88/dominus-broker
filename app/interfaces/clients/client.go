package clients

import (
	"dominus/app/domain/rules"
	"dominus/app/interfaces/database"

	"go.mongodb.org/mongo-driver/mongo"
)

type clients struct {
	dsn      string
	database string
	client   *mongo.Client
	rl       rules.RuleInt
}

func (c clients) NewMongoClient() database.RepositoryInt {
	return database.NewRepository(c.dsn, c.database, c.rl, c.client)
}

type ClientInt interface {
	//New mongo client to connect with mongoDB
	NewMongoClient() database.RepositoryInt
}

func NewClient(dsn, database string, c *mongo.Client) ClientInt {
	return &clients{
		dsn:      dsn,
		database: database,
		client:   c,
		rl:       rules.NewRule(),
	}
}
