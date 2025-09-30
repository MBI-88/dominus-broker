package manager

import (
	"dominus-project/internal/domain/adapters"
	"dominus-project/internal/domain/entities"
	"fmt"
)

func (m *managerService) UpdateSubscribers(ctx adapters.RestDto) error {
	topic := entities.NewTopic(ctx, m.queueLimit)

	if topic == nil {
		return fmt.Errorf("invalid topic")
	}
	name := ctx.Param("name")
	if err := topic.ParseTopic(); err != nil {
		return err
	}
	topic.SetName(name)
	if err := topic.ValidateTopic(); err != nil {
		return err
	}
	return m.topics.Update(topic)
}
