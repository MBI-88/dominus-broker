package manager

import (
	"dominus-project/app/domain/entities"
	"dominus-project/app/domain/repos"
	"dominus-project/app/domain/rules"
)

type managerService struct {
	topics entities.TopicsInt
	rls    rules.RulesInt
	limit  int64
}



type ManagerInt interface {
	AddTopic(ctx repos.RestContextInt) error
	UpdatePartition(ctx repos.RestContextInt) error
	DeleteTopic(ctx repos.RestContextInt) error
}

func NewManagerService( topics entities.TopicsInt, rls rules.RulesInt, limit int64) ManagerInt {
	return &managerService{
		topics: topics,
		rls:    rls,
		limit:  limit,
	}
}