package system

import (
	"dominus-project/internal/domain/adapters"
)

type systemService struct {
	lg adapters.ILogs
}

type SystemServiceInt interface {
	GetLogs() ([]string, error)
}

func NewSystemService(clog adapters.ILogs) SystemServiceInt {
	return &systemService{
		lg: clog,
	}
}
