package main

import (
	"fmt"
	"net"
	"os"
)

var serverTile = ServerTile{NewTile()}

func init() {
	serverTile.Cells[3][3].Rune = 'y'
}

type Server struct {
	rpc RPC
}

func runServer(sockPath string) {
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
		go serveConn(conn)
	}
}

func serveConn(conn net.Conn) {
	fmt.Println("Serving connection!", conn.LocalAddr().String(), conn.LocalAddr().Network())
	rpc := NewRPC(conn)
	rpc.Connect()

	rpc.SendQueue <- serverTile

	for msg := range rpc.RecvQueue {
		switch r := msg.(type) {
		case *ClientReplace:
			fmt.Println("got a ClientReplace")
			handleClientReplace(*r)
		case *ClientMove:
			rpc.SendQueue <- ServerMove{r.X, r.Y}
		}
	}
}

func handleClientReplace(cr ClientReplace) {
	serverTile.Cells[cr.Y][cr.X].Rune = cr.Rune
}
