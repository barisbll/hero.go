package main

import (
	"sync"
	"time"

	tcell "github.com/gdamore/tcell/v2"
)

type Direction string
type VerticalDirection string

const (
	TOP         Direction = "TOP"
	RIGHT       Direction = "RIGHT"
	LEFT        Direction = "LEFT"
	BOTTOM      Direction = "BOTTOM"
	TOPRIGHT    Direction = "TOPRIGHT"
	TOPLEFT     Direction = "TOPLEFT"
	BOTTOMRIGHT Direction = "BOTTOMRIGHT"
	BOTTOMLEFT  Direction = "BOTTOMLEFT"
)

const (
	X VerticalDirection = "X"
	Y VerticalDirection = "Y"
)

const BombEmoji rune = '💣'
const ExplosionEmoji rune = '💥'

type Coordinates struct {
	x int
	y int
}

type Bomb struct {
	bombId            int
	currentX          int
	currentY          int
	distanceX         int
	distanceY         int
	xDonePercent      float32
	yDonePercent      float32
	direction         Direction
	initialX          int
	initialY          int
	finalX            int
	finalY            int
	lastDrawnPosition Coordinates
	explodeIn         time.Duration
	isDead            bool

	speed uint
}

func calculateDirection(xDifference, yDifference int) Direction {
	switch true {
	case xDifference == 0 && yDifference > 0:
		return BOTTOM
	case xDifference == 0 && yDifference < 0:
		return TOP
	case xDifference > 0 && yDifference == 0:
		return RIGHT
	case xDifference < 0 && yDifference == 0:
		return LEFT
	case xDifference > 0 && yDifference > 0:
		return BOTTOMRIGHT
	case xDifference > 0 && yDifference < 0:
		return TOPRIGHT
	case xDifference < 0 && yDifference > 0:
		return BOTTOMLEFT
	case xDifference < 0 && yDifference < 0:
		return TOPLEFT
	default:
		return ""
	}
}

// will be used later when we add options to the game (non-stop bombs options for exp)
func calculateFinalPosition(maxX, maxY, currentX, currentY, distanceX, distanceY int, direction Direction) (int, int) {
	switch true {
	case direction == TOP:
		return currentX, 0
	case direction == BOTTOM:
		return currentX, maxY
	case direction == RIGHT:
		return maxX, currentY
	case direction == LEFT:
		return 0, currentY
	case direction == TOPRIGHT:
		leftForY := currentY
		leftForX := maxX - currentX
		if leftForY < leftForX {
			timesToMultiply := leftForX / distanceX
			return currentX + (distanceX * timesToMultiply), currentY - (distanceY * timesToMultiply)
		}
		timesToMultiply := leftForY / distanceY
		return currentX + (distanceX * timesToMultiply), currentY - (distanceY * timesToMultiply)
	default:
		return 0, 0
	}
}

func (b *Bomb) draw(s tcell.Screen, style tcell.Style, ticker *time.Ticker, bombExploded chan string, explosionComplete chan struct{}, explosionWaitGroup *sync.WaitGroup) {
	go func() {
		time.Sleep(b.explodeIn)
		bombExploded <- ""
		bombExploded <- ""
		time.Sleep(1 * time.Second)
		close(explosionComplete)
	}()

	go func() {
		for {
			select {
			case <-ticker.C:
				if b.lastDrawnPosition.x != 0 && b.lastDrawnPosition.y != 0 {
					s.SetContent(b.lastDrawnPosition.x, b.lastDrawnPosition.y, ' ', nil, style)
				}
				if !b.isDead {
					s.SetContent(b.currentX, b.currentY, BombEmoji, nil, style)
					b.lastDrawnPosition.x = b.currentX
					b.lastDrawnPosition.y = b.currentY
					b.move()
				} else {
					s.SetContent(b.currentX, b.currentY, ExplosionEmoji, nil, style)
				}

				s.Show()
			case <-bombExploded:
				s.SetContent(b.currentX, b.currentY, ExplosionEmoji, nil, style)
				b.lastDrawnPosition.x = b.currentX
				b.lastDrawnPosition.y = b.currentY
				explosionWaitGroup.Done()
				explosionWaitGroup.Wait()
			case <-explosionComplete:
				s.SetContent(b.currentX, b.currentY, ' ', nil, style)
				s.Show()
				ticker.Stop()
				return
			}

		}
	}()
}

func (b *Bomb) calculateDonePercent(verticalDirection VerticalDirection) {
	if verticalDirection == X {
		// do not exceed the finalX
		if b.xDonePercent >= 1 {
			return
		}

		// calculate xDonePercent for RIGHT
		if b.direction == RIGHT || b.direction == BOTTOMRIGHT || b.direction == TOPRIGHT {
			takenDistance := b.currentX - b.initialX
			b.xDonePercent = float32(takenDistance) / float32(b.distanceX)
			return
		}

		// calculate xDonePercent for LEFT
		takenDistance := (b.currentX - b.initialX) * -1
		b.xDonePercent = float32(takenDistance) / float32(b.distanceX)
		return
	} else if verticalDirection == Y {
		// do not exceed the finalY
		if b.yDonePercent >= 1 {
			return
		}

		// calculate yDonePercent for TOP
		if b.direction == TOP || b.direction == TOPRIGHT || b.direction == TOPLEFT {
			takenDistance := (b.currentY - b.initialY) * -1
			b.yDonePercent = float32(takenDistance) / float32(b.distanceY)
			return
		}

		// calculate yDonePercent for BOTTOM
		takenDistance := b.currentY - b.initialY
		b.yDonePercent = float32(takenDistance) / float32(b.distanceY)
		return
	}
}

func (b *Bomb) move() {
	if b.xDonePercent == 0 && b.yDonePercent == 0 {
		// First move
		if b.distanceX > b.distanceY {
			if b.direction == RIGHT || b.direction == BOTTOMRIGHT || b.direction == TOPRIGHT {
				b.currentX++
				b.calculateDonePercent(X)
			} else {
				b.currentX--
				b.calculateDonePercent(X)
			}
		} else {
			if b.direction == TOP || b.direction == TOPRIGHT || b.direction == TOPLEFT {
				b.currentY--
				b.calculateDonePercent(Y)
			} else {
				b.currentY++
				b.calculateDonePercent(Y)
			}
		}
		return
	}

	if b.xDonePercent >= 1 && b.yDonePercent >= 1 {
		return
	}

	// Move in X direction
	if b.yDonePercent >= b.xDonePercent {
		if b.direction == RIGHT || b.direction == BOTTOMRIGHT || b.direction == TOPRIGHT {
			b.currentX++
			b.calculateDonePercent(X)
		} else {
			b.currentX--
			b.calculateDonePercent(X)
		}
		return
	} else if b.xDonePercent > b.yDonePercent {
		// Move in Y direction
		if b.direction == TOP || b.direction == TOPRIGHT || b.direction == TOPLEFT {
			b.currentY--
			b.calculateDonePercent(Y)
		} else {
			b.currentY++
			b.calculateDonePercent(Y)
		}
	}
}
