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

	// for i := 0; i < 6; i++ {
	// 	args := Args{7, 8}

	// 	client.rpc.SendQueue <- &args
	// 	resp := <-client.rpc.RecvQueue

	// 	fmt.Println("Got response:", resp)
	// 	switch x := resp.(type) {
	// 	case *ServerMessage:
	// 		fmt.Println("++++++ SERVER message!")
	// 	case *ClientMessage:
	// 		fmt.Println("------ CLIENT message!")
	// 	case *ServerTile:
	// 		fmt.Println("++++++ ServerTile!")
	// 		tile = *x
	// 	default:
	// 		fmt.Println("shrugg??")
	// 	}
	// }
	// fmt.Println(tile)
	// resp.A()

	c.runGame(*tile)
}

func (c *Client) runGame(serverTile ServerTile) {
	s := MustScreen()

	// tile := NewTile()
	player := NewPlayer()

	tile := serverTile.Tile

	entities := drawers{&tile, &player}
	entities.Draw(s)

	cleanupOnce := sync.Once{}
	defer cleanupOnce.Do(s.Fini)

	// Hot loop
hot:
	for {

		entities.Draw(s)
		// Update screen
		s.Show()

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
				player.Move(-1, 0)
			case tcell.KeyRight:
				player.Move(1, 0)
			case tcell.KeyUp:
				player.Move(0, -1)
			case tcell.KeyDown:
				player.Move(0, 1)
			default:
				switch ev.Rune() {
				case 'C', 'c':
					panic("omg")
				case 'x':
					player.Insert(&tile, 'x')
					req := ClientReplace{
						X: player.X, Y: player.Y, Rune: 'x',
					}
					c.rpc.SendQueue <- req
				}
			}
		case *tcell.EventMouse:
			// x, y := ev.Position()
			button := ev.Buttons()
			// Only process button events, not wheel events
			button &= tcell.ButtonMask(0xff)

			// if button != tcell.ButtonNone && ox < 0 {
			// 	ox, oy = x, y
			// }
			// switch ev.Buttons() {
			// case tcell.ButtonNone:
			// 	if ox >= 0 {
			// 		label := fmt.Sprintf("%d,%d to %d,%d", ox, oy, x, y)
			// 		drawBox(s, ox, oy, x, y, boxStyle, label)
			// 		ox, oy = -1, -1
			// 	}
			// }
		}
	}
}
