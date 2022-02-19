package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)

// ==============================
// Server
// ==============================

type Server struct {
	Tile        Tile
	broadcastIn chan interface{}

	mu     sync.Mutex // Protects conns and id
	conns  map[*conn]struct{}
	nextID uint64
}

func NewServer() *Server {
	tile := NewTile()
	tile.Cells[3][3].Rune = 'y'
	return &Server{
		Tile:        tile,
		broadcastIn: make(chan interface{}, 10),
		conns:       make(map[*conn]struct{}),
		nextID:      1,
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
			for c := range s.conns {
				select {
				case c.broadcastOut <- msg:
				default:
					fmt.Println("Warning: attempted to broadcast to closed channel")
				}
			}
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("server error establishing connection!", err)
		}

		go s.newConn(conn).serve()
	}
}

func (s *Server) registerConn(c *conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conns == nil {
		s.conns = make(map[*conn]struct{})
	}
	s.conns[c] = struct{}{}

	// TODO: improve id allocation
	// - it is a side effect
	// - it does not garbage collect for disconnected clients
	// - it should probably happen in newConn?
	c.id = s.nextID
	s.nextID += 1
	fmt.Printf("New connection registered (id: %d)\n", c.id)
}

func (s *Server) deregisterConn(c *conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conns == nil {
		s.conns = make(map[*conn]struct{})
	}
	delete(s.conns, c)
	fmt.Printf("Connection deregistered (id: %d)\n", c.id)
}

// ==============================
// ClientState
// ==============================

type ClientState struct {
	Player
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

// ==============================
// conn
// ==============================

type conn struct {
	srv          *Server
	broadcastOut chan interface{}
	rwc          net.Conn
	rpc          RPC
	cs           ClientState
	id           uint64
}

func (s *Server) newConn(c net.Conn) *conn {
	ch := make(chan interface{}, 10)
	return &conn{
		srv:          s,
		broadcastOut: ch,
		rwc:          c,
		rpc:          NewRPC(c),
		cs:           NewClientState(),
	}
}

func (c *conn) serve() {
	c.srv.registerConn(c)
	defer c.srv.deregisterConn(c)
	c.rpc.Start()

	c.rpc.SendQueue <- ServerTile{c.srv.Tile}
	c.rpc.SendQueue <- ServerMove{X: c.cs.Player.X, Y: c.cs.Player.Y}

messageLoop:
	for {
		select {
		case msg, ok := <-c.rpc.RecvQueue:
			if !ok {
				break messageLoop
			}
			c.handleMessage(msg)
		case msg := <-c.broadcastOut:
			switch r := msg.(type) {
			case (ServerReplace):
				c.rpc.SendQueue <- r
			}
		}
	}
}

func (c *conn) handleMessage(msg interface{}) {
	switch r := msg.(type) {
	case *ClientPut:
		c.srv.Tile.Cells[r.Y][r.X].Rune = r.Rune
		c.srv.broadcastIn <- ServerReplace{
			X:    r.X,
			Y:    r.Y,
			Rune: r.Rune,
		}
	case *ClientMove:
		x := c.cs.Player.X + r.X
		y := c.cs.Player.Y + r.Y
		if x < 0 || x >= WIDTH || y < 0 || y >= WIDTH {
			return
		}
		c.cs.Player.X = x
		c.cs.Player.Y = y
		c.rpc.SendQueue <- ServerMove{x, y}
	}
}
