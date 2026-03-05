package middlewares

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"dominus-project/internal/domain/adapters"
	"dominus-project/internal/domain/enum"

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
}

type middlewares struct {
	token []byte
	logs  adapters.Logs
}

func NewMiddleware(t string, lg adapters.Logs) Middleware {
	return &middlewares{
		token: []byte(t),
		logs:  lg,
	}
}

func (m *middlewares) ApiToken(ctx context.Context) (context.Context, error) {
	ctx = m.logs.CheckID(ctx)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		go m.logs.WriteLog(ctx, enum.ERROR, "ApiToken", "Not found arguments")
		return nil, status.Errorf(codes.DataLoss, "not found arguments")
	}
	token := md.Get(enum.XapiKey)[0]
	hashedTokenRecived := sha256.Sum256([]byte(token))
	hashedKey := sha256.Sum256(m.token)
	if subtle.ConstantTimeCompare(hashedTokenRecived[:], hashedKey[:]) == 0 {
		go m.logs.WriteLog(ctx, enum.ERROR, "ApiToken", "Failed to match token")
		return nil, status.Errorf(codes.Unauthenticated, "failed to match token")
	}
	return ctx, nil
}

func (m *middlewares) UnaryLog(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	m.logs.WriteLog(ctx, enum.DEBUG, info.FullMethod, "Called")
	return handler(ctx, req)
}

func (m *middlewares) StreamLog(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	m.logs.WriteLog(ss.Context(), enum.DEBUG ,info.FullMethod, "Called")
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
