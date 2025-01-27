package interactors

import (
	"dominus/app/domain/entities"
	"dominus/app/domain/repos"
	"dominus/app/domain/rules"
)

type manager struct {
	repo       repos.RepositoryInt
	lg         LogsInt
	rls        rules.RulesInt
	collection string
}

func (m *manager) GetLogs(ctx repos.RestContextInt) ([]entities.Logs, error) {
	var (
		logs    = make([]entities.Logs, 0, 100)
		filters = ctx.Queries()
		page    = ctx.QueryInt("page")
		size    = ctx.QueryInt("size")
	)
	if err := m.repo.Filter(filters, &logs, m.collection, page, size); err != nil {
		go m.lg.WriteLog("GetLogs Filter", err.Error())
		return nil, err
	}
	return logs, nil

}

func (m *manager) GetTotalPages(ctx repos.RestContextInt) (int64, error) {
	total, err := m.repo.CountPages(m.collection)
	if err != nil {
		go m.lg.WriteLog("GetTotalPages CountPages", err.Error())
		return 0, err
	}
	return total, nil
}

func (m *manager) DelectLogs(ctx repos.RestContextInt) error {
	if err := m.repo.DeleteObjects(m.collection); err != nil {
		go m.lg.WriteLog("DelectLogs DeleteObjects", err.Error())
		return err
	}
	return nil
}

func (m *manager) GetStats() (any, error) {
	result, err := m.repo.Stats()
	if err != nil {
		go m.lg.WriteLog("GetStats Stats", err.Error())
		return nil, err
	}
	return result, nil
}

func (m *manager) GetBackup(ctx repos.RestContextInt) ([]entities.Logs, error) {
	var (
		logs    = make([]entities.Logs, 0, 1000)
		filters = ctx.Queries()
	)
	if err := m.repo.Backup(filters, &logs, m.collection); err != nil {
		go m.lg.WriteLog("GetBackup Filter", err.Error())
		return nil, err
	}
	return logs, nil
}

type ManagerInt interface {
	GetLogs(ctx repos.RestContextInt) ([]entities.Logs, error)
	GetTotalPages(ctx repos.RestContextInt) (int64, error)
	DelectLogs(ctx repos.RestContextInt) error
	GetStats() (any, error)
	GetBackup(ctx repos.RestContextInt) ([]entities.Logs, error)
}
