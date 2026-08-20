package middlewares

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"dominus-broker/internal/domain/repositories"
	"dominus-broker/internal/infrastructure/enum"
	"dominus-broker/internal/infrastructure/event"
	"fmt"

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
	IdemPotency(ctx context.Context) (context.Context, error)
}

type middlewares struct {
	token []byte
	logs  event.Event
	ch    repositories.CheckerClient
}

func NewMiddleware(t string, lg event.Event, checker repositories.CheckerClient) Middleware {
	return &middlewares{
		token: []byte(t),
		logs:  lg,
		ch:    checker,
	}
}

func (m *middlewares) ApiToken(ctx context.Context) (context.Context, error) {
	ctx = m.logs.CheckID(ctx)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		m.logs.WriteLog(ctx, enum.ERROR, "middlewares.ApiToken.FromIncomingContext", enum.NOT_FOUND)
		return nil, status.Errorf(codes.DataLoss, enum.NOT_FOUND)
	}
	tokens := md.Get(enum.X_API_KEY)
	if len(tokens) == 0 {
		m.logs.WriteLog(ctx, enum.ERROR, "middlewares.ApiToken.Get", enum.NOT_FOUND)
		return nil, status.Errorf(codes.DataLoss, enum.NOT_FOUND)
	}
	hashedTokenRecived := sha256.Sum256([]byte(tokens[0]))
	hashedKey := sha256.Sum256(m.token)
	if subtle.ConstantTimeCompare(hashedTokenRecived[:], hashedKey[:]) == 0 {
		m.logs.WriteLog(ctx, enum.ERROR, "middlewares.ApiToken", enum.MATCH_TOKEN)
		return nil, status.Errorf(codes.Unauthenticated, enum.MATCH_TOKEN)
	}
	return ctx, nil
}

func (m *middlewares) UnaryLog(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	m.logs.WriteLog(ctx, enum.DEBUG, fmt.Sprintf("middlewares.%s", info.FullMethod), enum.DEBUG_DESCRIPTION)
	return handler(ctx, req)
}

func (m *middlewares) StreamLog(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	m.logs.WriteLog(ss.Context(), enum.DEBUG, fmt.Sprintf("middlewares.%s", info.FullMethod), enum.DEBUG_DESCRIPTION)
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

func (m *middlewares) IdemPotency(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	m.logs.WriteLog(ctx, enum.DEBUG, "middlewares.IdemPotency", enum.DEBUG_DESCRIPTION)

	if !ok {
		m.logs.WriteLog(ctx, enum.ERROR, "middlewares.IdemPotency.FromIncomingContext", enum.NOT_FOUND)
		return nil, status.Errorf(codes.DataLoss, enum.NOT_FOUND)
	}

	keys := md.Get(enum.IDEM_POTENCY_HEADER)
	if len(keys) == 0 {
		m.logs.WriteLog(ctx, enum.ERROR, "middlewares.IdemPotency.Get", enum.IDEM_POTENCY_NOT_FOUND)
		return nil, status.Error(codes.DataLoss, enum.IDEM_POTENCY_NOT_FOUND)
	}

	if ok := m.ch.CheckConsumer(ctx, keys[0]); ok {
		m.logs.WriteLog(ctx, enum.INFO, "middlewares.IdemPotency.CheckConsumer", enum.RATE_LIMIT_REACHED)
		return nil, status.Error(codes.Aborted, enum.RATE_LIMIT_REACHED)
	}

	if err := m.ch.SaveConsumer(ctx, keys[0]); err != nil {
		m.logs.WriteLog(ctx, enum.ERROR, "middlewares.IdemPotency.SaveConsumer", err.Error())
	}

	return ctx, nil
}
