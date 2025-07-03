package manager

import (
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/domain/repos"
)

type managerService struct {
	topics entities.TopicsInt
}

type ManagerInt interface {
	AddTopic(ctx repos.RestContextInt) error
	UpdateSubscribers(ctx repos.RestContextInt) error
	DeleteTopic(ctx repos.RestContextInt) error
	GetQueueInfo() map[string]any
}

func NewManagerService(topics entities.TopicsInt) ManagerInt {
	return &managerService{
		topics: topics,
	}
}
