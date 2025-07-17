package adapters

import "context"

type IStreamServer interface {
	Send(payload []byte) error
	Context() context.Context
}
