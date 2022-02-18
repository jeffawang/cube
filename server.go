package main

import (
	"fmt"
	"net"
	"os"
)

type Server struct {
	Tile Tile
}

func NewServer() *Server {
	tile := NewTile()
	tile.Cells[3][3].Rune = 'y'
	return &Server{
		Tile: tile,
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

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("server error establishing connection!", err)
		}
		fmt.Println("New connection!")
		go s.serveConn(conn)
	}
}

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

func (s *Server) serveConn(conn net.Conn) {
	fmt.Println("Serving connection!", conn.LocalAddr())
	rpc := NewRPC(conn)
	rpc.Connect()

	cs := NewClientState()

	rpc.SendQueue <- ServerTile{s.Tile}
	rpc.SendQueue <- ServerMove{X: cs.Player.X, Y: cs.Player.Y}

	fmt.Println("started!")
	for msg := range rpc.RecvQueue {
		switch r := msg.(type) {
		case *ClientReplace:
			fmt.Println("got a ClientReplace")
			s.Tile.Cells[r.Y][r.X].Rune = r.Rune
			rpc.SendQueue <- ServerReplace{
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
			rpc.SendQueue <- ServerMove{x, y}
		}
	}
	fmt.Println("finished!")
}
