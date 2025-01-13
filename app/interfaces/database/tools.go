package database

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoTools struct {
}

func (*mongoTools) CreateIndex() []mongo.IndexModel {
	var indexModel []mongo.IndexModel
	for _, key := range []string{"created_at", "stage", "subscriber"} {
		index := mongo.IndexModel{
			Keys:    map[string]int{key: 1},
			Options: options.Index().SetName(key),
		}
		indexModel = append(indexModel, index)
	}
	return indexModel
}


func (m *mongoTools) MakeBackupFilter(filters map[string]string) []bson.D {
	return m.makeAgregation(filters)
}

func (*mongoTools) Paginator(page, size uint64) *options.FindOptions {
	if page <= 0 || size <= 0 {
		page = 2
		size = 4
	}
	opts := options.Find()
	opts.SetSkip(int64(page-1) * int64(size))
	opts.SetLimit(int64(size))
	return opts
}

func (m *mongoTools) MakeLogFiter(filters map[string]string, page, size uint64) []bson.D {
	var (
		opts        = m.Paginator(page, size)
		step        int
	)

	if *opts.Limit == 0 {
		step = 10
	} else {
		step = int(*opts.Limit)
	}
	arrayFilter := m.makeAgregation(filters)
	bsSkip := bson.D{primitive.E{Key: "$skip", Value: *opts.Skip}}
	bsLimit := bson.D{primitive.E{Key: "$limit", Value: step}}
	arrayFilter = append(arrayFilter, bsSkip, bsLimit)
	return arrayFilter
}


func (m *mongoTools) makeAgregation(filters map[string]string) []bson.D {
	var arrayFilter = make([]bson.D, 0, 100)
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
	return arrayFilter
}


type mongoToolsInt interface {
	CreateIndex() []mongo.IndexModel
	MakeBackupFilter(filters map[string]string) []bson.D
	Paginator(page, size uint64) *options.FindOptions
	MakeLogFiter(filters map[string]string, page, size uint64) []bson.D
}

func newMongoTools() mongoToolsInt {
	return new(mongoTools)
}