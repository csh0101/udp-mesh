package main

import (
	"context"
	"net"
	"sync"

	"github.com/golang/protobuf/proto"
)

var (
	ReqFlag  byte = 0
	RespFlag byte = 1
)

type Server struct {
	manager   *ConnectionManager
	peer      *Peer
	mesh      *Mesh
	writeCh   chan WriteUnit
	ctxCancel context.CancelFunc
	ctx       context.Context
	dispatch  *Dispatcher
	once      *sync.Once
	config    Config
	wg        *sync.WaitGroup
	// key is request_id and src
}

func NewServer(config Config) *Server {
	peer := &Peer{}
	c, cancel := context.WithCancel(context.Background())
	writeCh := make(chan WriteUnit)
	return &Server{
		config:    config,
		ctx:       c,
		peer:      peer,
		ctxCancel: cancel,
		once:      &sync.Once{},
		writeCh:   writeCh,
		mesh:      NewMesh(),
		dispatch:  NewDispatcher(),
		manager:   NewConnectionManager(),
	}
}

func (s *Server) handleHostInfo(data []byte, protoBuffer *proto.Buffer) error {
	// data 0:16 是http requestID
	protoBuffer.SetBuf(data[16:])
	mes := &HostInfoMessage{}
	err := protoBuffer.Unmarshal(mes)
	if err != nil {
		panic(err)
	}
	arp := mes.Arp
	if mes.Req {
		{
			laster := arp.Forwardaddress[len(arp.Forwardaddress)-1]
			temp := mes.Arp.Src
			mes.Arp.Src = s.config.LocalAddr
			mes.Arp.Dst = temp
			mes.Req = false
			mes.Host = s.config.LocalAddr
			conn := GetConnection(s.manager, laster, nil, &net.UDPAddr{
				IP:   net.ParseIP(laster),
				Port: 18888,
			})
			mes.Arp.Forwardaddress = mes.Arp.Forwardaddress[:len(arp.Forwardaddress)-1]
			var pack []byte
			pack = append(pack, []byte{1}...)
			pack = append(pack, data[:16]...)
			b, err := proto.Marshal(mes)
			if err != nil {
				panic(err)
			}
			pack = append(pack, b...)
			WriteFuncExample()(conn.conn, pack)
		}
		mesh := s.mesh.Copy()
		if len(mesh.GetRemotePeer()) >= 1 {
			for _, peer := range mesh.knownPeers {
				flag := false
				// avoid the reference send
				for _, forward := range arp.Forwardaddress {
					if peer == forward {
						flag = true
					}
				}
				if !flag {
					SendMessage(*mes, s, peer, data)
				}
			}
		}
	} else {
		if arp.Dst == s.config.LocalAddr {
			b := data[:16]
			s.mesh.AddPeer(mes.Host)
			s.dispatch.Dispatch(ReadUnit{
				message: []byte(mes.Host),
			}, string(b))
		} else {
			laster := arp.Forwardaddress[len(arp.Forwardaddress)-1]
			conn := GetConnection(s.manager, laster, nil, &net.UDPAddr{
				IP:   net.ParseIP(laster),
				Port: 18888,
			})
			temp := mes.Arp.Src
			mes.Arp.Src = s.config.LocalAddr
			mes.Arp.Dst = temp
			mes.Arp.Forwardaddress = mes.Arp.Forwardaddress[:len(arp.Forwardaddress)-1]
			var pack []byte
			pack = append(pack, []byte{1}...)
			pack = append(pack, data[:16]...)
			b, err := proto.Marshal(mes)
			if err != nil {
				panic(err)
			}
			pack = append(pack, b...)
			WriteFuncExample()(conn.conn, pack)
		}
	}
	return nil
}

func (s *Server) Start(wg *sync.WaitGroup) {
	c := s.peer.Run(s.ctx, s.writeCh)
	s.manager.Start(s.ctx, wg)
	go func(in <-chan ReadUnit, wg *sync.WaitGroup) {
		defer wg.Done()
		b := proto.NewBuffer(nil)
		for v := range in {
			data := v.message
			switch data[0] {
			case 0x01:
				s.handleHostInfo(data[1:], b)
			case 0x02:
			}
		}
	}(c, wg)
}

func (s *Server) Close() {
	s.once.Do(func() {
		s.ctxCancel()
		s.manager.Close()
		s.wg.Wait()
	})
}

func SendMessage(mes HostInfoMessage, s *Server, host string, data []byte) {
	conn := GetConnection(s.manager, host, nil, &net.UDPAddr{
		IP:   net.ParseIP(host),
		Port: 18888,
	})
	mes.Arp.Forwardaddress = append(mes.Arp.Forwardaddress, s.config.LocalAddr)
	var pack []byte
	pack = append(pack, []byte{1}...)
	pack = append(pack, data[:16]...)
	b, err := proto.Marshal(&mes)
	if err != nil {
		panic(err)
	}
	pack = append(pack)
	WriteFuncExample()(conn.conn, b)
}
