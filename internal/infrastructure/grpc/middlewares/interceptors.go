package middlewares

import (
	"context"
	"dominus-broker/internal/infrastructure/enum"
	"dominus-broker/internal/infrastructure/event"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Interceptor interface {
	StreamAuthInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
		streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error)
	UnaryAuthInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error
}

type interceptor struct {
	apiToken string
	log      event.Event
}

func NewInterceptor(token string, log event.Event) Interceptor {
	return &interceptor{
		apiToken: token,
		log:      log,
	}
}

func (i *interceptor) UnaryAuthInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ctx = metadata.AppendToOutgoingContext(ctx, enum.X_API_KEY, i.apiToken)
	go i.log.WriteLog(ctx, enum.DEBUG, "UnaryAuthInterceptor", enum.DEBUG_DESCRIPTION)
	return invoker(ctx, method, req, reply, cc, opts...)
}

func (i *interceptor) StreamAuthInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
	streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, enum.X_API_KEY, i.apiToken)
	s, err := streamer(ctx, desc, cc, method, opts...)
	go i.log.WriteLog(ctx, enum.DEBUG, "StreamAuthInterceptor", enum.DEBUG_DESCRIPTION)
	if err != nil {
		go i.log.WriteLog(ctx, enum.ERROR, "StreamAuthInterceptor", err.Error())
		return nil, err
	}
	return s, nil
}
