package main

import "sync"

type Mesh struct {
	rwLock     *sync.RWMutex
	knownPeers []string
	peerMap    map[string]struct{}
}

// Mesh Copy is a calculate the
func (m *Mesh) Copy() *Mesh {
	var dst []string
	peerMap := make(map[string]struct{})

	copy(dst, m.knownPeers)

	for _, v := range dst {
		peerMap[v] = struct{}{}
	}

	return &Mesh{
		rwLock:     &sync.RWMutex{},
		knownPeers: dst,
		peerMap:    peerMap,
	}
}

func (m *Mesh) AddPeer(peer string) {
	m.rwLock.RLock()
	if _, ok := m.peerMap[peer]; !ok {
		return
	}
	m.rwLock.RUnlock()

	m.rwLock.Lock()
	m.knownPeers = append(m.knownPeers, peer)
	m.peerMap[peer] = struct{}{}
	m.rwLock.Unlock()
}

func (m *Mesh) FilterForward(Peers []string) []string {

	res := make([]string, 0)

	m.rwLock.RLock()
	for _, v := range Peers {
		if _, ok := m.peerMap[v]; !ok {
			res = append(res, v)
		}
	}
	m.rwLock.Unlock()

	return res
}
