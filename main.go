package main

import (
	"flag"
	"os"
)

const defaultSockPath = "./test.sock"

type cmdFunc func([]string) error

func main() {
	cmd := server
	args := os.Args[1:]
	if len(os.Args) > 1 {
		switch cmdString := os.Args[1]; cmdString {
		case "server":
			cmd = server
			args = os.Args[2:]
		case "client":
			cmd = client
			args = os.Args[2:]
		}
	}
	err := cmd(args)
	if err != nil {
		os.Exit(1)
	}
}

func client(args []string) error {
	clientFlags := flag.NewFlagSet("client", flag.ExitOnError)
	socket := clientFlags.String("socket", defaultSockPath, "path to the socket to listen on")
	clientFlags.Parse(args)

	NewClient(*socket).Run()

	return nil
}

func server(args []string) error {
	serverFlags := flag.NewFlagSet("server", flag.ExitOnError)
	socket := serverFlags.String("socket", defaultSockPath, "path to the socket to listen on")
	serverFlags.Parse(args)

	NewServer(*socket).Run()
	return nil
}
