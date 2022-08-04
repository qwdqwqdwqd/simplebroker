package broker

import (
	"context"

	//"fmt"
	"sync"

	//"fmt"

	"simplebroker/pkg/broker"
)

type Broker struct {
	queueManager *QueueManager
	sync.RWMutex
}

func NewBroker() broker.Broker {
	return &Broker{
		queueManager: CreateQueue(),
	}
}

func (m *Broker) Publish(ctx context.Context, queue string, msg broker.Message) (int, error) {
	q, _ := m.queueManager.GetQueue(queue)
	id := q.PublishMessage(msg)
	return id, nil
}

func (m *Broker) Subscribe(ctx context.Context, queue string) (<-chan broker.Message, error) {
	q, _ := m.queueManager.GetOrCreateQueue(queue)
	channel := q.CreateSubscriber(ctx)
	return channel, nil
}
