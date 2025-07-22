package manager

import (
	"dominus-project/internal/domain/adapters"
	"dominus-project/internal/domain/entities"
)

type IManager interface {
	AddTopic(ctx adapters.IRestDto) error
	UpdateSubscribers(ctx adapters.IRestDto) error
	DeleteTopic(ctx adapters.IRestDto) error
	GetQueueInfo() map[string]any
}

type managerService struct {
	topics     entities.ITopics
	queueLimit int
}

func NewManagerService(topics entities.ITopics, limit int) IManager {
	return &managerService{
		topics:     topics,
		queueLimit: limit,
	}
}
