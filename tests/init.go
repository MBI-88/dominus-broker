package tests

import (
	"dominus/app/domain/rules"
	"dominus/app/interactors"
	"sync"
	"time"
)

var (
	rls = rules.NewRules()
	inter interactors.InteractorInt
	subsOk  = []string{"78.168.1.6:5001/api/test", "grpc.dominus.com/api","192.16.1.6:5001","grpc.dominus.com"}
	subsEr = []string{"error.err", "error.empty:5001"}
	data  = struct {
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
	ad = new(sync.Mutex)
)

func init() {
	log := NewEventMock()
	inter = interactors.NewInteractor(log, rls)
	inter = inter.Set(NewGrpcClientMock())
}
