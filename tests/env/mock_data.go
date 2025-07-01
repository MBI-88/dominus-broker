package env

import (
	"sync"
	"time"
)

var (
	SubsOk  = []string{"78.168.1.6:5001/api/test", "grpc.dominus.com/api","192.16.1.6:5001","grpc.dominus.com"}
	Tps = "test.example"
	Data  = struct {
		Id        string    `json:"id"`
		Name      string    `json:"name"`
		Dni       string    `json:"dni"`
		Age       int64     `json:"age"`
		CreatedAt time.Time `json:"created_at"`
	}{
		Id:        "3131f3aferewr13ewr466erw65",
		Name:      "Connection",
		Dni:       "4644646AB464",
		Age:       35,
		CreatedAt: time.Now(),
	}
	Ad = new(sync.Mutex)
)
