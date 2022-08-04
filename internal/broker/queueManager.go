package broker

import "sync"

type QueueManager struct {
	queues map[string]*Queue
	sync.RWMutex
}

func (qm *QueueManager) GetQueue(name string) (*Queue, bool) {
	qm.RLock()
	defer qm.RUnlock()
	queue, ok := qm.queues[name]
	return queue, ok
}

func (qm *QueueManager) GetOrCreateQueue(name string) (*Queue, bool) {
	queue, ok := qm.GetQueue(name)
	if !ok {
		queue := qm.CreateQueue(name)
		return queue, true
	}
	return queue, false
}

func (ts *QueueManager) CreateQueue(name string) *Queue {
	ts.Lock()
	defer ts.Unlock()
	queue := NewQueue(name)
	ts.queues[name] = queue
	return queue
}

func CreateQueue() *QueueManager {
	return &QueueManager{
		queues: map[string]*Queue{},
	}
}
