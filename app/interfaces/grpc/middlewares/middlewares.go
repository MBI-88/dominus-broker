package middlewares

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"dominus/app/domain/event"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type middlewares struct {
	token []byte
	logs  event.LogsInt
}

func (m *middlewares) ApiToken(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, status.Errorf(codes.DataLoss, "Not found arguments")
	}
	if t, ok := md["API_TOKEN"]; ok {
		if len(t) != 1 {
			return nil, status.Errorf(codes.InvalidArgument, "Format error")
		}
		hashedTokenRecived := sha256.Sum256([]byte(t[0]))
		hashedKey := sha256.Sum256(m.token)

		if subtle.ConstantTimeCompare(hashedTokenRecived[:], hashedKey[:]) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "Failed to match token")
		}
	}
	return ctx, nil
}

func (m *middlewares) UnaryLog(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	m.logs.Printf(info.FullMethod, "called")
	return handler(ctx, req)
}

func (m *middlewares) StreamLog(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	m.logs.Printf(info.FullMethod, "called")
	return handler(srv, ss)
}

func (m *middlewares) LogErrors() logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, level logging.Level, msg string, fields ...any) {
		if level == logging.LevelError {
			m.logs.WriteLog("Error", msg)
		}
	})
}



type MiddlewareInt interface {
	ApiToken(ctx context.Context) (context.Context, error)
	LogErrors() logging.Logger
	UnaryLog(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)
	StreamLog(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
}

func NewMiddleware(t string, lg event.LogsInt) MiddlewareInt {
	return &middlewares{
		token: []byte(t),
		logs:  lg,
	}
}
