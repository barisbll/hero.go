package main

import (
	"fmt"
)

type Menu struct {
	settings *Settings
	hero     *Hero
}

type Corner string

const (
	CORNER_TOP_RIGHT    Corner = "TOPRIGHT"
	CORNER_TOP_LEFT     Corner = "TOPLEFT"
	CORNER_BOTTOM_RIGHT Corner = "BOTTOMRIGHT"
	CORNER_BOTTOM_LEFT  Corner = "BOTTOMLEFT"
)

func NewMenu(settings *Settings, hero *Hero) *Menu {
	return &Menu{
		settings: settings,
		hero:     hero,
	}
}

func (m *Menu) printString(text string, corner Corner) {
	xMax := m.settings.xMax
	yMax := m.settings.yMax
	screen := *m.settings.screen
	style := *m.settings.theme

	switch corner {
	case CORNER_TOP_RIGHT:
		for i, char := range text {
			screen.SetContent(xMax-len(text)+i, 0, char, nil, style)
		}
	case CORNER_TOP_LEFT:
		for i, char := range text {
			screen.SetContent(i+1, 0, char, nil, style)
		}
	case CORNER_BOTTOM_RIGHT:
		for i, char := range text {
			screen.SetContent(xMax-len(text)+i, yMax-1, char, nil, style)
		}
	case CORNER_BOTTOM_LEFT:
		for i, char := range text {
			screen.SetContent(i+1, yMax-1, char, nil, style)
		}
	}
}

func (m *Menu) drawScore() {
	m.printString(fmt.Sprintf("SCORE: %d", m.hero.enemyIdCounter), CORNER_TOP_RIGHT)
}

func (m *Menu) drawPause() {
	m.printString("PAUSE", CORNER_TOP_LEFT)
}

func (m *Menu) clearPause() {
	m.printString("     ", CORNER_TOP_LEFT)
}
