package system


func (m *systemService) GetLogs() ([]string, error) {
	return m.lg.GetLogs()
}