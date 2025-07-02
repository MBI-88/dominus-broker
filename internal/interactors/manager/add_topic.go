package manager

import (
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/domain/repos"
	"fmt"
	"sync"
)

func (m *managerService) AddTopic(ctx repos.RestContextInt) error {
	topic := &entities.Topic{
		Lck:   new(sync.Mutex),
	}

	if err := ctx.BodyParser(topic); err != nil {
		return err
	}
	if err := m.rls.ValidateStruct(topic); err != nil {
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
