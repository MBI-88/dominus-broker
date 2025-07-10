package interactors_test

import (
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/interactors/manager"
	"dominus-project/tests/env"
	"testing"
)

func TestManagerCases(t *testing.T) {
	topics := entities.NewTopics(env.TopicLimit)
	manager := manager.NewManagerService(topics, env.QueueLimit)

	// Create
	t.Run("AddTopic_Ok", func(t *testing.T) {
		dto := env.NewRestConext(true, "./../../mocks/rest_create_body.json", "")

		if err := manager.AddTopic(dto); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("AddTopic_Error_Parse", func(t *testing.T) {
		dto := env.NewRestConext(false, "", "")

		if err := manager.AddTopic(dto); err == nil {
			t.Fatal("Error must be different from nil")
		}
	})
	t.Run("AddTopic_Error_ValidStruct", func(t *testing.T) {
		dto := env.NewRestConext(true, "./../../mocks/rest_create_body_error.json", "")

		if err := manager.AddTopic(dto); err == nil {
			t.Fatal("Error must be different from nil")
		}
	})

	t.Run("AddTopic_Error_EmptyName", func(t *testing.T) {
		dto := env.NewRestConext(true, "./../../mocks/rest_create_body_empty_name.json", "")

		if err := manager.AddTopic(dto); err == nil {
			t.Fatal("Error must be different from nil")
		}
	})

	t.Run("AddTopic_Error_TopicExist", func(t *testing.T) {
		dto := env.NewRestConext(true, "./../../mocks/rest_create_body_exist.json", "")

		if err := manager.AddTopic(dto); err == nil {
			t.Fatal("Error must be different from nil")
		}
	})

	// Update
	t.Run("UpdatePartition_Ok", func(t *testing.T) {
		dto := env.NewRestConext(true, "./../../mocks/rest_update_body.json", "test")

		if err := manager.UpdateSubscribers(dto); err != nil {
			t.Fatal(err)
		}

	})
	t.Run("UpdatePartition_ErrorParse", func(t *testing.T) {
		dto := env.NewRestConext(false, "", "")

		if err := manager.UpdateSubscribers(dto); err == nil {
			t.Fatal("Error must be different from nil")
		}

	})
	t.Run("UpdatePartition_Error_ValidStruct", func(t *testing.T) {
		dto := env.NewRestConext(true, "./../../mocks/rest_update_error.json", "test")

		if err := manager.UpdateSubscribers(dto); err == nil {
			t.Fatal("Error must be different from nil")
		}

	})
	t.Run("UpdateSubscribers_Error", func(t *testing.T) {
		dto := env.NewRestConext(true, "./../../mocks/rest_update_body.json", "")

		if err := manager.UpdateSubscribers(dto); err == nil {
			t.Fatal("Error must be different from nil")
		}

	})

	// Delete
	t.Run("DeleteTopic_Ok", func(t *testing.T) {
		dto := env.NewRestConext(true, "", "test")

		if err := manager.DeleteTopic(dto); err != nil {
			t.Fatal(err)
		}

	})

	t.Run("DeleteTopic_Error", func(t *testing.T) {
		dto := env.NewRestConext(false, "", "")

		if err := manager.DeleteTopic(dto); err == nil {
			t.Fatal("Error must be different from nil")
		}
	})

}
