package manager

import (
	"dominus-project/internal/domain/repos"
	"fmt"
)

func (m *managerService) DeleteTopic(ctx repos.RestContextInt) error {
	name := ctx.Param("name")
	if name == "" {
		return fmt.Errorf("Param empty")
	}
	return m.topics.Delete(name)
}
