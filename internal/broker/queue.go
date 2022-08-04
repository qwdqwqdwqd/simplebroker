package broker

import (
	"context"

	"simplebroker/pkg/broker"
	"sync"
)

type Queue struct {
	Name          string
	Subscribers   map[int]*Subscriber
	subDeleteChan chan *Subscriber
	subAddChan    chan *Subscriber
	msgPubChan    chan *broker.Message
	sync.Mutex
}

func (t *Queue) CreateSubscriber(ctx context.Context) chan broker.Message {
	ch := make(chan broker.Message)
	newSub := CreateNewSubscriber(ctx, ch, t.subDeleteChan)
	t.subAddChan <- newSub
	return ch
}

func (t *Queue) PublishMessage(msg broker.Message) int {
	id := 0 // TODO: id generator
	t.msgPubChan <- &msg
	return id
}
func (t *Queue) Listener() {
	for {
		select {
		case msg := <-t.msgPubChan:
			var wg sync.WaitGroup
			for _, sub := range t.Subscribers {
				sub := sub
				wg.Add(1)
				go func() {
					sub.subChanel <- msg
					wg.Done()
				}()
			}
			wg.Wait()
		case newSub := <-t.subAddChan:
			t.Subscribers[newSub.Id] = newSub
		case subscriber := <-t.subDeleteChan:
			delete(t.Subscribers, subscriber.Id)
		}
	}
}

func NewQueue(name string) *Queue {
	queue := &Queue{
		Name:          name,
		Subscribers:   map[int]*Subscriber{},
		subDeleteChan: make(chan *Subscriber),
		subAddChan:    make(chan *Subscriber),
		msgPubChan:    make(chan *broker.Message),
	}
	go queue.Listener()
	return queue
}
