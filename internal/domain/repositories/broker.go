package repositories

import "context"

// Broker outputs
type BrokerClient interface {
	ClientStream(urls []string, msg <-chan []byte, ctx context.Context)
	ServerStream(urls []string, initalMsg []byte, msg chan<- []byte, ctx context.Context, done chan<- struct{})
	BidirectionalStream(urls []string, provMsg <-chan []byte, subMsg chan<- []byte, tx chan<- struct{}, ctx context.Context)
}
