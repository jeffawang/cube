package main

import (
	"os"
)

const sockPath = "./test.sock"

type Args struct {
	A int
	B int
}

func main() {
	// fmt.Println(os.Args)
	if len(os.Args) < 2 || os.Args[1] == "server" {
		NewServer().Run(sockPath)
	} else {
		if len(os.Args) > 2 {
			NewClient(os.Args[2]).Run()
		} else {
			NewClient(sockPath).Run()
		}
	}
}
