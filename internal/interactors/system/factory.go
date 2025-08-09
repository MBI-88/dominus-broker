package system

import (
	"dominus-project/internal/domain/adapters"
)

type SystemService interface {
	GetLogs() ([]string, error)
}

type systemService struct {
	lg adapters.Logs
}

func NewSystemService(clog adapters.Logs) SystemService {
	return &systemService{
		lg: clog,
	}
}
