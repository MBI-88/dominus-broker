package manager

import (
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/domain/repos"
)

func (m *managerService) UpdateSubscribers(ctx repos.IRestContext) error {
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
