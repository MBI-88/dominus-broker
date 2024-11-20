package rules

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type rules struct {
	re *regexp.Regexp
}

func (*rules) CreateIndex(client *mongo.Client, database, name string, ctx context.Context) {
	db := client.Database(database)
	if err := db.CreateCollection(ctx, name); err != nil {
		panic(err)
	}
	collection := db.Collection(name)
	indexModel := mongo.IndexModel{
		Keys:    map[string]int{"topic": 1},
		Options: options.Index().SetUnique(true),
	}

	if _, err := collection.Indexes().CreateOne(ctx, indexModel); err != nil {
		panic(err)
	}
}

func (r *rules) CheckURI(uri string) bool {
	return r.re.MatchString(uri)
}

func (*rules) Paginator(page, size uint64) *options.FindOptions {
	if page < 0 || size < 0 {
		page = 2
		size = 4
	}
	opts := options.Find()
	opts.SetSkip(int64(page-1) * int64(size))
	opts.SetLimit(int64(size))
	return opts
}

func (r *rules) MakeLogFiter(filters map[string]string, page, size uint64) []bson.D {
	var (
		arrayFilter = make([]bson.D, 0, 50)
		opts = r.Paginator(page, size) 
		step int
	)

	if *opts.Limit == 0 {
		step = 10
	} else {
		step = int(*opts.Limit)
	}

	for key, val := range filters {
		switch key {
		case "subscriber":
			if val != "" {
				arrayFilter = append(arrayFilter, bson.D{
					primitive.E{Key: key, Value: primitive.Regex{Pattern: fmt.Sprintf("^%s", val)}},
				})
			}
		case "stage":
			if val != "" {
				arrayFilter = append(arrayFilter, bson.D{
					primitive.E{Key: key, Value: primitive.Regex{Pattern: fmt.Sprintf("^%s", val)}},
				})
			}
		case "start":
			if val != "" {
				if timestap, err := time.Parse("2006-01-02 15:04:05", val); err == nil {
					arrayFilter = append(arrayFilter, bson.D{
						primitive.E{Key: "$match", Value: bson.D{
							primitive.E{Key: "created_at", Value: bson.D{
								primitive.E{Key: "$gte", Value: timestap},
							}},
						}}},
					)
				}
			}
		case "end":
			if val != "" {
				if timestap, err := time.Parse("2006-01-02 15:04:05", val); err == nil {
					arrayFilter = append(arrayFilter, bson.D{
						primitive.E{Key: "$match", Value: bson.D{
							primitive.E{Key: "created_at", Value: bson.D{
								primitive.E{Key: "$lte", Value: timestap},
							}},
						}}},
					)
				}
			}
		}

	}
	bsSkip := bson.D{primitive.E{Key: "$skip", Value: *opts.Skip}}
	bsLimit := bson.D{primitive.E{Key: "$limit", Value: step}}
	arrayFilter = append(arrayFilter, bsSkip, bsLimit)
	return arrayFilter
}

func (*rules) MakeEmptyFilter() bson.D {
	return bson.D{}
}

type RulesInt interface {
	CreateIndex(client *mongo.Client, database, name string, ctx context.Context)
	CheckURI(uri string) bool
	MakeLogFiter(filters map[string]string, page, size uint64) []bson.D
	Paginator(page, size uint64) *options.FindOptions
	MakeEmptyFilter() bson.D
}

func NewRule() RulesInt {
	return &rules{
		re: regexp.MustCompile(`^(https?:\/\/[a-zA-Z0-9.-]+)(:\d{1,5})?(\/[^\s]*)?$`),
	}
}
