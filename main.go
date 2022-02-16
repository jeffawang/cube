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

func main() {
	s := mustScreen()

	t := NewTile()

	t.Draw(s)

	// Event loop
	// ox, oy := -1, -1
	cleanupOnce := sync.Once{}
	defer cleanupOnce.Do(s.Fini)

	// Hot loop
hot:
	for {
		// Update screen
		s.Show()

		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				cleanupOnce.Do(s.Fini)
				break hot
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
				panic("omg")
				s.Clear()
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
