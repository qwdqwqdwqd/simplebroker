package broker

import (
	"context"
)

type Message struct {
	id   int
	Body string
}

type Broker interface {
	Publish(ctx context.Context, queue string, msg Message) (int, error)
	Subscribe(ctx context.Context, queue string) (<-chan Message, error)
}
