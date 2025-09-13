package manager_test

import (
	"dominus-project/internal/interactors/manager"
	"dominus-project/mocks"
	"fmt"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestAddTopic(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*mocks.MockRestDto, *mocks.MockTopics, *mocks.MockTopic)
		output    error
	}{
		{
			name: "AddTopic Ok",
			setupMock: func(mrd *mocks.MockRestDto, mt1 *mocks.MockTopics, mt2 *mocks.MockTopic) {
				mrd.EXPECT().
					BodyParser(gomock.All()).
					DoAndReturn(func(inst any) error {
						temp := reflect.ValueOf(inst).Elem()
						temp.FieldByName("Name").SetString("test")
						subs := temp.FieldByName("Subscribers")
						n1 := reflect.ValueOf("server1.api.com")
						n2 := reflect.ValueOf("server2.api.com")
						news := reflect.Append(subs, n1, n2)
						subs.Set(news)
						return nil
					}).Times(1)

				mt1.EXPECT().
					Find(gomock.All()).
					Return(mt2, fmt.Errorf("not found")).Times(1)

				mt1.EXPECT().
					Append(gomock.All()).
					Return(nil).Times(1)
			},
			output: nil,
		},
		{
			name: "AddTopic parser error",
			setupMock: func(mrd *mocks.MockRestDto, mt1 *mocks.MockTopics, mt2 *mocks.MockTopic) {
				mrd.EXPECT().
					BodyParser(gomock.All()).
					Return(fmt.Errorf("error parser")).Times(1)
			},
			output: fmt.Errorf("parser error"),
		},
		{
			name: "AddTopic validate struct error",
			setupMock: func(mrd *mocks.MockRestDto, mt1 *mocks.MockTopics, mt2 *mocks.MockTopic) {
				mrd.EXPECT().
					BodyParser(gomock.All()).
					DoAndReturn(func(inst any) error {
						temp := reflect.ValueOf(inst).Elem()
						temp.FieldByName("Name").SetString("test-1")
						subs := temp.FieldByName("Subscribers")
						n1 := reflect.ValueOf("server1.api.com")
						n2 := reflect.ValueOf("server2.api.com")
						news := reflect.Append(subs, n1, n2)
						subs.Set(news)
						return nil
					}).Times(1)

			},
			output: fmt.Errorf("validate struct error"),
		},
		{
			name: "AddTopic get name error",
			setupMock: func(mrd *mocks.MockRestDto, mt1 *mocks.MockTopics, mt2 *mocks.MockTopic) {
				mrd.EXPECT().
					BodyParser(gomock.All()).
					DoAndReturn(func(inst any) error {
						temp := reflect.ValueOf(inst).Elem()
						temp.FieldByName("Name").SetString("")
						subs := temp.FieldByName("Subscribers")
						n1 := reflect.ValueOf("server1.api.com")
						n2 := reflect.ValueOf("server2.api.com")
						news := reflect.Append(subs, n1, n2)
						subs.Set(news)
						return nil
					}).Times(1)
			},
			output: fmt.Errorf("get name error"),
		},
		{
			name: "AddTopic find topic ok",
			setupMock: func(mrd *mocks.MockRestDto, mt1 *mocks.MockTopics, mt2 *mocks.MockTopic) {
				mrd.EXPECT().
					BodyParser(gomock.All()).
					DoAndReturn(func(inst any) error {
						temp := reflect.ValueOf(inst).Elem()
						temp.FieldByName("Name").SetString("test")
						subs := temp.FieldByName("Subscribers")
						n1 := reflect.ValueOf("server1.api.com")
						n2 := reflect.ValueOf("server2.api.com")
						news := reflect.Append(subs, n1, n2)
						subs.Set(news)
						return nil
					}).Times(1)

				mt1.EXPECT().
					Find(gomock.All()).
					Return(mt2, nil).Times(1)
			},
			output: fmt.Errorf("find topic ok"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockResDto := mocks.NewMockRestDto(ctrl)
			mockTopics := mocks.NewMockTopics(ctrl)
			mockTopic := mocks.NewMockTopic(ctrl)
			tt.setupMock(mockResDto, mockTopics, mockTopic)

			service := manager.NewManagerService(mockTopics, 10)

			err := service.AddTopic(mockResDto)
			if tt.output != nil && err == nil {
				t.Fatalf("Expected output different from output %s != %s", tt.output, err)
			}
			if tt.output == nil && err != nil {
				t.Fatalf("Expected output different from output %s != %s", tt.output, err)
			}

		})
	}
}

