package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

func runClient(sockPath string) {
	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	buf := bufio.NewWriter(conn)
	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(buf)

	for i := 0; i < 6; i++ {
		args := Args{7, 8}
		enc.Encode(args)
		buf.Flush()

		var resp interface{}
		err = dec.Decode(&resp)
		if err != nil {
			fmt.Println("error decoding response", err)
		}
		fmt.Println("Got response:", resp)
		switch resp.(type) {
		case *ServerMessage:
			fmt.Println("++++++ SERVER message!")
		case *ClientMessage:
			fmt.Println("------ CLIENT message!")
		case *ServerTile:
			fmt.Println("++++++ ServerTile!")
		default:
			fmt.Println("shrugg??")
		}
	}
	// resp.A()

	// runGame()
}
