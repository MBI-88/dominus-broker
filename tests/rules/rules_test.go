package rules_test

import (
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/domain/rules"
	"dominus-project/tests/env"
	"encoding/json"
	"testing"
)

func TestRules(t *testing.T) {
	rls := rules.NewValidator()

	t.Run("Rules_OK", func(t *testing.T) {
		var topic entities.Topic 
		if err := json.Unmarshal(env.ReadJson("./../mocks/rest_create_body.json"), &topic); err != nil {
			t.Fatal(err)
		}

		if err := rls.ValidateStruct(&topic); err != nil {
			t.Fatal(err)
		}

	})

	t.Run("Rules_ERROR", func(t *testing.T) {
		var topic entities.Topic 
		if err := json.Unmarshal(env.ReadJson("./../mocks/rest_create_body_error.json"), &topic); err != nil {
			t.Fatal(err)
		}

		if err := rls.ValidateStruct(&topic); err == nil {
			t.Fatal("Error is empty")
		}
	})
}