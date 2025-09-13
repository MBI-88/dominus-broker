package manager

func (m *managerService) GetQueueInfo() map[string]any {
	return m.topics.GetTopicsInfo()
}
