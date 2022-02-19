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
	c := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			socket, err := cmd.Flags().GetString("socket")
			if err != nil {
				return err
			}
			NewServer(socket).Run()
			return nil
		},
	}
	c.Flags().StringP("socket", "s", defaultSockPath, "socket to listen on")
	return c
}

func clientCommand() *cobra.Command {
	c := &cobra.Command{
		Use: "client",
		RunE: func(cmd *cobra.Command, args []string) error {
			socket, err := cmd.Flags().GetString("socket")
			if err != nil {
				return err
			}
			NewClient(socket).Run()
			return nil
		},
	}
	c.Flags().StringP("socket", "s", defaultSockPath, "socket to listen on")
	return c
}
