package interactors

import (
	"dominus-project/app/domain/entities"
	"dominus-project/app/domain/repos"
	"dominus-project/app/domain/rules"
)

type managerService struct {
	topics entities.TopicsInt
	rls    rules.RulesInt
}

func (m *managerService) AddTopic(ctx repos.RestContextInt) error {
	topic := new(entities.Topic)

	if err := ctx.BodyParser(topic); err != nil {
		return err
	}
	if err := m.rls.ValidateStruct(topic); err != nil {
		return err
	}

	m.topics.Append(topic)
	return nil
}

func (m *managerService) UpdatePartition(ctx repos.RestContextInt) error {
	topic := new(entities.Topic)
	name := ctx.Params("name")

	if err := ctx.BodyParser(topic); err != nil {
		return err
	}
	if err := m.rls.ValidateStruct(topic); err != nil {
		return err
	}
	
	topic.Name = name
	return m.topics.Update(topic)
}

func (m *managerService) DeleteTopic(ctx repos.RestContextInt) error {
	topic := new(entities.Topic)
	name := ctx.Params("name")
	topic.Name = name
	return m.topics.Delete(topic)
}

type ManagerInt interface {
	AddTopic(ctx repos.RestContextInt) error
	UpdatePartition(ctx repos.RestContextInt) error
	DeleteTopic(ctx repos.RestContextInt) error
}
