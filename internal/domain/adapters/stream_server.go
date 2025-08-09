package adapters

import "context"

type StreamServer interface {
	Send(payload []byte) error
	Context() context.Context
}
