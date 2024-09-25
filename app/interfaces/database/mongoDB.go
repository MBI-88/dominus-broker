package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type repository struct {
	client *mongo.Client
	database string
	dsn string
}

func (r *repository) connectDB(dsn string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(dsn)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil  {
		panic(err)
	}

	r.client = client
	if err = r.client.Ping(ctx,nil); err != nil {
		panic(err)
	}

}

func (r *repository) Connection() {
	r.connectDB(r.dsn)
}

func (r repository) createIndexes(name, filed string) {
	if err := r.client.Database(r.database).CreateCollection(context.TODO(), name); err != nil {
		panic(err)
	}

	collection := r.client.Database(r.database).Collection(name)
	indexModel := mongo.IndexModel{
		Keys: map[string]int{filed:1},
		Options: options.Index().SetName(filed),
	}

	if _, err := collection.Indexes().CreateOne(context.TODO(), indexModel); err != nil {
		panic(err)
	}
}

func (r repository) Migrations(strCollections string) {
	collections := strings.Split(strCollections, ",")
	r.connectDB(r.dsn)

	for _, name := range collections {
		if name == "" {
			r.createIndexes(name, "")
		}else {
			if err := r.client.Database(r.database).CreateCollection(context.TODO(),name); err != nil {
				panic(err)
			}
		}
	}

	fmt.Println("[*] Migration successful!")
}

func (r repository) Client(collection string) *mongo.Collection {
	return r.client.Database(r.database).Collection(collection)
}



func (r *repository) SetVars(dsn, database string) {
	r.dsn = dsn 
	r.database = database
}




type RepositoryInt interface {
	// Create migrations 
	//
	// strCollections: string of collections name join by ",". It's going to be split by ","
	Migrations(strCollections string)
	// Create a client and set database name
	Connection()
	// Client to use in any connection
	Client(collection string) *mongo.Collection
	// Set system variables
	//
	// dsn: url for connectiong
	//
	// database: name of the database to connect
	SetVars(dsn, database string)
}


func NewRepository(dsn, database string) RepositoryInt {
	return  &repository{
		dsn : dsn,
		database: database,
	}
}