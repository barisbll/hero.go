package main

import tcell "github.com/gdamore/tcell/v2"

const HeroEmoji rune = '\U0001F9B8' // 🦸‍♂️

type Hero struct {
	x       int
	y       int
	speed   uint8
	bullets []Bullet
}

func (h *Hero) goRight(maxWidth int) {
	if h.x < maxWidth-2 {
		h.x += int(h.speed)
	}
}

func (h *Hero) goLeft() {
	if h.x > 0 {
		h.x -= int(h.speed)
	}
}

func (h *Hero) goUp() {
	if h.y > 0 {
		h.y -= int(h.speed)
	}
}

func (h *Hero) goDown(maxHeight int) {
	if h.y < maxHeight-1 {
		h.y += int(h.speed)
	}
}

func (h *Hero) draw(s tcell.Screen, style tcell.Style) {
	s.SetContent(h.x, h.y, HeroEmoji, nil, style)

}
