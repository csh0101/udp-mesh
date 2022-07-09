package main

import (
	"context"
	"net"
	"sync"
)

type ConnectionManager struct {
	rwLock     *sync.RWMutex
	connection map[string]Connection
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		rwLock:     &sync.RWMutex{},
		connection: make(map[string]Connection),
	}
}

func (manager *ConnectionManager) Start(ctx context.Context) error {
	go func(ctx context.Context) {
		<-ctx.Done()
		manager.Close()
	}(ctx)
	return nil
}

// Close this Close() Should be call after all connetion are not used
func (manager *ConnectionManager) Close() {
	for _, v := range manager.connection {
		v.Close()
	}
}

// the function uuid is by google uuid algothrim for remote host
func GetConnection(manager *ConnectionManager, uuid string) Connection {
	manager.rwLock.RLock()
	if conn, ok := manager.connection[uuid]; ok {
		return conn
	}
	manager.rwLock.RUnlock()
	return Connection{}
}

func SetConnection(manager *ConnectionManager, uuid string, connection Connection) {
	manager.rwLock.Lock()
	manager.connection[uuid] = connection
	manager.rwLock.Unlock()
}

type Connection struct {
	uuid string
	conn *net.UDPConn
	once *sync.Once
}

func NewConnection(manager *ConnectionManager, uuid string, conn *net.UDPConn) Connection {
	return Connection{
		uuid: uuid,
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
