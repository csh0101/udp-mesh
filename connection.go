package main

import (
	"context"
	"net"
	"sync"
)

type ConnectionManager struct {
	rwLock     *sync.RWMutex
	connection map[string]*Connection
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		rwLock:     &sync.RWMutex{},
		connection: make(map[string]*Connection),
	}
}

func (manager *ConnectionManager) Start(ctx context.Context, wg *sync.WaitGroup) error {
	wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup) {
		<-ctx.Done()
		manager.Close()
		wg.Done()
	}(ctx, wg)
	return nil
}

// Close this Close() Should be call after all connetion are not used
func (manager *ConnectionManager) Close() {
	for _, v := range manager.connection {
		v.Close()
	}
}

// the function uuid is by google uuid algothrim for remote host
func GetConnection(manager *ConnectionManager, host string, laddr *net.UDPAddr, raddr *net.UDPAddr) *Connection {
	manager.rwLock.RLock()
	conn, ok := manager.connection[host]
	if ok {
		return conn

	}
	manager.rwLock.RUnlock()

	if conn == nil {
		manager.rwLock.Lock()
		if _, ok := manager.connection[host]; ok {
			return manager.connection[host]
		}
		u, err := net.DialUDP("udp", nil, raddr)
		if err != nil {
			panic(err)
		}
		manager.connection[host] = NewConnection(manager, host, u)
		manager.rwLock.Unlock()
	}
	return manager.connection[host]

}

type Connection struct {
	host string
	conn *net.UDPConn
	once *sync.Once
}

func NewConnection(manager *ConnectionManager, host string, conn *net.UDPConn) *Connection {
	return &Connection{
		host: host,
		conn: conn,
		once: &sync.Once{},
	}
}

func (c *Connection) Close() error {
	c.once.Do(func() {
		if c.conn != nil {
			c.conn.Close()
		}
	})
	return nil
}
