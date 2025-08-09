package events

import (
	"dominus-project/internal/domain/adapters"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type logs struct {
	path  string
	lg    *log.Logger
	locck *sync.Mutex
}

func NewLogs(path string) adapters.Logs {
	file, err := os.CreateTemp(path, "log-*.log")
	if err != nil {
		panic(err)
	}
	return &logs{
		path:  path,
		lg:    log.New(file, "LOG: ", log.LstdFlags),
		locck: new(sync.Mutex),
	}
}

func (l *logs) WriteLog(op, dsc string) {
	l.locck.Lock()
	defer l.locck.Unlock()
	l.lg.Printf("%s: %s\n", op, dsc)
}

func (l *logs) Printf(format string, args ...any) {
	l.lg.Printf(format, args...)
}

func (l *logs) GetLogs() ([]string, error) {
	files := make([]string, 0, 100)
	payloads := make([]string, 0, 1000)
	err := filepath.Walk(l.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, info.Name())
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	for _, val := range files {
		if val != "" {
			p := fmt.Sprintf("%s/%s", l.path, val)
			f, err := os.Open(p)
			if err != nil {
				return nil, err
			}
			body, err := io.ReadAll(f)
			if err != nil {
				return nil, err
			}
			payloads = append(payloads, string(body))
			f.Close()
		}
	}
	return payloads, nil
}
