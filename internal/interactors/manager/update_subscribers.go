package manager

import (
	"dominus-project/internal/domain/adapters"
	"dominus-project/internal/domain/entities"
)

func (m *managerService) UpdateSubscribers(ctx adapters.RestDto) error {
	topic := entities.NewTopic(ctx, m.queueLimit)
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