func TestUpdateSubscribers(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*mocks.MockRestDto, *mocks.MockTopics, *mocks.MockTopic)
		output    error
	}{
		{
			name: "UpdateSubscribers Ok",
			setupMock: func(mrd *mocks.MockRestDto, mt1 *mocks.MockTopics, mt2 *mocks.MockTopic) {
				mrd.EXPECT().
					Param("name").
					Return("test").Times(1)

				mrd.EXPECT().
					BodyParser(gomock.All()).
					DoAndReturn(func(inst any) error {
						temp := reflect.ValueOf(inst).Elem()
						subs := temp.FieldByName("Subscribers")
						n1 := reflect.ValueOf("server1.api.com")
						n2 := reflect.ValueOf("server2.api.com")
						news := reflect.Append(subs, n1, n2)
						subs.Set(news)
						return nil
					}).Times(1)

				mt1.EXPECT().
					Update(gomock.All()).
					Return(nil).Times(1)
			},
			output: nil,
		},
		{
			name: "UpdateSubscribers parser error",
			setupMock: func(mrd *mocks.MockRestDto, mt1 *mocks.MockTopics, mt2 *mocks.MockTopic) {
				mrd.EXPECT().
					Param("name").
					Return("test").Times(1)

				mrd.EXPECT().
					BodyParser(gomock.All()).
					Return(fmt.Errorf("error parser")).Times(1)
			},
			output: fmt.Errorf("parser error"),
		},
		{
			name: "UpdateSubscribers validate struct error",
			setupMock: func(mrd *mocks.MockRestDto, mt1 *mocks.MockTopics, mt2 *mocks.MockTopic) {
				mrd.EXPECT().
					Param("name").
					Return("test").Times(1)

				mrd.EXPECT().
					BodyParser(gomock.All()).
					DoAndReturn(func(inst any) error {
						temp := reflect.ValueOf(inst).Elem()
						subs := temp.FieldByName("Subscribers")
						n1 := reflect.ValueOf("server1.api.com")
						n2 := reflect.ValueOf("server2.api.com")
						n3 := reflect.ValueOf("http://:80/api-3")
						news := reflect.Append(subs, n1, n2, n3)
						subs.Set(news)
						return nil
					}).Times(1)

			},
			output: fmt.Errorf("validate struct error"),
		},
		{
			name: "UpdateSubscribers update error",
			setupMock: func(mrd *mocks.MockRestDto, mt1 *mocks.MockTopics, mt2 *mocks.MockTopic) {
				mrd.EXPECT().
					Param("name").
					Return("test").Times(1)

				mrd.EXPECT().
					BodyParser(gomock.All()).
					DoAndReturn(func(inst any) error {
						temp := reflect.ValueOf(inst).Elem()
						subs := temp.FieldByName("Subscribers")
						n1 := reflect.ValueOf("server1.api.com")
						n2 := reflect.ValueOf("server2.api.com")
						news := reflect.Append(subs, n1, n2)
						subs.Set(news)
						return nil
					}).Times(1)

				mt1.EXPECT().
					Update(gomock.All()).
					Return(fmt.Errorf("update error")).Times(1)
			},
			output: fmt.Errorf("update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockResDto := mocks.NewMockRestDto(ctrl)
			mockTopics := mocks.NewMockTopics(ctrl)
			mockTopic := mocks.NewMockTopic(ctrl)
			tt.setupMock(mockResDto, mockTopics, mockTopic)

			service := manager.NewManagerService(mockTopics, 10)

			err := service.UpdateSubscribers(mockResDto)
			if tt.output != nil && err == nil {
				t.Fatalf("Expected output different from output %s != %s", tt.output, err)
			}
			if tt.output == nil && err != nil {
				t.Fatalf("Expected output different from output %s != %s", tt.output, err)
			}

		})
	}
}

