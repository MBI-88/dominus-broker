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
	subsOk  = []string{"192.168.1.6:5001", "192.168.1.6:5002"}
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
	adquired = new(sync.Mutex)
)

func init() {
	rp := NewRepoMock(true, "")
	log := NewEventMock()
	inter = interactors.NewInteractor(rp, log, rls)
	inter = inter.Set(NewGrpcClientMock())
}
