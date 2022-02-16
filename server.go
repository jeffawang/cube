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
	var resp interface{}
	sm := new(ServerMessage)
	cm := new(ClientMessage)

	fmt.Println("Serving connection!", conn.LocalAddr().String(), conn.LocalAddr().Network())

	defer fmt.Println("omg")

	i := 0

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

		if i%2 == 0 {
			sm.Number = i
			sm.Even = i%2 == 0
			resp = sm
		} else {
			cm.Number = i
			cm.Even = i%2 == 0
			resp = cm
		}
		enc.Encode(&resp)
		buf.Flush()
		i += 1
	}
}
