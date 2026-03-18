package middlewares

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"dominus-project/internal/adapters/event"
	"dominus-project/internal/adapters/enum"
	"dominus-project/internal/domain/repositories"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Middleware interface {
	ApiToken(ctx context.Context) (context.Context, error)
	LogErrors() logging.Logger
	UnaryLog(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)
	StreamLog(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
	IdPotency(ctx context.Context) (context.Context, error)
}

type middlewares struct {
	token []byte
	logs  event.Event
	ch   repositories.CheckerClient
}

func NewMiddleware(t string, lg event.Event, checker repositories.CheckerClient) Middleware {
	return &middlewares{
		token: []byte(t),
		logs:  lg,
		ch: checker,
	}
}

func (m *middlewares) ApiToken(ctx context.Context) (context.Context, error) {
	ctx = m.logs.CheckID(ctx)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		go m.logs.WriteLog(ctx, enum.ERROR, "ApiToken", enum.NOT_FOUND)
		return nil, status.Errorf(codes.DataLoss, enum.NOT_FOUND)
	}
	token := md.Get(enum.X_API_KEY)[0]
	hashedTokenRecived := sha256.Sum256([]byte(token))
	hashedKey := sha256.Sum256(m.token)
	if subtle.ConstantTimeCompare(hashedTokenRecived[:], hashedKey[:]) == 0 {
		go m.logs.WriteLog(ctx, enum.ERROR, "ApiToken", enum.MATCH_TOKEN)
		return nil, status.Errorf(codes.Unauthenticated, enum.MATCH_TOKEN)
	}
	return ctx, nil
}

func (m *middlewares) UnaryLog(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	go m.logs.WriteLog(ctx, enum.DEBUG, info.FullMethod, enum.DEBUG_DESCRIPTION)
	return handler(ctx, req)
}

func (m *middlewares) StreamLog(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	go m.logs.WriteLog(ss.Context(), enum.DEBUG ,info.FullMethod, enum.DEBUG_DESCRIPTION)
	return handler(srv, ss)
}

func (m *middlewares) LogErrors() logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		var l string
		switch level {
			case logging.LevelDebug:
			l = enum.DEBUG
			case logging.LevelInfo:
			l = enum.INFO
			case logging.LevelWarn:
			l = enum.WARN
			default:
			l = enum.ERROR

		}
		go m.logs.WriteLog(ctx, l, "LogErrors", msg)
	})
}

func (m *middlewares) IdPotency(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	go m.logs.WriteLog(ctx, enum.DEBUG, "IdPotency", enum.DEBUG_DESCRIPTION)

	if !ok {
		go m.logs.WriteLog(ctx, enum.ERROR, "IdPotency", enum.NOT_FOUND)
		return nil, status.Errorf(codes.DataLoss, enum.NOT_FOUND)
	}

	key := md.Get(enum.ID_POTENCY_HEADER)[0]
	if key == "" {
		go m.logs.WriteLog(ctx, enum.ERROR, "IdPotency", enum.NOT_FOUND)
		return nil, status.Error(codes.DataLoss, enum.NOT_FOUND)
	}

	if ok := m.ch.CheckConsumer(ctx, key); ok {
		go m.logs.WriteLog(ctx, enum.INFO, "CheckConsumer", "id potency found")
		return nil, status.Error(codes.Aborted, "id potency found")
	}

	go func(key string) {
		if err := m.ch.SaveConsumer(ctx, key); err != nil {
			m.logs.WriteLog(ctx, enum.ERROR, "SaveConsumer", err.Error())
		}
	}(key)

	return ctx, nil
}