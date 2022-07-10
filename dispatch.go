package main

import (
	"sync"
)

// Message DESIGNG  version + reuqestid + content
// version 1byte
// requestid 32zijie
// content protobuf
// req/resp handler+type  is_need_forward request+id

// unused
type Dispatcher struct {
	rwLock            *sync.RWMutex
	requestDispatcher map[string]chan ReadUnit
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		rwLock:            &sync.RWMutex{},
		requestDispatcher: make(map[string]chan ReadUnit),
	}
}

// message only have requestID and message
func (m *Dispatcher) Dispatch(message ReadUnit, uuid string) {
	m.rwLock.RLock()
	if ch, ok := m.requestDispatcher[uuid]; ok {
		ch <- message
	}
	m.rwLock.RUnlock()
}

func (m *Dispatcher) Add(uuid string) <-chan ReadUnit {
	m.rwLock.Lock()
	ch := make(chan ReadUnit)
	m.requestDispatcher[uuid] = ch
	m.rwLock.Unlock()
	return ch
}

func (m *Dispatcher) CloseChannel(uuid string) {
	m.rwLock.Lock()
	close(m.requestDispatcher[uuid])
	delete(m.requestDispatcher, uuid)
	m.rwLock.Unlock()
}
