package manager

import (
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/domain/repos"
)

type managerService struct {
	topics entities.ITopics
	queueLimit int
}

type IManager interface {
	AddTopic(ctx repos.IRestContext) error
	UpdateSubscribers(ctx repos.IRestContext) error
	DeleteTopic(ctx repos.IRestContext) error
	GetQueueInfo() map[string]any
}

func NewManagerService(topics entities.ITopics, limit int) IManager {
	return &managerService{
		topics: topics,
		queueLimit: limit,
	}
}
