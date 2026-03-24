package event

import (
	"context"
	"dominus-project/internal/infrastructure/enum"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

type Event interface {
	WriteLog(ctx context.Context, level, op, dsc string)
	CheckID(ctx context.Context) context.Context
}

type log struct {
	lg     *slog.Logger
	mode   string // feature flag, selecting log type, cmd, client (send to other place)
	url    string // feature flag, url to send logs
	devMod bool
}

func NewEvent(mode, url string, devMode bool) Event {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError | slog.LevelDebug | slog.LevelInfo,
	})

	return &log{
		mode:   mode,
		url:    url,
		devMod: !devMode,
		lg:     slog.New(handler),
	}
}

func (l *log) WriteLog(ctx context.Context, level, op, dsc string) {
	var id string

	temp := ctx.Value(enum.ID)
	if temp == nil {
		id = ""
	} else {
		id = temp.(string)
	}
	message, err := jsoniter.Marshal(&Logs{
		ID:          id,
		Description: dsc,
		Op:          op,
	})
	if err != nil {
		return
	}

	l.typeLog(level, message)
}

func (l *log) CheckID(ctx context.Context) context.Context {
	id := ctx.Value(enum.ID)
	if id == nil {
		id = uuid.NewString()
		ctx = context.WithValue(ctx, enum.ID, id)
	}
	return ctx
}

func (l *log) typeLog(level string, message []byte) {
	switch l.mode {
	case enum.LOGCLIENT:
		l.clientLog(level, message)
	default:
		l.cmdLog(level, message)
	}
}

func (l *log) clientLog(level string, message []byte) {
}

func (l *log) cmdLog(level string, message []byte) {
	switch level {
	case enum.INFO:
		if l.devMod {
			l.lg.Info(fmt.Sprintf("LOG: %s\n", message))
		}
	case enum.DEBUG:
		if l.devMod {
			l.lg.Debug(fmt.Sprintf("LOG: %s\n", message))
		}
	case enum.WARN:
		if l.devMod {
			l.lg.Warn(fmt.Sprintf("LOG: %s\n", message))
		}
	default:
		l.lg.Error(fmt.Sprintf("LOG: %s\n", message))
	}
}
