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
	settings := NewSettings(&s, &defStyle)

	(*settings.screen).SetStyle(defStyle)
	(*settings.screen).EnableMouse()
	(*settings.screen).EnablePaste()
	(*settings.screen).Clear()

	hero := NewHero(settings)
	menu := NewMenu(settings, hero)

	quit := func() {
		maybePanic := recover()
		(*settings.screen).Fini()
		if maybePanic != nil {
			fmt.Println(maybePanic)
			// panic(maybePanic)
		}
	}
	defer quit()

	// Event loop
	for {
		// Update screen
		(*settings.screen).Show()

		// Poll event
		ev := (*settings.screen).PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			(*settings.screen).Sync()
			xmax, ymax := (*settings.screen).Size()
			settings.xMax = xmax
			settings.yMax = ymax
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
					hero = NewHero(settings)
					menu.hero = hero
					(*settings.screen).Clear()
				}
			} else if ev.Rune() == 'P' || ev.Rune() == 'p' {
				hero.isPaused = !hero.isPaused
				if hero.isPaused {
					menu.drawPause()
				} else {
					menu.clearPause()
				}

			} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
				(*settings.screen).Clear()
			} else if ev.Key() == tcell.KeyRight || ev.Rune() == 'D' || ev.Rune() == 'd' {
				hero.goRight()
			} else if ev.Key() == tcell.KeyLeft || ev.Rune() == 'A' || ev.Rune() == 'a' {
				hero.goLeft()
			} else if ev.Key() == tcell.KeyUp || ev.Rune() == 'W' || ev.Rune() == 'w' {
				hero.goUp()
			} else if ev.Key() == tcell.KeyDown || ev.Rune() == 'S' || ev.Rune() == 's' {
				hero.goDown()
			}
		case *tcell.EventMouse:
			clickedX, clickedY := ev.Position()

			switch ev.Buttons() {
			case tcell.Button1, tcell.Button2:
				hero.addBomb(clickedX, clickedY)
			}
		}

		hero.draw()
		menu.drawScore()
	}
}
