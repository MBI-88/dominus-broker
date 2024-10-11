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
	//WriteLog. Writes logs created by database and system operations
	//
	//Parameters
	//
	//-> op: type of operation 
	//
	//-> dsc: description of error
	WriteLog(op, dsc string)
	//Used for fasthttp
	//
	//Parameters
	//
	//-> format: type of format chosen by fasthttp
	//
	//-> args: any countity of argumentes chosen by fasthttp
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
