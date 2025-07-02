package rules

import (
	"fmt"
	"net"

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
	customValidate := validator.New()
	customValidate.RegisterValidation("hostname_port", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		
		host, port, _ := net.SplitHostPort(value)

		isValidPort := func(port string) bool {
			if p, err := net.LookupPort("tcp", port); err == nil {
				return p > 0 && p <= 65535
			}
			return false
		}

		if net.ParseIP(host) != nil && isValidPort(port) {
			return  true
		}

		if err := validator.New().Var(value, "hostname"); err == nil {
			return  true
		}

		return false

	})
	return &rules{
		validate: customValidate,
	}
}
