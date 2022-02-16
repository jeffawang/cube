package main

import (
	"fmt"
	"os"
)

const sockPath = "./test.sock"

type Args struct {
	A int
	B int
}

func main() {
	fmt.Println(os.Args)
	if len(os.Args) < 2 || os.Args[1] == "server" {
		runServer(sockPath)
	} else {
		NewClient(sockPath).Run()
	}
}
