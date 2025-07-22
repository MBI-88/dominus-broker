package system

import (
	"dominus-project/internal/domain/adapters"
)

type ISystemService interface {
	GetLogs() ([]string, error)
}

type systemService struct {
	lg adapters.ILogs
}

func NewSystemService(clog adapters.ILogs) ISystemService {
	return &systemService{
		lg: clog,
	}
}
