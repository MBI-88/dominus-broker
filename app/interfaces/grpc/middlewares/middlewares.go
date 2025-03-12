package middlewares

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"github.com/PR0C0D3-MBI/dominus-project/app/domain/repos"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type middlewares struct {
	token []byte
	logs  repos.LogsInt
}

func (m *middlewares) ApiToken(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "Not found arguments")
	}
	token := md.Get("API_TOKEN")[0]
	hashedTokenRecived := sha256.Sum256([]byte(token))
	hashedKey := sha256.Sum256(m.token)
	if subtle.ConstantTimeCompare(hashedTokenRecived[:], hashedKey[:]) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Failed to match token")
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

func NewMiddleware(t string, lg repos.LogsInt) MiddlewareInt {
	return &middlewares{
		token: []byte(t),
		logs:  lg,
	}
}
