package main

import tcell "github.com/gdamore/tcell/v2"

const WIDTH = 25
const SCREENWIDTH = WIDTH + 2

// ==============================
// Tile
// ==============================

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

func (t *Tile) Replace(x, y int, r rune) {
	t.Cells[y][x].Rune = r
}

// ==============================
// Cell
// ==============================

type Cell struct {
	Rune rune
	// Color tcell.Style
}

// ==============================
// Player
// ==============================

type Player struct {
	X, Y int
	Rune rune
}

func NewPlayer() Player {
	return Player{}
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
