package main

import (
	"log"
	"time"

	tcell "github.com/gdamore/tcell/v2"
)

func main() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()
	xmax, ymax := s.Size()

	hero := Hero{x: xmax / 2, y: ymax / 2, speed: 1}

	hero.draw(s, defStyle)

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	// Event loop
	for {
		// Update screen
		s.Show()

		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
			xmax, ymax = s.Size()
		case *tcell.EventKey:
			// System logic keys
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
				s.Clear()
			} else if ev.Key() == tcell.KeyRight || ev.Rune() == 'D' || ev.Rune() == 'd' {
				hero.goRight(xmax)
				s.Clear()
			} else if ev.Key() == tcell.KeyLeft || ev.Rune() == 'A' || ev.Rune() == 'a' {
				hero.goLeft()
				s.Clear()
			} else if ev.Key() == tcell.KeyUp || ev.Rune() == 'W' || ev.Rune() == 'w' {
				hero.goUp()
				s.Clear()
			} else if ev.Key() == tcell.KeyDown || ev.Rune() == 'S' || ev.Rune() == 's' {
				hero.goDown(ymax)
				s.Clear()
			}
		case *tcell.EventMouse:
			clickedX, clickedY := ev.Position()

			switch ev.Buttons() {
			case tcell.Button1, tcell.Button2:
				finalX, finalY := clickedX, clickedY

				bomb := Bomb{
					currentX:  hero.x,
					currentY:  hero.y,
					distanceX: makePositive(finalX - hero.x),
					distanceY: makePositive(finalY - hero.y),
					direction: calculateDirection(clickedX-hero.x, clickedY-hero.y),
					speed:     1,
					initialX:  hero.x,
					initialY:  hero.y,
					finalX:    finalX,
					finalY:    finalY,
					explodeIn: time.Second * 1,
					isDead:    false,
				}

				hero.bombs = append(hero.bombs, bomb)
				bomb.draw(s, defStyle)
			}
		}

		hero.draw(s, defStyle)
	}
}
