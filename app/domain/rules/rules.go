package rules

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Error Response
type errorResponse struct {
	Error       bool
	FailedField string
	Tag         string
	Value       any
}



type rules struct {
	validate *validator.Validate
}

func (r *rules) ValidateStruct(data any) error {
	var errors []errorResponse
	if err := r.validate.Struct(data); err != nil {
		for _, er := range err.(validator.ValidationErrors) {
			var element errorResponse
			element.FailedField = er.Field()
			element.Tag = er.Tag()
			element.Value = er.Value()
			element.Error = true
			errors = append(errors, element)
		}
	}
	if errors != nil {
		errMsgs := make([]string, 0, 50)
		for _, er := range errors {
			errMsgs = append(errMsgs, fmt.Sprintf("%s: %v disagreement with %s", er.FailedField, er.Value, er.Tag))
		}
		return fmt.Errorf("%v", errMsgs)
	}
	return nil
}

type RulesInt interface {
	ValidateStruct(data any) error 
}

func NewValidator() RulesInt {
	return &rules{
		validate: validator.New(),
	}
}