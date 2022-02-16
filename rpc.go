package main

import (
	"sync"

	tcell "github.com/gdamore/tcell/v2"
)

func runGame() {
	s := MustScreen()

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
}
