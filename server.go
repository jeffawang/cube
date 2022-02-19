package main

import (
	"fmt"
	"net"
	"os"
)

type Server struct {
	Tile         Tile
	broadcastIn  chan interface{}
	broadcastOut []chan interface{}
}

func NewServer() *Server {
	tile := NewTile()
	tile.Cells[3][3].Rune = 'y'
	return &Server{
		Tile:         tile,
		broadcastIn:  make(chan interface{}, 10),
		broadcastOut: make([]chan interface{}, 0),
	}
}

func (s *Server) Run(sockPath string) {
	os.Remove(sockPath)
	listener, err := net.Listen("unix", sockPath)
	if err != nil {
		fmt.Println("error listening on socket", sockPath)
		os.Exit(1)
	}

	fmt.Println("Listening for new connections on", sockPath)

	go func() {
		for msg := range s.broadcastIn {
			for _, ch := range s.broadcastOut {
				ch <- msg
			}
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("server error establishing connection!", err)
		}
		fmt.Println("New connection!")

		// connChan := make(chan interface{}, 10)
		// s.broadcastOut = append(s.broadcastOut, connChan)

		go s.newConn(conn).serve()
	}
}

type ClientState struct {
	Player
}

type conn struct {
	srv          *Server
	broadcastOut chan interface{}
	rwc          net.Conn
	rpc          RPC
}

func (s *Server) newConn(c net.Conn) *conn {
	ch := make(chan interface{}, 10)
	s.broadcastOut = append(s.broadcastOut, ch)
	return &conn{
		srv:          s,
		broadcastOut: ch,
		rwc:          c,
		rpc:          NewRPC(c),
	}
}

func (c *conn) serve() {
	fmt.Println("Serving connection!", c.rwc.LocalAddr())
	c.rpc.Connect()

	cs := NewClientState()

	c.rpc.SendQueue <- ServerTile{c.srv.Tile}
	c.rpc.SendQueue <- ServerMove{X: cs.Player.X, Y: cs.Player.Y}

	fmt.Println("started!")
messageLoop:
	for {
		select {
		case msg, ok := <-c.rpc.RecvQueue:
			if !ok {
				break messageLoop
			}
			switch r := msg.(type) {
			case *ClientReplace:
				fmt.Println("got a ClientReplace")
				c.srv.Tile.Cells[r.Y][r.X].Rune = r.Rune
				c.srv.broadcastIn <- ServerReplace{
					X:    r.X,
					Y:    r.Y,
					Rune: r.Rune,
				}

			case *ClientMove:
				x := cs.Player.X + r.X
				y := cs.Player.Y + r.Y
				if x < 0 || x >= WIDTH || y < 0 || y >= WIDTH {
					continue
				}
				cs.Player.X = x
				cs.Player.Y = y
				c.rpc.SendQueue <- ServerMove{x, y}
			}
		case msg := <-c.broadcastOut:
			switch r := msg.(type) {
			case (ServerReplace):
				c.rpc.SendQueue <- r
			}
		}
	}
	fmt.Println("finished!")
}

func NewClientState() ClientState {
	return ClientState{
		Player: Player{
			X:    WIDTH / 2,
			Y:    WIDTH / 2,
			Rune: 'p',
		},
	}
}
