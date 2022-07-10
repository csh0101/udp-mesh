package main

import (
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type EchoServer struct {
	udpServer *Server
	echo      *echo.Echo
}

func NewEchoServer(config Config) *EchoServer {
	return &EchoServer{
		udpServer: NewServer(config),
		echo:      echo.New(),
	}
}

func (e *EchoServer) Start() {
	serve := e.echo
	udp := e.udpServer
	serve.GET("/host/list", func(c echo.Context) error {
		m := udp.mesh.Copy()
		for _, v := range m.GetRemotePeer() {
			conn := GetConnection(udp.manager, v, nil, &net.UDPAddr{
				IP:   net.ParseIP(v),
				Port: 18888,
			})
			wu := WriteUnit{
				conn: conn.conn,
				f:    WriteFuncExample(),
			}
			message := &HostInfoMessage{
				Req: true,
				Arp: &ARP{
					Src:            udp.config.LocalAddr,
					Dst:            v,
					Forwardaddress: []string{udp.config.LocalAddr},
				},
			}
			u, err := uuid.NewUUID()
			if err != nil {
				panic(err)
			}

			var pack []byte

			pack = append(pack, []byte{1}...)

			pack = append(pack, []byte(u.String())...)

			b, err := proto.Marshal(message)

			if err != nil {
				panic(err)
			}
			pack = append(pack, b...)

			readCh := udp.dispatch.Add(u)
			udp.writeCh <- wu

			af := time.After(time.Second * 10)

		}
		return nil
	})
}
