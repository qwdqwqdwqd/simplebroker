package broker

import (
	"context"
	"simplebroker/pkg/broker"
)

type Subscriber struct {
	Id          int
	Channel     chan broker.Message
	Ctx         context.Context
	subChanel   chan *broker.Message
	unSubSignal chan *Subscriber
	messages    []*broker.Message
}

func (s *Subscriber) SendMessages() {
	for {

		select {
		case <-s.Ctx.Done():
			go func() { s.unSubSignal <- s }()
			return
		case msg := <-s.subChanel:
			s.Channel <- *msg
		}

	}
}
func CreateNewSubscriber(ctx context.Context, ch chan broker.Message, unSubSignal chan *Subscriber) *Subscriber {
	id := 0 // TODO: create ID maker
	newSub := &Subscriber{
		Id:          id,
		Channel:     ch,
		Ctx:         ctx,
		subChanel:   make(chan *broker.Message),
		unSubSignal: unSubSignal,
		messages:    make([]*broker.Message, 0),
	}
	go newSub.SendMessages()
	return newSub
}
