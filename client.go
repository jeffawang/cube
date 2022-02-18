package main

import (
	"fmt"
	"net"
	"os"
	"sync"

	tcell "github.com/gdamore/tcell/v2"
)

type Client struct {
	tile Tile
	rpc  RPC
}

func NewClient(sockpath string) *Client {
	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &Client{
		tile: NewTile(),
		rpc:  NewRPC(conn),
	}
}

func (c *Client) Run() {
	c.rpc.Connect()

	var tile *ServerTile
	resp := <-c.rpc.RecvQueue
	tile, ok := resp.(*ServerTile)
	if !ok {
		panic("no server tile!")
	}

	c.runGame(*tile)
}

func (c *Client) runGame(serverTile ServerTile) {
	s := MustScreen()

	player := NewPlayer()
	player.Rune = 'p'

	tile := serverTile.Tile

	entities := drawers{&tile, &player}
	entities.Draw(s)

	cleanupOnce := sync.Once{}
	defer cleanupOnce.Do(s.Fini)

	// Hot loop
hot:
	for {
		select {
		case msg := <-c.rpc.RecvQueue:
			switch m := msg.(type) {
			case (*ServerMove):
				player.X = m.X
				player.Y = m.Y
			}
		default:
		}

		entities.Draw(s)

		// Update screen
		s.Show()

		if s.HasPendingEvent() {
			// Poll event
			ev := s.PollEvent()

			// Process event
			switch ev := ev.(type) {
			case *tcell.EventResize:
				s.Sync()
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyCtrlC:
					cleanupOnce.Do(s.Fini)
					break hot
				case tcell.KeyCtrlL:
					s.Sync()
				case tcell.KeyLeft:
					// player.Move(-1, 0)
					c.rpc.SendQueue <- ClientMove{-1, 0}
				case tcell.KeyRight:
					// player.Move(1, 0)
					c.rpc.SendQueue <- ClientMove{1, 0}
				case tcell.KeyUp:
					// player.Move(0, -1)
					c.rpc.SendQueue <- ClientMove{0, -1}
				case tcell.KeyDown:
					// player.Move(0, 1)
					c.rpc.SendQueue <- ClientMove{0, 1}
				default:
					switch ev.Rune() {
					case 'C', 'c':
						panic("omg")
					case 'x':
						player.Insert(&tile, 'x')
						c.rpc.SendQueue <- ClientReplace{
							X: player.X, Y: player.Y, Rune: 'x',
						}
					}
				}
			}
		}
	}
}
