package manager

import (
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/domain/repos"
)

func (m *managerService) UpdatePartition(ctx repos.RestContextInt) error {
	topic := new(entities.Topic)
	name := ctx.Param("name")

	if err := ctx.BodyParser(topic); err != nil {
		return err
	}
	if err := m.rls.ValidateStruct(topic); err != nil {
		return err
	}

	topic.Name = name
	return m.topics.Update(topic)
}
