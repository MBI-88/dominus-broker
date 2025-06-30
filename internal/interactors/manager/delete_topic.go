package manager

import (
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/domain/repos"
	"fmt"
)

func (m *managerService) DeleteTopic(ctx repos.RestContextInt) error {
	topic := new(entities.Topic)
	name := ctx.Param("name")
	if name == "" {
		return fmt.Errorf("Param empty")
	}
	topic.Name = name
	return m.topics.Delete(topic)
}
