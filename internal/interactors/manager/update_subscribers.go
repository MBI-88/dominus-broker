package manager

import (
	"dominus-project/internal/domain/adapters"
	"dominus-project/internal/domain/entities"
)

func (m *managerService) UpdateSubscribers(ctx adapters.IRestDto) error {
	topic := entities.NewTopic(ctx, m.queueLimit)
	name := ctx.Param("name")
	if err := topic.FillTopic(); err != nil {
		return err
	}
	if err := topic.ValidateTopic(); err != nil {
		return err
	}
	topic.SetName(name)
	return m.topics.Update(topic)
}
