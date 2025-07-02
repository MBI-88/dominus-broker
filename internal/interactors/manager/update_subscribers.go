package manager

import (
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/domain/repos"
)

func (m *managerService) UpdateSubscribers(ctx repos.RestContextInt) error {
	topic := new(entities.Topic)
	name := ctx.Param("name")

	if err := ctx.BodyParser(topic); err != nil {
		return err
	}
	if err := m.rls.ValidateStruct(topic); err != nil {
		return err
	}

	topic.SetName(name)
	return m.topics.Update(topic)
}
