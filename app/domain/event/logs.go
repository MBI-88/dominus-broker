package event

import (
	"log"
	"os"
)

type logs struct {
	path string
	lg   *log.Logger
}

func (l *logs) WriteLog(op, dsc string) {
	l.lg.Printf("%s: %s\n", op,dsc)
}

func (l *logs) Printf(format string, args ...any) {
	l.lg.Printf(format, args...)
}

type LogsInt interface {
	WriteLog(op, dsc string)
	Printf(format string, args ...any)
}

func NewLogs(path string) LogsInt {
	file, err := os.CreateTemp(path, "log-*.log")
	if err != nil {
		panic(err)
	}
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()

	return &logs{
		path: path,
		lg:   log.New(file, "LOG: ", log.LstdFlags),
	}
}
