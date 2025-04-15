package interactors

import (
	"dominus-project/app/domain/entities"
	"dominus-project/app/domain/repos"
	"dominus-project/app/domain/rules"
	"fmt"
	"sync"
)

type managerService struct {
	topics entities.TopicsInt
	rls    rules.RulesInt
	limit  int64
}

func (m *managerService) AddTopic(ctx repos.RestContextInt) error {
	topic := &entities.Topic{
		Queue: entities.NewQueue(),
		Lck: new(sync.Mutex),
		Limit: m.limit,
	}
	
	if err := ctx.BodyParser(topic); err != nil {
		return err
	}
	if err := m.rls.ValidateStruct(topic); err != nil {
		return err
	}
	if topic.Name == "" {
		return fmt.Errorf("Topic name is empty")
	}
	if _, err := m.topics.Find(topic.Name); err != nil {
		return m.topics.Append(topic)
	}
	return fmt.Errorf("Topic exists")
}

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

func (m *managerService) DeleteTopic(ctx repos.RestContextInt) error {
	topic := new(entities.Topic)
	name := ctx.Param("name")
	if name == "" {
		return fmt.Errorf("Param empty")
	}
	topic.Name = name
	return m.topics.Delete(topic)
}

type ManagerInt interface {
	AddTopic(ctx repos.RestContextInt) error
	UpdatePartition(ctx repos.RestContextInt) error
	DeleteTopic(ctx repos.RestContextInt) error
}
