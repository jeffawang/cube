package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
)

// TODO: move the client message somewhere?
func init() {
	gob.Register(&ClientMessage{})
	gob.Register(&ServerMessage{})
	gob.Register(&ServerTile{})
}

// ServerTile is a tile that the server sends down to the client
type ServerTile struct {
	Tile
}

// ServerMessage is a message sent by the server down to clients.
type ServerMessage struct {
	Number int
	Even   bool
}

// ClientMessage is a message sent by the client up to the server.
type ClientMessage struct {
	Number int
	Even   bool
}

var serverTile = ServerTile{NewTile()}

func init() {
	serverTile.Cells[3][3].Rune = 'y'
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
	buf := bufio.NewWriter(conn)
	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(buf)

	var req Args
	var msg interface{}

	fmt.Println("Serving connection!", conn.LocalAddr().String(), conn.LocalAddr().Network())

	for {
		err := dec.Decode(&req)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client closed connection")
			} else {
				fmt.Println("error decoding request", err)
			}
			return
		}
		fmt.Println("got request:", req)
		msg = serverTile
		err = enc.Encode(&msg)
		if err != nil {
			fmt.Println("uh oh encoding", err)
		}
		err = buf.Flush()
		if err != nil {
			fmt.Println("uh oh", err)
		}

	}
}
