package interactors

import (
	"dominus-project/app/domain/repos"
)

type systemService struct {
	lg         repos.LogsInt
}

func (m *systemService) GetLogs() ([]string, error) {
	return m.lg.GetLogs()
}

type SystemServiceInt interface {
	GetLogs() ([]string, error)
}
