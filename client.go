package main

import (
	"fmt"
	"net"
	"os"
)

type Client struct {
	tile Tile
	rpc  RPC
}

func NewClient(sockpath string) Client {
	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return Client{
		tile: NewTile(),
		rpc:  NewRPC(conn),
	}
}

func runClient(sockPath string) {
	client := NewClient(sockPath)
	client.rpc.Connect()

	var tile ServerTile

	for i := 0; i < 6; i++ {
		args := Args{7, 8}

		client.rpc.SendQueue <- &args
		resp := <-client.rpc.RecvQueue

		fmt.Println("Got response:", resp)
		switch x := resp.(type) {
		case *ServerMessage:
			fmt.Println("++++++ SERVER message!")
		case *ClientMessage:
			fmt.Println("------ CLIENT message!")
		case *ServerTile:
			fmt.Println("++++++ ServerTile!")
			tile = *x
		default:
			fmt.Println("shrugg??")
		}
	}
	fmt.Println(tile)
	// resp.A()

	// runGame(tile)
}
