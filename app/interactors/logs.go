package interactors

import (
	"github.com/PR0C0D3-MBI/dominus-project/app/domain/repos"
)

type systemService struct {
	lg         repos.LogsInt
}

func (m *systemService) GetLogs(ctx repos.RestContextInt) ([]string, error) {
	return m.lg.GetLogs()
}

func (m *systemService) GetBackup(ctx repos.RestContextInt) ([]string, error) {
	return m.lg.GetLogs()
}

type SystemServiceInt interface {
	GetLogs(ctx repos.RestContextInt) ([]string, error)
	GetBackup(ctx repos.RestContextInt) ([]string, error)
}
