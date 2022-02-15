package main

import (
	"fmt"
	"net"
	"net/rpc"
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

func server() {
	os.Remove(sockPath)
	listener, err := net.Listen("unix", sockPath)

	fmt.Printf("accepting connections on %s:%s\n", listener.Addr().Network(), listener.Addr())
	arith := new(Arith)
	rpc.Register(arith)

	rpc.Accept(listener)

	fmt.Println(listener.Addr(), err)
}

func client() {
	client, err := rpc.Dial("unix", sockPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	args := Args{7, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Got response:", reply)
}
