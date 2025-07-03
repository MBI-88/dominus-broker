package repos

import "context"

type StreamServerInt interface {
	Send(payload []byte) error
	Context() context.Context
}
