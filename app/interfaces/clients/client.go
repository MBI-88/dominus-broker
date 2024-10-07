package clients

import (
	"dominus/app/domain/rules"
	"dominus/app/interfaces/database"
	"dominus/app/interfaces/rest/output"

	"go.mongodb.org/mongo-driver/mongo"
)

type clients struct {
	dsn      string
	database string
	client   *mongo.Client
	rl       rules.RuleInt
}

func (c clients) MongoClient() database.RepositoryInt {
	return database.NewRepository(c.dsn, c.database, c.rl, c.client)
}

func (c clients) RestClient() output.RestClientInt {
	return output.NewRestClient()
}

type ClientInt interface {
	//Mongo client to connect with mongoDB
	MongoClient() database.RepositoryInt
	//Rest client
	RestClient() output.RestClientInt
}

func NewClient(dsn, database string, c *mongo.Client) ClientInt {
	return &clients{
		dsn:      dsn,
		database: database,
		client:   c,
		rl:       rules.NewRule(),
	}
}
