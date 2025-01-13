package rules

import (
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type rules struct {
	re *regexp.Regexp
}


func (r *rules) CheckURI(uri string) bool {
	return r.re.MatchString(uri)
}

func (*rules) MakeMongoID() primitive.ObjectID {
	return primitive.NewObjectID()
}

type RulesInt interface {
	CheckURI(uri string) bool
	MakeMongoID() primitive.ObjectID
}

func NewRules() RulesInt {
	return &rules{
		re: regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}:\d{1,5}$`),
	}
}
