package main

import (
	"os"

	"github.com/jessevdk/go-flags"
)

const defaultSockPath = "./test.sock"

type serverOpts struct {
	sockPath string `short:"s" long:"socket"`
}

func newServerOpts(args []string) *serverOpts {
	opts := serverOpts{
		sockPath: defaultSockPath,
	}
	flags.ParseArgs(&opts, args)
	return &opts
}

func (s *serverOpts) run() {
	NewServer(s.sockPath).Run()
}

func main() {
	if len(os.Args) == 1 {
		newServerOpts(os.Args[1:]).run()
	} else {
		if os.Args[1] == "server" {
		}
	}
}
