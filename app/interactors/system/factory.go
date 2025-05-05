package system

import (
	"dominus-project/app/domain/repos"
)

type systemService struct {
	lg         repos.LogsInt
}


type SystemServiceInt interface {
	GetLogs() ([]string, error)
}


func NewSystemService(clog repos.LogsInt) SystemServiceInt {
	return &systemService{
		lg: clog,
	}
}