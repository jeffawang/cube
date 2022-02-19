package main

import "github.com/spf13/cobra"

const defaultSockPath = "./test.sock"

type cmdFunc func([]string) error

func main() {
	cmd().Execute()
}

func cmd() *cobra.Command {

	root := &cobra.Command{
		Use:   "cube",
		Short: "hi!",
	}
	root.AddCommand(serverCommand())
	root.AddCommand(clientCommand())
	return root
}

func serverCommand() *cobra.Command {
	var socket string
	c := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			NewServer(socket).Run()
			return nil
		},
	}
	c.Flags().StringVarP(&socket, "socket", "s", defaultSockPath, "socket to listen on")
	return c
}

func clientCommand() *cobra.Command {
	var socket string
	c := &cobra.Command{
		Use: "client",
		RunE: func(cmd *cobra.Command, args []string) error {
			NewClient(socket).Run()
			return nil
		},
	}
	c.Flags().StringVarP(&socket, "socket", "s", defaultSockPath, "socket to listen on")
	return c
}
