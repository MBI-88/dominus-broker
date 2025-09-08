package grpcconn_test

import (
	"dominus-project/internal/domain/entities"
	grpcconn "dominus-project/internal/interactors/grpc_conn"
	"dominus-project/mocks"
	"dominus-project/tests/env"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

func TestSimpleConn(t *testing.T) {
	tests := []struct {
		name string 
		setupMock func (*mocks.MockGrpcDto, *mocks.MockTopics, *mocks.MockTopic)
		output error
	}{
		{
			name: "SimpleConn Ok",
			setupMock: func(mgd *mocks.MockGrpcDto, mt1 *mocks.MockTopics, mt2 *mocks.MockTopic) {
				mgd.EXPECT(). 
				GetSubscribers(). 
				Return([]string{"test"}).Times(1)

				mt1.EXPECT(). 
				Find(gomock.All()). 
				Return(mt2, nil).Times(1)

				mt2.EXPECT(). 
				GetSubscribers(). 
				Return([]string{"test"}).Times(1)

				mgd.EXPECT(). 
				GetPayload(). 
				Return([]byte("test-1")).Times(1)

				mt2.EXPECT(). 
				SetMessage(gomock.All()). 
				Return(nil).Times(1)
			},
			output: nil,
		},
		{
			name: "SimpleConn subscribers empty",
			setupMock: func(mgd *mocks.MockGrpcDto, mt1 *mocks.MockTopics, mt2 *mocks.MockTopic) {
				mgd.EXPECT(). 
				GetSubscribers(). 
				Return([]string{}).Times(1)
			},
			output: fmt.Errorf("subscribers empty"),
		},
		{
			name: "SimpleConn find topic error",
			setupMock: func(mgd *mocks.MockGrpcDto, mt1 *mocks.MockTopics, mt2 *mocks.MockTopic) {
				mgd.EXPECT(). 
				GetSubscribers(). 
				Return([]string{"test"}).Times(1)

				mt1.EXPECT(). 
				Find(gomock.All()). 
				Return(mt2, fmt.Errorf("find topic error")).Times(1)

			},
			output: fmt.Errorf("find topic errror"),
		},
		{
			name: "SimpleConn subscribers topic added empty",
			setupMock: func(mgd *mocks.MockGrpcDto, mt1 *mocks.MockTopics, mt2 *mocks.MockTopic) {
				mgd.EXPECT(). 
				GetSubscribers(). 
				Return([]string{"test"}).Times(1)

				mt1.EXPECT(). 
				Find(gomock.All()). 
				Return(mt2, nil).Times(1)

				mt2.EXPECT(). 
				GetSubscribers(). 
				Return([]string{}).Times(1)
			},
			output: fmt.Errorf("subscribers topic added empty"),
		},
		{
			name: "SimpleConn payload empty",
			setupMock: func(mgd *mocks.MockGrpcDto, mt1 *mocks.MockTopics, mt2 *mocks.MockTopic) {
				mgd.EXPECT(). 
				GetSubscribers(). 
				Return([]string{"test"}).Times(1)

				mt1.EXPECT(). 
				Find(gomock.All()). 
				Return(mt2, nil).Times(1)

				mt2.EXPECT(). 
				GetSubscribers(). 
				Return([]string{"test"}).Times(1)

				mgd.EXPECT(). 
				GetPayload(). 
				Return([]byte("")).Times(1)
			},
			output: fmt.Errorf("payload empty"),
		},
		{
			name: "SimpleConn set message error",
			setupMock: func(mgd *mocks.MockGrpcDto, mt1 *mocks.MockTopics, mt2 *mocks.MockTopic) {
				mgd.EXPECT(). 
				GetSubscribers(). 
				Return([]string{"test"}).Times(1)

				mt1.EXPECT(). 
				Find(gomock.All()). 
				Return(mt2, nil).Times(1)

				mt2.EXPECT(). 
				GetSubscribers(). 
				Return([]string{"test"}).Times(1)

				mgd.EXPECT(). 
				GetPayload(). 
				Return([]byte("test-1")).Times(1)

				mt2.EXPECT(). 
				SetMessage(gomock.All()). 
				Return(fmt.Errorf("message error")).Times(1)
			},
			output: fmt.Errorf("set message error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish() 

			mockGrpcDto := mocks.NewMockGrpcDto(ctrl)
			mockTopics := mocks.NewMockTopics(ctrl)
			mockLogs := mocks.NewMockLogs(ctrl)
			mockGrpcClient := mocks.NewMockGrpcClient(ctrl)
			mockTopic := mocks.NewMockTopic(ctrl)

			tt.setupMock(mockGrpcDto, mockTopics, mockTopic) 
			service := grpcconn.NewGrpcService(mockLogs, mockGrpcClient, mockTopics)

			err := service.SimpleConn(mockGrpcDto)
			if tt.output != nil && err == nil {
				t.Fatalf("Expected output different from output %s != %s", tt.output, err)
			}
			if tt.output == nil && err != nil {
				t.Fatalf("Expected output different from output %s != %s", tt.output, err)
			}
		})
	}
}



func TestRunQueue(t *testing.T) {
	log := env.NewEventMock(true)
	topics := entities.NewTopics(100)
	topic := entities.NewTopic(nil, env.QueueLimit)
	topic.SetName(env.Tps)
	topic.SetSubscribers(env.SubsOk)
	body, _ := json.Marshal(env.Data)
	topic.SetMessage(body)
	topics.Append(topic)

	t.Run("Find-Topic-OK", func(t *testing.T) {
		simple := grpcconn.NewGrpcService(log, env.NewGrpcClientMock(true, false), topics)
		queue := make(chan struct{})
		go func() {
			time.Sleep(3 * time.Second)
			queue <- struct{}{}
		}()

		if err := simple.RunQueue(queue); err.Error() != "Queue closed" {
			t.Fatal("Queue error")
		}

		close(queue)
	})

	t.Run("Find-Topic-Connection-Error", func(t *testing.T) {
		simple := grpcconn.NewGrpcService(log, env.NewGrpcClientMock(false, false), topics)
		queue := make(chan struct{})
		go func() {
			time.Sleep(3 * time.Second)
			queue <- struct{}{}
		}()

		if err := simple.RunQueue(queue); err.Error() != "Queue closed" {
			t.Fatal("Queue error")
		}

		close(queue)
	})
	
	t.Run("Find-Topic-Response-Error", func(t *testing.T) {
		simple := grpcconn.NewGrpcService(log, env.NewGrpcClientMock(true, false), topics)
		queue := make(chan struct{})
		go func() {
			time.Sleep(3 * time.Second)
			queue <- struct{}{}
		}()

		if err := simple.RunQueue(queue); err.Error() != "Queue closed" {
			t.Fatal("Queue error")
		}

		close(queue)
	})

}
