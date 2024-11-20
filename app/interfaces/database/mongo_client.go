package database

import (
	"context"
	"dominus/app/domain/rules"
	"dominus/app/interactors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
	client      *mongo.Client
	database    string
	dsn         string
	rls         rules.RuleInt
	collections string
}

func (r *repository) Migrations() {
	collections := strings.Split(r.collections, ",")
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

func (r *repository) FindObject(f any, collection string, object any) error {
	cl := r.client.Database(r.database).Collection(collection)
	cursor := cl.FindOne(context.TODO(), f)

	if err := cursor.Decode(object); err != nil {
		return err
	}
	return nil
}

func (r *repository) FindObjects(collection string, objects any, filter any, page, sizze int) error {
	cl := r.client.Database(r.database).Collection(collection)
	cursor, err := cl.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return err
	}
	if err := cursor.All(context.TODO(), objects); err != nil {
		return err
	}
	return nil
}

func (r *repository) InsertObject(object any, collection string) error {
	cl := r.client.Database(r.database).Collection(collection)
	_, err := cl.InsertOne(context.TODO(), object)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateObject(filter any, updadte any, collection string) error {
	cl := r.client.Database(r.database).Collection(collection)
	_, err := cl.UpdateOne(context.TODO(), filter, updadte)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) DeleteObject(f any, collection string) error {
	cl := r.client.Database(r.database).Collection(collection)
	_, err := cl.DeleteOne(context.TODO(), f)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) DeleteObjects(f any, collection string) error {
	cl := r.client.Database(r.database).Collection(collection)
	_, err := cl.DeleteMany(context.TODO(), f)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) CountPages(collection string) (int64, error) {
	cl := r.client.Database(r.database).Collection(collection)
	total, err := cl.EstimatedDocumentCount(context.TODO())
	if err != nil {
		return 0, err
	}
	return total, nil
}

func NewRepository(dsn, database, cols string, r rules.RuleInt, c *mongo.Client) interactors.RepositoryInt {
	return &repository{
		dsn:         dsn,
		database:    database,
		collections: cols,
		rls:         r,
		client:      c,
	}
}
