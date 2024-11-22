package interactors

import (
	"dominus/app/domain/entities"
	"dominus/app/domain/event"
	"dominus/app/domain/rules"
)

type manager struct {
	repo RepositoryInt
	lg   event.LogsInt
	rls  rules.RulesInt
	collection string
}

func (m *manager) GetLogs(ctx RestContextInt) ([]entities.Logs, error) {
	var (
		logs []entities.Logs
		filters = ctx.Queries()
		page  = ctx.QueryInt("page")
		size = ctx.QueryInt("size")
	)

	arrfilter := m.rls.MakeLogFiter(filters, page, size)

	if err := m.repo.Filter(arrfilter, &logs, m.collection); err != nil {
		go m.lg.WriteLog("GetLogs Filter", err.Error())
		return nil, err
	}
	return logs, nil 

}

func (m *manager) GetTotalPages(ctx RestContextInt) (uint64, error) {
	total, err := m.repo.CountPages(m.collection)
	if err != nil {
		go m.lg.WriteLog("GetTotalPages CountPages", err.Error())
		return 0, err
	}
	return total, nil
}

func (m *manager) DelectLogs(ctx RestContextInt) error {
	if err := m.repo.DeleteObjects(m.rls.MakeEmptyFilter(), m.collection); err != nil {
		go m.lg.WriteLog("DelectLogs DeleteObjects", err.Error())
		return err
	}
	return nil 
}

func (m *manager) GetStats() (any, error) {
	result, err := m.repo.Stats()
	if err != nil {
		return nil, err
	}
	return result, nil 
}

func (m *manager) GetBackup(ctx RestContextInt) ([]entities.Logs, error) {
	var (
		logs []entities.Logs
		filters = ctx.Queries()
	)
	arrfilter := m.rls.MakeBackupFilter(filters)
	if err := m.repo.Filter(arrfilter, &logs, m.collection); err != nil {
		go m.lg.WriteLog("GetBackup Filter", err.Error())
		return nil, err
	}
	return logs, nil 
}


type ManagerInt interface {
	GetLogs(ctx RestContextInt) ([]entities.Logs, error)
	GetTotalPages(ctx RestContextInt) (uint64, error)
	DelectLogs(ctx RestContextInt) error
	GetStats() (any, error)
	GetBackup(ctx RestContextInt) ([]entities.Logs, error)
}