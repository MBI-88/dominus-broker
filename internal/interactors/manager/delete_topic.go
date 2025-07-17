package manager

import (
	"dominus-project/internal/domain/adapters"
	"fmt"
)

func (m *managerService) DeleteTopic(ctx adapters.IRestDto) error {
	name := ctx.Param("name")
	if name == "" {
		return fmt.Errorf("Param empty")
	}
	return m.topics.Delete(name)
}
