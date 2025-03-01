package mongoDB

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
	active      bool
}

func (r *repository) Migrations() {
	if r.active {
		collections := strings.SplitSeq(r.collections, ",")
		for name := range collections {
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
		return
	}
	fmt.Println("[-] Database inactive")
}

func (r *repository) InsertObject(obj any, collection string) error {
	if r.active {
		cl := r.client.Database(r.database).Collection(collection)
		_, err := cl.InsertOne(context.TODO(), obj)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("[-] Database inactive")
}

func (r *repository) DeleteObjects(collection string) error {
	if r.active {
		cl := r.client.Database(r.database).Collection(collection)
		_, err := cl.DeleteMany(context.TODO(), bson.D{})
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("[-] Database inactive")
}

func (r *repository) CountPages(collection string) (int64, error) {
	if r.active {
		cl := r.client.Database(r.database).Collection(collection)
		return cl.EstimatedDocumentCount(context.TODO())
	}
	return 0, fmt.Errorf("[-] Database inactive")
}

func (r *repository) Filter(filter map[string]string, object any, collection string, page, size int) error {
	if r.active {
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
	return fmt.Errorf("[-] Database inactive")
}

func (r *repository) Backup(filter map[string]string, object any, collection string) error {
	if r.active {
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
	return fmt.Errorf("[-] Database inactive")
}

func (r *repository) Stats() (any, error) {
	if r.active {
		var result bson.M
		if err := r.client.Database(r.database).
			RunCommand(context.TODO(), bson.D{{Key: "dbStats", Value: 1}}).
			Decode(&result); err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, fmt.Errorf("[-] Database inactive")
}

func NewRepository(dsn, database, cols string, c *mongo.Client, active bool) repos.RepositoryInt {
	return &repository{
		dsn:         dsn,
		database:    database,
		collections: cols,
		client:      c,
		tl:          newMongoTools(),
		active:      active,
	}
}
