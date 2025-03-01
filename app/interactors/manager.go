package interactors

import (
	"dominus/app/domain/repos"
)

type manager struct {
	lg         repos.LogsInt
}

func (m *manager) GetLogs(ctx repos.RestContextInt) ([]string, error) {
	return m.lg.GetLogs()
}

func (m *manager) GetBackup(ctx repos.RestContextInt) ([]string, error) {
	return m.lg.GetLogs()
}

type ManagerInt interface {
	GetLogs(ctx repos.RestContextInt) ([]string, error)
	GetBackup(ctx repos.RestContextInt) ([]string, error)
}
