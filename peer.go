package main

import (
	"context"
	"fmt"
	"log"
	"net"
)

type Peer struct {
}

func (p *Peer) Run(ctx context.Context, in <-chan WriteUnit) <-chan ReadUnit {
	log.Println("running...")
	if in != nil {
		go WriteLoop(ctx, in)
	}
	return ReadLoop(ctx)
}

type WriteFunc func(conn *net.UDPConn, data []byte)

func WriteFuncExample() WriteFunc {
	return func(conn *net.UDPConn, data []byte) {
		if n, err := conn.Write(data); err != nil {
			log.Println(err)
			return
		} else {
			log.Println(fmt.Sprintf("write %d byte data", n))
		}
	}
}

type WriteUnit struct {
	conn *net.UDPConn
	data []byte
	f    WriteFunc
}

func WriteLoop(ctx context.Context, in <-chan WriteUnit) {
	for u := range in {
		u.f(u.conn, u.data)
	}
}

type ReadUnit struct {
	message []byte
	addr    net.Addr
}

func ReadLoop(ctx context.Context) <-chan ReadUnit {
	listener, err := net.ListenPacket("udp4", ":18888")
	if err != nil {
		panic(err)
	}
	out := make(chan ReadUnit)
	go func(ctx context.Context, out chan ReadUnit, listener net.PacketConn) {
	LOOP:
		for {
			select {
			case <-ctx.Done():
				break LOOP
			default:
			}
			buf := make([]byte, 1024)
			n, addr, err := listener.ReadFrom(buf)
			if err != nil {
				panic(err)
			}
			out <- ReadUnit{
				message: buf[:n],
				addr:    addr,
			}
		}
		close(out)
		listener.Close()
	}(ctx, out, listener)
	return out
}
