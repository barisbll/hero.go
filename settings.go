package main

import tcell "github.com/gdamore/tcell/v2"

type Settings struct {
	screen *tcell.Screen
	theme  *tcell.Style
	xMax   int
	yMax   int
}

func NewSettings(screen *tcell.Screen, theme *tcell.Style) *Settings {
	xMax, yMax := (*screen).Size()

	return &Settings{
		screen: screen,
		theme:  theme,
		xMax:   xMax,
		yMax:   yMax,
	}
}
