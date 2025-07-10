package system

import (
	"dominus-project/internal/domain/repos"
)

type systemService struct {
	lg repos.ILogs
}

type SystemServiceInt interface {
	GetLogs() ([]string, error)
}

func NewSystemService(clog repos.ILogs) SystemServiceInt {
	return &systemService{
		lg: clog,
	}
}
