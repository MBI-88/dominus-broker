package database

import (
	"context"
	"dominus/app/domain/rules"
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



func (r repository) Migrations(strCollections string) {
	collections := strings.Split(strCollections, ",")
	
	for _, name := range collections {
		if name == "topic_db" {
			r.rls.CreateIndex(r.client, r.database, name, context.TODO())
		} else {
			if err := r.client.Database(r.database).CreateCollection(context.TODO(), name); err != nil {
				panic(err)
			}
		}
	}

	fmt.Println("[*] Migration successful!")
}


func (r repository) FindObject(f primitive.D, collection string, object any) error {
	cl := r.client.Database(r.database).Collection(collection)
	cursor := cl.FindOne(context.TODO(), f)
	
	if err := cursor.Decode(object); err != nil {
		return err
	}

	return nil
}

func (r repository) FindObjects(collection string, objects any, filter bson.D, op ...*options.FindOptions) error {
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

func (r repository) InsertObject(object any, collection string) (*mongo.InsertOneResult, error) {
	cl := r.client.Database(r.database).Collection(collection)
	result, err := cl.InsertOne(context.TODO(), object)
	if err != nil {
		return nil, err
	}
	return result, nil
}


func (r repository) UpdateObject(filter primitive.D, updadte primitive.D ,collection string) (*mongo.UpdateResult, error) {
	cl := r.client.Database(r.database).Collection(collection)
	result, err := cl.UpdateOne(context.TODO(), filter, updadte)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r repository) DeleteObject(f primitive.D, collection string) (*mongo.DeleteResult, error) {
	cl := r.client.Database(r.database).Collection(collection)
	result, err := cl.DeleteOne(context.TODO(), f)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r repository) DeleteObjects(f primitive.D, collection string) (*mongo.DeleteResult, error) {
	cl := r.client.Database(r.database).Collection(collection)
	result, err := cl.DeleteMany(context.TODO(), f)
	if err != nil {
		return nil, err
	}
	return result, nil
}


func (r repository) CountPages(collection string) (int64, error) {
	cl := r.client.Database(r.database).Collection(collection)
	total, err := cl.EstimatedDocumentCount(context.TODO())
	if err != nil {
		return 0, err
	}
	return total, nil
}






type RepositoryInt interface {
	//Make migrations in the database
	//
	//Parameters
	//
	//* strCollection: string that contains collection names splited by ","
	Migrations(strCollection string)
	//Delete an Object in the database
	//
	//Parameters
	//
	//* f: filter to use
	//
	//* collection: name of the collection
	DeleteObject(f primitive.D, collection string) (*mongo.DeleteResult, error)
	//Update an object in the database
	//
	//Parameters
	//
	//* filter: filter to select objects to update
	//
	//* update: the object and key to update
	//
	//* collection: collection name to use
	UpdateObject(filter primitive.D, updadte primitive.D ,collection string) (*mongo.UpdateResult, error)
	//Create an object in the database
	//
	//Parameters
	//
	//* obj: the object to be updated
	//
	//* collection: collection name to use
	InsertObject(obj any, collection string) (*mongo.InsertOneResult, error)
	//Find objects in the database
	// 
	//Parameters
	//
	//* collection: collection name to use
	//
	//* objects: array object to fill
	//
	//* filter: the filter to find objects
	//
	//* op: contains options to use in the query
	FindObjects(collection string, objects any, filter bson.D, op ...*options.FindOptions) error
	//Find and object in the database
	//
	//Parameters
	//
	//* f: filter to match with objects
	//
	//* collection: collection name to use
	//
	//* object: the object to fill
	FindObject(f primitive.D, collection string, object any) error
	//Count pages in the database
	//
	//Parameters
	//
	//* collection: collection name to use
	CountPages(collection string) (int64, error)
	//Delete objects in the database
	//
	//Parameters
	//
	//* f: filter to find objects
	//
	//* collection: collection name to use
	DeleteObjects(f primitive.D, collection string) (*mongo.DeleteResult, error)
}

func NewRepository(dsn, database string, r rules.RuleInt, c *mongo.Client) RepositoryInt {
	return &repository{
		dsn:      dsn,
		database: database,
		rls: r,
		client: c,
	}
}
