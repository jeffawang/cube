package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
)

const sockPath = "./test.sock"

func main() {
	fmt.Println(os.Args)
	if len(os.Args) < 2 || os.Args[1] == "server" {
		server()
	} else {
		client()
	}
}

type Response struct {
	Message string
}

func serveConn(conn net.Conn) {
	buf := bufio.NewWriter(conn)
	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(buf)

	var req Args
	resp := new(Response)

	fmt.Println("Serving connection!")

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
		resp.Message = "hi!"
		enc.Encode(resp)
		buf.Flush()
	}
}

func server() {
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

func client() {
	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	buf := bufio.NewWriter(conn)
	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(buf)

	args := Args{7, 8}
	enc.Encode(args)
	buf.Flush()

	var resp Response
	err = dec.Decode(&resp)
	if err != nil {
		fmt.Println("error decoding response", err)
	}
	fmt.Println("Got response:", resp)
}
