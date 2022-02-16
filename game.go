package main

import (
	"log"

	tcell "github.com/gdamore/tcell/v2"
)

type Drawer interface {
	Draw(tcell.Screen)
}

type drawers []Drawer

func (ds drawers) Draw(s tcell.Screen) {
	for _, d := range ds {
		d.Draw(s)
	}
}

func MustScreen() tcell.Screen {
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
