package main

import (
	"fmt"
	"log"

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

	hero := Hero{x: xmax / 2, y: ymax / 2, speed: 1, isDead: false, enemyIdCounter: 24}

	hero.draw(s, defStyle)
	hero.spanNewEnemies(s, defStyle, xmax, ymax)

	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			fmt.Println(maybePanic)
			// panic(maybePanic)
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
			} else if ev.Rune() == 'R' || ev.Rune() == 'r' {
				if hero.isDead {
					for _, enemy := range hero.enemies {
						enemy.isDead = true
					}
					hero = Hero{x: xmax / 2, y: ymax / 2, speed: 1, isDead: false}
					s.Clear()
				}
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
				hero.addBomb(s, defStyle, clickedX, clickedY)
			}
		}

		hero.draw(s, defStyle)

	}
}
