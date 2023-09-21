package main

import (
	"fmt"

	tcell "github.com/gdamore/tcell/v2"
)

// TODO: add logic to delete the bullets when they are out of scope
const BulletEmoji rune = '>'

type Bullet struct {
	x1           int
	y1           int
	x2           int
	y2           int
	currentX     int
	currentY     int
	distanceX    int
	distanceY    int
	xDonePercent float32
	yDonePercent float32
	bulletRune   tcell.Key

	speed uint
}

type Direction string

const (
	X Direction = "X"
	Y Direction = "Y"
)

func (b *Bullet) setBulletRune() {
	switch true {
	case b.distanceX > 0 && b.distanceY > 0:
		b.bulletRune = tcell.KeyUpLeft
	case b.distanceX > 0 && b.distanceY < 0:
		b.bulletRune = tcell.KeyDownLeft
	case b.distanceX < 0 && b.distanceY > 0:
		b.bulletRune = tcell.KeyUpRight
	case b.distanceX < 0 && b.distanceY < 0:
		b.bulletRune = tcell.KeyDownRight
	}
}

func (b *Bullet) calculateDonePercent(direction Direction) {
	if direction == "X" {
		takenDistance := makePositive(b.currentX - makePositive(b.x1))
		b.xDonePercent = float32(takenDistance) / makePositiveFloat(float32(b.distanceX))
		fmt.Printf("X done percent %f", b.xDonePercent)
	}
	if direction == "Y" {
		takenDistance := makePositive(b.currentY - makePositive(b.y1))
		b.yDonePercent = float32(takenDistance) / makePositiveFloat(float32(b.distanceY))
		fmt.Printf("Y done percent %f", b.yDonePercent)
	}
}

func (b *Bullet) moveOne(direction Direction) {
	if direction == "X" {
		if b.distanceX > 0 {
			b.currentX++
		}
		b.currentX--
		fmt.Printf("\nCurrent X %d\n", b.currentX)
	}
	if direction == "Y" {
		if b.distanceY > 0 {
			b.currentY++
		}
		b.currentY--
	}
	b.calculateDonePercent(direction)
}

func (b *Bullet) move() {
	if b.xDonePercent == 0 && b.yDonePercent == 0 {
		b.distanceX = b.x2 - b.x1
		b.distanceY = b.y2 - b.y1

		// fmt.Printf("disX %d, disY %d", b.distanceX, b.distanceY)

		if makePositive(b.distanceX) > makePositive(b.distanceY) {
			b.moveOne(X)
		} else {
			b.moveOne(Y)
		}
	}

	if b.xDonePercent > b.yDonePercent {
		b.moveOne(Y)
	} else {
		b.moveOne(X)
	}
}

func (b *Bullet) draw(s tcell.Screen, style tcell.Style) {
	for i := 0; i <= int(3); i++ {
		b.move()
	}
	s.SetContent(b.currentX, b.currentY, rune(b.bulletRune), nil, style)
}
