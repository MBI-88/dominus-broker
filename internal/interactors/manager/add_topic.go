package manager

import (
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/domain/repos"
	"fmt"
)

func (m *managerService) AddTopic(ctx repos.RestContextInt) error {
	topic := entities.NewTopic(ctx)

	if err := topic.FillTopic(); err != nil {
		return err
	}
	if err := topic.ValidateTopic(); err != nil {
		return err
	}
	if topic.GetName() == "" {
		return fmt.Errorf("Topic name is empty")
	}
	if _, err := m.topics.Find(topic.GetName()); err != nil {
		return m.topics.Append(topic)
	}
	return fmt.Errorf("Topic exists")
}
