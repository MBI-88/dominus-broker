package rules

import (
	"context"
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type errorResponse struct {
	Error        bool
	FaildedField string
	Tag          string
	Value        any
}

type rule struct {
	v  *validator.Validate
	re  *regexp.Regexp
}

func (r *rule) ValidateStruct(data any) error {
	var (
		errors  []errorResponse
		element errorResponse
	)
	if err := r.v.Struct(data); err != nil {
		for _, er := range err.(validator.ValidationErrors) {
			element.FaildedField = er.Field()
			element.Tag = er.Tag()
			element.Value = er.Value()
			element.Error = true
			errors = append(errors, element)
		}
	}

	if errors != nil {
		errMsg := make([]string, 0)
		for _, er := range errors {
			errMsg = append(errMsg, fmt.Sprintf("%s: %v disagreement with %s", er.FaildedField, er.Value, er.Tag))
		}
		return fmt.Errorf("%v", errMsg)
	}
	return nil
}


func (*rule) CreateIndex(client *mongo.Client, database, name string, ctx context.Context) {
	db := client.Database(database)
	if err := db.CreateCollection(ctx, name); err != nil {
		panic(err)
	}
	collection := db.Collection(name)
	indexModel := mongo.IndexModel{
		Keys: map[string]int{"topic": 1},
		Options: options.Index().SetUnique(true),
	}

	if _, err := collection.Indexes().CreateOne(ctx, indexModel); err != nil {
		panic(err)
	}
}

func (r *rule) CheckURI(uri string) bool {
	return r.re.MatchString(uri)
}

type RuleInt interface {
	//Create index respect to business logic
	//
	//Parameters
	//
	//-> client: databas client to use
	//
	//-> database: database name to use
	//
	//-> name: collection name selected
	//
	//-> ctx: context to use in the query
	CreateIndex(client *mongo.Client, database, name string, ctx context.Context)
	//Validates message using bussiness logic
	//
	//Parameters
	//
	//-> data: the object to validate
	ValidateStruct(data any) error
	//Validates uri
	//
	//Parameters
	//
	//-> uri: uri to validate
	CheckURI(uri string) bool
}


func NewRule() RuleInt {
	return &rule{
		v: validator.New(),
		re: regexp.MustCompile(`^(https?:\/\/[a-zA-Z0-9.-]+)(:\d{1,5})?(\/[^\s]*)?$`),
	}
}




