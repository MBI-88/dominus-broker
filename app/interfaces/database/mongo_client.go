package database

import (
	"context"
	"dominus/app/domain/rules"
	"dominus/app/interactors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type repository struct {
	client   *mongo.Client
	database string
	dsn      string
	rls      rules.RuleInt
}



func (r *repository) Migrations(strCollections string) {
	collections := strings.Split(strCollections, ",")
	for _, name := range collections {
		if name == "topic" {
			r.rls.CreateIndex(r.client, r.database, name, context.TODO())
		} else {
			if err := r.client.Database(r.database).CreateCollection(context.TODO(), name); err != nil {
				panic(err)
			}
		}
	}
	fmt.Println("[*] Migration successful!")
}


func (r *repository) FindObject(f primitive.D, collection string, object any) error {
	cl := r.client.Database(r.database).Collection(collection)
	cursor := cl.FindOne(context.TODO(), f)
	
	if err := cursor.Decode(object); err != nil {
		return err
	}

	return nil
}

func (r *repository) FindObjects(collection string, objects any, filter bson.D, op ...*options.FindOptions) error {
	cl := r.client.Database(r.database).Collection(collection)
	cursor, err := cl.Find(context.TODO(), filter , op...)
	defer cursor.Close(context.TODO())

	if err != nil {
		return err
	}
	if err := cursor.All(context.TODO(), objects); err != nil {
		return err
	}

	return nil
}

func (r *repository) InsertObject(object any, collection string) (*mongo.InsertOneResult, error) {
	cl := r.client.Database(r.database).Collection(collection)
	result, err := cl.InsertOne(context.TODO(), object)
	if err != nil {
		return nil, err
	}
	return result, nil
}


func (r *repository) UpdateObject(filter primitive.D, updadte primitive.D ,collection string) (*mongo.UpdateResult, error) {
	cl := r.client.Database(r.database).Collection(collection)
	result, err := cl.UpdateOne(context.TODO(), filter, updadte)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repository) DeleteObject(f primitive.D, collection string) (*mongo.DeleteResult, error) {
	cl := r.client.Database(r.database).Collection(collection)
	result, err := cl.DeleteOne(context.TODO(), f)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repository) DeleteObjects(f primitive.D, collection string) (*mongo.DeleteResult, error) {
	cl := r.client.Database(r.database).Collection(collection)
	result, err := cl.DeleteMany(context.TODO(), f)
	if err != nil {
		return nil, err
	}
	return result, nil
}


func (r *repository) CountPages(collection string) (int64, error) {
	cl := r.client.Database(r.database).Collection(collection)
	total, err := cl.EstimatedDocumentCount(context.TODO())
	if err != nil {
		return 0, err
	}
	return total, nil
}



func NewRepository(dsn, database string, r rules.RuleInt, c *mongo.Client) interactors.RepositoryInt {
	return &repository{
		dsn:      dsn,
		database: database,
		rls: r,
		client: c,
	}
}
