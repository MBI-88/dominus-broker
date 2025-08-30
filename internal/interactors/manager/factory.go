package manager

import (
	"dominus-project/internal/domain/adapters"
	"dominus-project/internal/domain/entities"
)

type ManagerService interface {
	AddTopic(ctx adapters.RestDto) error
	UpdateSubscribers(ctx adapters.RestDto) error
	DeleteTopic(ctx adapters.RestDto) error
	GetQueueInfo() map[string]any
}

type managerService struct {
	topics     entities.Topics
	queueLimit int
}

func NewManagerService(topics entities.Topics, limit int) ManagerService {
	return &managerService{
		topics:     topics,
		queueLimit: limit,
	}
}
