package repos

import "context"

type IGrpClient interface {
	Simple(url string, msg []byte) (IGrpResponse, error)
	ClientStream(urls []string, msg <-chan []byte, ctx context.Context)
	ServerStream(urls []string, initalMsg []byte, msg chan<- []byte, ctx context.Context, done chan<- struct{})
	BidirectionalStream(urls []string, provMsg <-chan []byte, subMsg chan<- []byte, errMsg chan<- error, tx chan<- struct{}, ctx context.Context, done chan<- struct{})
}