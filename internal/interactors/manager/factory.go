package manager

import (
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/domain/repos"
	"dominus-project/internal/domain/rules"
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

func NewManagerService(topics entities.TopicsInt, rls rules.RulesInt, limit int64) ManagerInt {
	return &managerService{
		topics: topics,
		rls:    rls,
		limit:  limit,
	}
}
