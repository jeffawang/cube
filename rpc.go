package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"net"
)

func init() {
	// Messages from server
	gob.Register(&ServerTile{})
	gob.Register(&ServerMove{})
	gob.Register(&ServerReplace{})

	// Messages from client
	gob.Register(&ClientPut{})
	gob.Register(&ClientMove{})
}

// ==============================
// Server messages
// ==============================

// ServerTile is a tile that the server sends down to the client
type ServerTile struct {
	Tile
}

type ServerReplace struct {
	X, Y int
	Rune rune
}

type ServerMove struct {
	X, Y int
}

// ==============================
// Client messages
// ==============================

// ClientPut represents the intent to replace a cell's rune
// in a tile.
type ClientPut struct {
	X, Y int
	Rune rune
}

type ClientMove struct {
	X, Y int
}

// ==============================
// RPC framework
// ==============================

type RPC struct {
	SendQueue chan interface{}
	RecvQueue chan interface{}

	enc *gob.Encoder
	dec *gob.Decoder
	buf *bufio.Writer
}

func NewRPC(conn net.Conn) RPC {
	buf := bufio.NewWriter(conn)

	return RPC{
		SendQueue: make(chan interface{}, 10),
		RecvQueue: make(chan interface{}, 10),
		buf:       buf,
		dec:       gob.NewDecoder(conn),
		enc:       gob.NewEncoder(buf),
	}
}

func (r *RPC) Start() {
	go func() {
		defer close(r.RecvQueue)
		for {
			var resp interface{}
			err := r.dec.Decode(&resp)
			if err == io.EOF {
				return
			} else if err != nil {
				fmt.Println("uh oh", err)
				continue
			}
			r.RecvQueue <- resp
		}
	}()
	go func() {
		defer close(r.SendQueue)
		for val := range r.SendQueue {
			err := r.enc.Encode(&val)
			if err == io.EOF {
				return
			} else if err != nil {
				fmt.Printf("Problem encoding message (%v): %s\n", val, err.Error())
				continue
			}
			err = r.buf.Flush()
			if err != nil {
				fmt.Printf("Problem flushing send buffer (%v): %s\n", val, err.Error())
			}
		}
	}()
}
