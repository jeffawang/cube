package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"

	tl "github.com/JoelOtter/termloop"
)

const sockPath = "./test.sock"

type Player struct {
	*tl.Entity
}

func (p *Player) Tick(event tl.Event) {
	if event.Type == tl.EventKey {
		x, y := p.Position()
		switch event.Key {
		case tl.KeyArrowRight:
			p.SetPosition(x+1, y)
		case tl.KeyArrowLeft:
			p.SetPosition(x-1, y)
		case tl.KeyArrowUp:
			p.SetPosition(x, y-1)
		case tl.KeyArrowDown:
			p.SetPosition(x, y+1)
		}
	}
}

func main() {
	game := tl.NewGame()
	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorGreen,
		Fg: tl.ColorBlack,
		Ch: 'v',
	})
	level.AddEntity(tl.NewRectangle(10, 10, 50, 20, tl.ColorBlue))

	player := Player{tl.NewEntity(1, 1, 1, 1)}
	player.SetCell(0, 0, &tl.Cell{Fg: tl.ColorRed, Ch: '옷'})
	level.AddEntity(&player)

	game.Screen().SetLevel(level)

	game.Start()

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
