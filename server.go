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

func (s *Server) serveConn(conn net.Conn) {
	fmt.Println("Serving connection!", conn.LocalAddr().String(), conn.LocalAddr().Network())
	rpc := NewRPC(conn)
	rpc.Connect()

	rpc.SendQueue <- &ServerTile{s.Tile}

	for msg := range rpc.RecvQueue {
		switch r := msg.(type) {
		case *ClientReplace:
			fmt.Println("got a ClientReplace")
			s.Tile.Cells[r.Y][r.X].Rune = r.Rune
		case *ClientMove:
			rpc.SendQueue <- ServerMove{r.X, r.Y}
		}
	}
}
