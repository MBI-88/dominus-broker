package manager

import (
	"dominus-project/internal/domain/adapters"
	"dominus-project/internal/domain/entities"
	"fmt"
)

func (m *managerService) AddTopic(ctx adapters.RestDto) error {
	topic := entities.NewTopic(ctx, m.queueLimit)

	if topic == nil {
		return fmt.Errorf("invalid topic")
	}
	if err := topic.ParseTopic(); err != nil {
		return err
	}
	if err := topic.ValidateTopic(); err != nil {
		return err
	}
	if topic.GetName() == "" {
		return fmt.Errorf("topic name is empty")
	}
	if _, err := m.topics.Find(topic.GetName()); err != nil {
		return m.topics.Append(topic)
	}
	return fmt.Errorf("topic exists")
}
