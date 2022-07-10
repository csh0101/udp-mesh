package main

import "sync"

var (
	BroadCastDomain string = "255.255.255.255"
)

type Mesh struct {
	rwLock     *sync.RWMutex
	knownPeers []string
	peerMap    map[string]struct{}
}

func NewMesh() *Mesh {
	return &Mesh{
		rwLock:     &sync.RWMutex{},
		knownPeers: make([]string, 0),
		peerMap:    make(map[string]struct{}),
	}
}
func (m *Mesh) GetRemotePeer() []string {
	if len(m.knownPeers) == 0 {
		return []string{BroadCastDomain}
	}
	return m.knownPeers
}

// Mesh Copy is a calculate the
func (m *Mesh) Copy() *Mesh {
	m.rwLock.Lock()
	var dst []string
	peerMap := make(map[string]struct{})
	copy(dst, m.knownPeers)
	m.rwLock.Unlock()
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
