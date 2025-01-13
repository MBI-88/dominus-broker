package database

import (
	"context"
	"dominus/app/domain/repos"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type repository struct {
	client      *mongo.Client
	database    string
	dsn         string
	tl          mongoToolsInt
	collections string
}

func (r *repository) Migrations() {
	collections := strings.Split(r.collections, ",")
	for _, name := range collections {
		if name == "logs" {
			index := r.tl.CreateIndex()
			db := r.client.Database(name)
			if err := db.CreateCollection(context.TODO(), name); err != nil {
				panic(err)
			}
			if _, err := db.Collection(name).Indexes().
				CreateMany(context.TODO(), index); err != nil {
				panic(err)
			}
		} else {
			if err := r.client.Database(r.database).CreateCollection(context.TODO(), name); err != nil {
				panic(err)
			}
		}
	}
	fmt.Println("[*] Migration successful!")
}


func (r *repository) InsertObject(obj any, collection string) error {
	cl := r.client.Database(r.database).Collection(collection)
	_, err := cl.InsertOne(context.TODO(), obj)
	if err != nil {
		return err
	}
	return nil
}


func (r *repository) DeleteObjects(collection string) error {
	cl := r.client.Database(r.database).Collection(collection)
	_, err := cl.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) CountPages(collection string) (int64, error) {
	cl := r.client.Database(r.database).Collection(collection)
	return cl.EstimatedDocumentCount(context.TODO())
}

func (r *repository) Filter(filter map[string]string, object any, collection string, page, size uint64) error {
	allowDiskUse := true
	options := &options.AggregateOptions{
		AllowDiskUse: &allowDiskUse,
	}
	rfilter := r.tl.MakeLogFiter(filter, page, size)
	cl := r.client.Database(r.database).Collection(collection)
	cursor, err := cl.Aggregate(context.TODO(), rfilter, options)
	defer cursor.Close(context.TODO())
	if err != nil {
		return err
	}

	if err := cursor.All(context.TODO(), object); err != nil {
		return err
	}
	return nil
}

func (r *repository) Backup(filter map[string]string, object any, collection string) error {
	allowDiskUse := true
	options := &options.AggregateOptions{
		AllowDiskUse: &allowDiskUse,
	}
	rfilter := r.tl.MakeBackupFilter(filter)
	cl := r.client.Database(r.database).Collection(collection)
	cursor, err := cl.Aggregate(context.TODO(), rfilter, options)
	defer cursor.Close(context.TODO())
	if err != nil {
		return err
	}
	if err := cursor.All(context.TODO(), object); err != nil {
		return err
	}
	return nil
}
 
func (r *repository) Stats() (any, error) {
	var result bson.M
	if err := r.client.Database(r.database).
		RunCommand(context.TODO(), bson.D{{Key: "dbStats", Value: 1}}).
		Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func NewRepository(dsn, database, cols string,c *mongo.Client) repos.RepositoryInt {
	return &repository{
		dsn:         dsn,
		database:    database,
		collections: cols,
		client:      c,
		tl:          newMongoTools(),
	}
}
