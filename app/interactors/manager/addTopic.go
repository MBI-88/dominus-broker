package manager

import (
	"dominus-project/app/domain/entities"
	"dominus-project/app/domain/repos"
	"fmt"
	"sync"
)

func (m *managerService) AddTopic(ctx repos.RestContextInt) error {
	topic := &entities.Topic{
		Queue: entities.NewQueue(),
		Lck:   new(sync.Mutex),
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