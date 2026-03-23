package sqs

import "context"

func (q *sqs) CreateGroup() error {
	return q.client.Group(context.Background())
}