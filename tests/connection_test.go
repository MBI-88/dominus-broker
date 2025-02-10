package tests

import (
	"encoding/json"
	"testing"
)

func TestSimpleConn(t *testing.T) {
	t.Run("Connection-OK", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		simple := inter.NewConnection()
		result := simple.SimpleConn(NewSimpleConnMock(payload, subsOk))
		if result != nil {
			t.Fatal(result)
		}
	})

	t.Run("Connection-Error", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		simple := inter.NewConnection()
		result := simple.SimpleConn(NewSimpleConnMock(payload, make([]string, 0)))
		if result == nil {
			t.Fatalf("Expected error but received nil")
		}
	})
}

func TestCheckURI(t *testing.T) {
	arrayErr := make([]bool, 0, 2)
	t.Run("Check-OK", func(t *testing.T) {
		for _, uri := range subsOk {
			arrayErr = append(arrayErr, rls.CheckURI(uri))
		}

		if arrayErr[0] != true && arrayErr[1] != true {
			t.Fatalf("Expected true but received false")
		}
	})

	t.Run("Check-Error", func(t *testing.T) {
		arrayErr = make([]bool, 0, 2)
		for _, uri := range subsEr {
			arrayErr = append(arrayErr, rls.CheckURI(uri))
		}

		if arrayErr[0] != false && arrayErr[1] != false {
			t.Fatalf("Expected false but received true")
		}
	})
}

func TestClientStream(t *testing.T) {
	t.Run("ClientStream-OK", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		stream := inter.NewConnection()
		result := stream.StreamClientConn(NewClienConnMock(payload, subsOk))
		if result.Error() != "End" {
			t.Fatal(result)
		}
	})

	t.Run("ClientStream-Error", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		stream := inter.NewConnection()
		result := stream.StreamClientConn(NewClienConnMock(payload, make([]string, 0)))
		if result == nil {
			t.Fatalf("Expected error but received nil")
		}
	})
}

func TestServerStream(t *testing.T) {
	t.Run("ServerStream-OK", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		stream := inter.NewConnection()
		msg, ctx := NewServerConnMock(payload, subsOk)
		result := stream.StreamServerConn(msg, ctx)
		if result.Error() != "Connection closed" {
			t.Fatal(result)
		}

	})

	t.Run("ServerStream-Error", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		stream := inter.NewConnection()
		msg, ctx := NewServerConnMock(payload, make([]string, 0))
		result := stream.StreamServerConn(msg, ctx)
		if result.Error() == "End" {
			t.Fatalf("Expected error but received nil")
		}
	})

}

func TestBidirectionalStream(t *testing.T) {
	t.Run("BidirectionalStream-OK", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		stream := inter.NewConnection()
		result := stream.StreamBiConn(NewBiConnMock(payload, subsOk))
		if result.Error() != "Connection closed" {
			t.Fatal(result)
		}
	})

	t.Run("BidirectionalStream-Error", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		stream := inter.NewConnection()
		result := stream.StreamBiConn(NewBiConnMock(payload, make([]string, 0)))
		if result.Error() == "End" {
			t.Fatalf("Expected error but received nil")
		}
	})
}
