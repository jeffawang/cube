package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"

	tcell "github.com/gdamore/tcell/v2"
)

const sockPath = "./test.sock"

type Cell struct {
	Rune  rune
	Color tcell.Style
}

const WIDTH = 25
const SCREENWIDTH = WIDTH + 2

// 25 rows of 80 columns
type Tile struct {
	Cells [WIDTH][WIDTH]Cell
}

func NewTile() Tile {
	return Tile{[WIDTH][WIDTH]Cell{}}
}

func (t *Tile) Draw(s tcell.Screen) {
	for y, row := range t.Cells {
		for x, col := range row {
			s.SetContent(1+x, 1+y, col.Rune, nil, tcell.StyleDefault)
		}
	}
	e := SCREENWIDTH - 1
	for i := 1; i < e; i++ {
		s.SetContent(i, 0, tcell.RuneHLine, nil, tcell.StyleDefault)
		s.SetContent(i, e, tcell.RuneHLine, nil, tcell.StyleDefault)
		s.SetContent(0, i, tcell.RuneVLine, nil, tcell.StyleDefault)
		s.SetContent(e, i, tcell.RuneVLine, nil, tcell.StyleDefault)
	}
	s.SetContent(0, 0, tcell.RuneULCorner, nil, tcell.StyleDefault)
	s.SetContent(0, e, tcell.RuneLLCorner, nil, tcell.StyleDefault)
	s.SetContent(e, 0, tcell.RuneURCorner, nil, tcell.StyleDefault)
	s.SetContent(e, e, tcell.RuneLRCorner, nil, tcell.StyleDefault)
}

func mustScreen() tcell.Screen {
	defaultStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defaultStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()
	return s
}

type Drawer interface {
	Draw(tcell.Screen)
}

type Player struct {
	X, Y int
	Rune rune
}

func (p *Player) Draw(s tcell.Screen) {
	s.SetContent(1+p.X, 1+p.Y, p.Rune, nil, tcell.StyleDefault)
}

func (p *Player) Move(dx, dy int) {
	newX := p.X + dx
	newY := p.Y + dy
	if 0 <= newX && newX < WIDTH {
		p.X = newX
	}
	if 0 <= newY && newY < WIDTH {
		p.Y = newY
	}
}

func (p *Player) Insert(t *Tile, r rune) {
	t.Cells[p.Y][p.X].Rune = r
	p.Move(1, 0)
}

func NewPlayer() Player {
	return Player{
		X:    WIDTH / 2,
		Y:    WIDTH / 2,
		Rune: 'p',
	}
}

type drawers []Drawer

func (ds drawers) Draw(s tcell.Screen) {
	for _, d := range ds {
		d.Draw(s)
	}
}

func main() {
	s := mustScreen()

	tile := NewTile()
	player := NewPlayer()

	entities := drawers{&tile, &player}
	entities.Draw(s)

	// tile.Draw(s)

	// Event loop
	// ox, oy := -1, -1
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

	return
	fmt.Println(os.Args)
	if len(os.Args) < 2 || os.Args[1] == "server" {
		server()
	} else {
		client()
	}
}

type Blah interface {
	A()
}

// ServerMessage is a message sent by the server down to clients.
type ServerMessage struct {
	Number int
	Even   bool
}

func (s *ServerMessage) A() {
	fmt.Println("I'm a ServerMessage!")
}

// ClientMessage is a message sent by the client up to the server.
type ClientMessage struct {
	Number int
	Even   bool
}

func (s *ClientMessage) A() {
	fmt.Println("I'm a ClientMessage!")
}

func init() {
	gob.Register(&ClientMessage{})
	gob.Register(&ServerMessage{})
}

func serveConn(conn net.Conn) {
	buf := bufio.NewWriter(conn)
	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(buf)

	var req Args
	var resp Blah
	sm := new(ServerMessage)
	cm := new(ClientMessage)

	fmt.Println("Serving connection!", conn.LocalAddr().String(), conn.LocalAddr().Network())

	defer fmt.Println("omg")

	i := 0

	for {
		err := dec.Decode(&req)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client closed connection")
			} else {
				fmt.Println("error decoding request", err)
			}
			return
		}
		fmt.Println("got request:", req)

		if i%2 == 0 {
			sm.Number = i
			sm.Even = i%2 == 0
			resp = sm
		} else {
			cm.Number = i
			cm.Even = i%2 == 0
			resp = cm
		}
		enc.Encode(&resp)
		enc.Encode(&resp)
		enc.Encode(&resp)
		enc.Encode(&resp)
		buf.Flush()
		i += 1
	}
}

func server() {
	os.Remove(sockPath)
	listener, err := net.Listen("unix", sockPath)
	if err != nil {
		fmt.Println("error listening on socket", sockPath)
		os.Exit(1)
	}

	fmt.Println("Listening for new connections on", sockPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("server error establishing connection!", err)
		}
		fmt.Println("New connection!")
		go serveConn(conn)
	}
}

func client() {
	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	buf := bufio.NewWriter(conn)
	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(buf)

	args := Args{7, 8}
	enc.Encode(args)
	buf.Flush()

	for {

		var resp Blah
		err = dec.Decode(&resp)
		if err != nil {
			fmt.Println("error decoding response", err)
		}
		fmt.Println("Got response:", resp)
		resp.A()

	}
}
