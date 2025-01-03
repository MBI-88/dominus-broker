package tests

import (
	"encoding/json"
	"testing"
)

func TestSimpleConnOK(t *testing.T) {
	payload, _ := json.Marshal(data)
	simple := inter.NewConnection()
	result := simple.SimpleConn(NewSimpleConnMock(payload, subsOk))
	if result != nil {
		t.Fatal(result)
	}
}

func TestSimpleConnER(t *testing.T) {
	payload, _ := json.Marshal(data)
	simple := inter.NewConnection()
	result := simple.SimpleConn(NewSimpleConnMock(payload, make([]string, 0)))
	if result == nil {
		t.Fatalf("Expected error but received nil")
	}
}

func TestCheckURI(t *testing.T) {
	arrayErr := make([]bool, 0, 2)
	for _, uri := range subsEr {
		arrayErr = append(arrayErr, rls.CheckURI(uri))
	}

	if arrayErr[0] != false && arrayErr[1] != false {
		t.Fatalf("Expected false but received true")
	} 

	arrayErr = make([]bool, 0, 2)
	for _, uri := range subsOk {
		arrayErr = append(arrayErr, rls.CheckURI(uri))
	}

	if arrayErr[0] != true && arrayErr[1] != true {
		t.Fatalf("Expected true but received false")
	}
}

func TestClientStreamOK(t *testing.T) {
	payload, _ := json.Marshal(data)
	stream := inter.NewConnection()
	result := stream.StreamClientConn(NewClienConnMock(payload, subsOk))
	if result.Error() != "End"  {
		t.Fatal(result)
	}
}

func TestClientStreamER(t *testing.T) {
	payload, _ := json.Marshal(data)
	stream := inter.NewConnection()
	result := stream.StreamClientConn(NewClienConnMock(payload, make([]string, 0)))
	if result == nil {
		t.Fatalf("Expected error but received nil")
	}
}

func TestServerStreamOK(t *testing.T) {
	payload, _ := json.Marshal(data)
	stream := inter.NewConnection()
	msg, ctx := NewServerConnMock(payload, subsOk)
	result := stream.StreamServerConn(msg, ctx)
	if result.Error() != "Connection closed"  {
		t.Fatal(result)
	}
}

func TestServerStreamER(t *testing.T) {
	payload, _ := json.Marshal(data)
	stream := inter.NewConnection()
	msg, ctx := NewServerConnMock(payload, make([]string, 0))
	result := stream.StreamServerConn(msg, ctx)
	if result.Error() == "End"  {
		t.Fatalf("Expected error but received nil")
	}
}

func TestBidirectionalStreamOK(t *testing.T) {
	payload, _ := json.Marshal(data)
	stream := inter.NewConnection()
	result := stream.StreamBiConn(NewBiConnMock(payload, subsOk))
	if result.Error() != "Connection closed"  {
		t.Fatal(result)
	}
}

func TestBidirectionalStreamER(t *testing.T) {
	payload, _ := json.Marshal(data)
	stream := inter.NewConnection()
	result := stream.StreamBiConn(NewBiConnMock(payload, make([]string, 0)))
	if result.Error() == "End"  {
		t.Fatalf("Expected error but received nil")
	}
}