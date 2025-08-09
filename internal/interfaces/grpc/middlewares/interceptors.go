package middlewares

import (
	"context"

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
}

func NewInterceptor(token string) Interceptor {
	return &interceptor{
		apiToken: token,
	}
}

func (i *interceptor) UnaryAuthInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", i.apiToken)
	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}

func (i *interceptor) StreamAuthInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
	streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", i.apiToken)
	s, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, err
	}
	return s, nil
}
