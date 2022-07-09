package main

import (
	"sync"
)

// unused
type HostInfoDispatcher struct {
	rwLock            *sync.RWMutex
	requestDispatcher map[string]chan HostInfoMessageResponse
}

func (m *HostInfoDispatcher) Dispatch(response HostInfoMessageResponse, uuid string) {
	m.rwLock.RLock()
	m.requestDispatcher[uuid] <- response
	m.rwLock.RUnlock()
}

func (m *HostInfoDispatcher) Add(uuid string) {
	m.rwLock.Lock()
	ch := make(chan HostInfoMessageResponse)
	m.requestDispatcher[uuid] = ch
	m.rwLock.Unlock()
}

func (m *HostInfoDispatcher) CloseChannel(uuid string) {
	m.rwLock.Lock()
	close(m.requestDispatcher[uuid])
	delete(m.requestDispatcher, uuid)
	m.rwLock.Unlock()
}

// Filter filet the smae tuple (source + reuqestID)
type Filter struct {
	rwLock     *sync.RWMutex
	requestMap map[string]struct{}
}