func TestDeleteTopic(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*mocks.MockRestDto, *mocks.MockTopics)
		output    error
	}{
		{
			name: "DeleteTopic Ok",
			setupMock: func(mrd *mocks.MockRestDto, mt1 *mocks.MockTopics) {
				mrd.EXPECT().
					Param("name").
					Return("test").Times(1)

				mt1.EXPECT().
					Delete(gomock.All()).
					Return(nil).Times(1)
			},
			output: nil,
		},
		{
			name: "DeleteTopic empty param",
			setupMock: func(mrd *mocks.MockRestDto, mt1 *mocks.MockTopics) {
				mrd.EXPECT().
					Param("name").
					Return("").Times(1)
			},
			output: fmt.Errorf("empty param"),
		},
		{
			name: "DeleteTopic delete error",
			setupMock: func(mrd *mocks.MockRestDto, mt1 *mocks.MockTopics) {
				mrd.EXPECT().
					Param("name").
					Return("test").Times(1)

				mt1.EXPECT().
					Delete(gomock.All()).
					Return(fmt.Errorf("delete error")).Times(1)
			},
			output: fmt.Errorf("delete error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockResDto := mocks.NewMockRestDto(ctrl)
			mockTopics := mocks.NewMockTopics(ctrl)
			tt.setupMock(mockResDto, mockTopics)

			service := manager.NewManagerService(mockTopics, 10)

			err := service.DeleteTopic(mockResDto)
			if tt.output != nil && err == nil {
				t.Fatalf("Expected output different from output %s != %s", tt.output, err)
			}
			if tt.output == nil && err != nil {
				t.Fatalf("Expected output different from output %s != %s", tt.output, err)
			}

		})
	}
}

func TestGetQueueInfo(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*mocks.MockRestDto, *mocks.MockTopics)
		output    int
	}{
		{
			name: "GetQueueInfo resp length 3",
			setupMock: func(mrd *mocks.MockRestDto, mt *mocks.MockTopics) {
				mt.EXPECT().
					GetTopicsInfo().
					Return(map[string]any{
						"test-1": []string{"test1", "test2"},
						"test-2": []string{"test1", "test2"},
						"test-3": []string{"test1", "test2"},
					}).Times(1)
			},
			output: 3,
		},
		{
			name: "GetQueueInfo length 2",
			setupMock: func(mrd *mocks.MockRestDto, mt *mocks.MockTopics) {
				mt.EXPECT().
					GetTopicsInfo().
					Return(map[string]any{
						"test-1": []string{"test1", "test2"},
						"test-2": []string{"test1", "test2"},
					}).Times(1)
			},
			output: 2,
		},
		{
			name: "GetQueueInfo length 1",
			setupMock: func(mrd *mocks.MockRestDto, mt *mocks.MockTopics) {
				mt.EXPECT().
					GetTopicsInfo().
					Return(map[string]any{
						"test-1": []string{"test1", "test2"},
					}).Times(1)
			},
			output: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockResDto := mocks.NewMockRestDto(ctrl)
			mockTopics := mocks.NewMockTopics(ctrl)
			tt.setupMock(mockResDto, mockTopics)

			service := manager.NewManagerService(mockTopics, 10)

			resp := service.GetQueueInfo()
			if tt.output != len(resp) {
				t.Fatalf("Total key expected %d", tt.output)
			}
		})
	}

}
