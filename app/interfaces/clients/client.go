package clients

import "dominus/app/interfaces/database"

type clients struct {
	dsn, database string
}

func (c clients) NewMongoClient() database.RepositoryInt {
	return database.NewRepository(c.dsn, c.database)
}

type ClientInt interface {
	NewMongoClient() database.RepositoryInt
}

func NewClient(dsn, database string) ClientInt {
	return &clients{
		dsn:      dsn,
		database: database,
	}
}