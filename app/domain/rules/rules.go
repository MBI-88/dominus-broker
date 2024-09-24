package rules

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type errorResponse struct {
	Error        bool
	FaildedField string
	Tag          string
	Value        any
}

type rule struct {
	v  *validator.Validate
	ex regexp.Regexp
}

func (r rule) ValidateMessage(data any) error {
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

func (r rule) ValidateTopic(topic map[string][]string) error {
	for key, val := range topic {
		for _, url := range val {
			if !r.ex.MatchString(url) {
				return fmt.Errorf("Not acceptable %s:[%s]", key, url)
			}
		}
	}
	return nil
}

type RuleInt interface {
	ValidateMessage(data any) error
	ValidateTopic(sub map[string][]string) error
}

func NewRule() RuleInt {
	return &rule{v: validator.New(), ex: *regexp.MustCompile(`^(https?)://([a-zA-Z0-9-.]+)\.([a-zA-Z]{2,})`)}
}
