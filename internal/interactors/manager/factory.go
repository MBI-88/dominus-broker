package manager

import (
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/domain/repos"
	"dominus-project/internal/domain/rules"
)

type managerService struct {
	topics entities.TopicsInt
	rls    rules.RulesInt
}

type ManagerInt interface {
	AddTopic(ctx repos.RestContextInt) error
	UpdateSubscribers(ctx repos.RestContextInt) error
	DeleteTopic(ctx repos.RestContextInt) error
	GetQueueInfo() map[string]any
}

func NewManagerService(topics entities.TopicsInt, rls rules.RulesInt) ManagerInt {
	return &managerService{
		topics: topics,
		rls:    rls,
	}
}
