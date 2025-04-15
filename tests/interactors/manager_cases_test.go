package interactors_test

import (
	"dominus-project/app/domain/entities"
	"dominus-project/app/interactors"
	"dominus-project/tests/mocks"
	"testing"
)

func TestManagerCases(t *testing.T) {
	log := mocks.NewEventMock(true)
	topics := entities.NewTopics(100)
	inter := interactors.NewInteractor(log, topics, 100)
	manager := inter.NewManagerService()

	// Create
	t.Run("AddTopic_Ok", func(t *testing.T) {
		dto := mocks.NewRestConext(true, "./../mocks/rest_create_body.json", "", nil)

		if err := manager.AddTopic(dto); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("AddTopic_Error_Parse", func(t *testing.T) {
		dto := mocks.NewRestConext(false, "", "", nil)

		if err := manager.AddTopic(dto); err == nil {
			t.Fatal("Error must be different from nil")
		}
	})
	t.Run("AddTopic_Error_ValidStruct", func(t *testing.T) {
		dto := mocks.NewRestConext(true, "./../mocks/rest_create_body_error.json", "", nil)

		if err := manager.AddTopic(dto); err == nil {
			t.Fatal("Error must be different from nil")
		}
	})

	t.Run("AddTopic_Error_EmptyName", func(t *testing.T) {
		dto := mocks.NewRestConext(true, "./../mocks/rest_create_body_empty_name.json", "", nil)

		if err := manager.AddTopic(dto); err == nil {
			t.Fatal("Error must be different from nil")
		}
	})

	t.Run("AddTopic_Error_TopicExist", func(t *testing.T) {
		dto := mocks.NewRestConext(true, "./../mocks/rest_create_body_exist.json", "", nil)

		if err := manager.AddTopic(dto); err == nil {
			t.Fatal("Error must be different from nil")
		}
	})

	// Update 
	t.Run("UpdatePartition_Ok", func(t *testing.T) {
		dto := mocks.NewRestConext(true, "./../mocks/rest_update_body.json", "test", nil)

		if err := manager.UpdatePartition(dto); err != nil {
			t.Fatal(err)
		}

	})
	t.Run("UpdatePartition_ErrorParse", func(t *testing.T) {
		dto := mocks.NewRestConext(false, "", "", nil)

		if err := manager.UpdatePartition(dto); err == nil {
			t.Fatal("Error must be different from nil")
		}

	})
	t.Run("UpdatePartition_Error_ValidStruct", func(t *testing.T) {
		dto := mocks.NewRestConext(true, "./../mocks/rest_update_error.json", "test", nil)

		if err := manager.UpdatePartition(dto); err == nil {
			t.Fatal("Error must be different from nil")
		}

	})
	t.Run("UpdatePartition_Error", func(t *testing.T) {
		dto := mocks.NewRestConext(true, "./../mocks/rest_update_body.json", "", nil)

		if err := manager.UpdatePartition(dto); err == nil {
			t.Fatal("Error must be different from nil")
		}

	})

	// Delete 
	t.Run("DeleteTopic_Ok", func(t *testing.T) {
		dto := mocks.NewRestConext(true, "", "test", nil) 

		if err := manager.DeleteTopic(dto); err != nil {
			t.Fatal(err)
		}

	})
	
	t.Run("DeleteTopic_Error", func(t *testing.T) {
		dto := mocks.NewRestConext(false, "", "", nil)

		if err := manager.DeleteTopic(dto); err == nil {
			t.Fatal("Error must be different from nil")
		}
	})
	
}