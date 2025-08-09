package adapters

type Logs interface {
	WriteLog(op, dsc string)
	Printf(format string, args ...any)
	GetLogs() ([]string, error)
}
