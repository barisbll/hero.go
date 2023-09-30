package main

import (
	"sync"
	"time"

	tcell "github.com/gdamore/tcell/v2"
)

const HeroEmoji rune = '\U0001F9B8' // 🦸‍♂️
const DeadEmoji rune = '\U0001F480' // 💀

type Hero struct {
	x             int
	y             int
	speed         uint8
	bombs         []Bomb
	bombIdCounter int
	isDead        bool
}

func (h *Hero) goRight(maxWidth int) {
	if h.isDead {
		return
	}

	if h.x < maxWidth-2 {
		h.x += int(h.speed)
	}
}

func (h *Hero) goLeft() {
	if h.isDead {
		return
	}

	if h.x > 0 {
		h.x -= int(h.speed)
	}
}

func (h *Hero) goUp() {
	if h.isDead {
		return
	}

	if h.y > 0 {
		h.y -= int(h.speed)
	}
}

func (h *Hero) goDown(maxHeight int) {
	if h.isDead {
		return
	}

	if h.y < maxHeight-1 {
		h.y += int(h.speed)
	}
}

func (h *Hero) draw(s tcell.Screen, style tcell.Style) {
	if h.isDead {
		s.SetContent(h.x, h.y, DeadEmoji, nil, style)
		return
	}

	s.SetContent(h.x, h.y, HeroEmoji, nil, style)
}

func (h *Hero) addBomb(s tcell.Screen, style tcell.Style, clickedX, clickedY int) {
	if len(h.bombs) >= 5 || h.isDead {
		return
	}

	finalX, finalY := clickedX, clickedY

	bomb := Bomb{
		bombId:    h.bombIdCounter,
		currentX:  h.x,
		currentY:  h.y,
		distanceX: makePositive(finalX - h.x),
		distanceY: makePositive(finalY - h.y),
		direction: calculateDirection(clickedX-h.x, clickedY-h.y),
		speed:     1,
		initialX:  h.x,
		initialY:  h.y,
		finalX:    finalX,
		finalY:    finalY,
		explodeIn: time.Second * 1,
		isDead:    false,
	}

	h.bombIdCounter++
	h.bombs = append(h.bombs, bomb)

	ticker := time.NewTicker(20 * time.Millisecond)
	bombExploded := make(chan string)
	explosionComplete := make(chan struct{})

	var explosionWaitGroup sync.WaitGroup
	explosionWaitGroup.Add(2)

	bomb.draw(s, style, ticker, bombExploded, explosionComplete, &explosionWaitGroup)

	go func() {
		for {
			<-bombExploded

			var indexToRemove int = -1

			for i, heroBomb := range h.bombs {
				if heroBomb.bombId == bomb.bombId {
					bomb.isDead = true
					h.killTheThingsInTheExplosionArea(s, style, bomb.currentX, bomb.currentY)
					indexToRemove = i
					break
				}
			}

			if indexToRemove != -1 {
				// Create a new slice without the element to remove
				h.bombs = append(h.bombs[:indexToRemove], h.bombs[indexToRemove+1:]...)
			}

			explosionWaitGroup.Done()
			explosionWaitGroup.Wait()
		}
	}()

}

func (h *Hero) isNearBomb(bombX, bombY, actorX, actorY int) bool {
	distanceX := makePositive(bombX - actorX)
	distanceY := makePositive(bombY - actorY)

	if distanceX <= 2 && distanceY <= 2 {
		return true
	}

	return false
}

func (h *Hero) killTheThingsInTheExplosionArea(s tcell.Screen, style tcell.Style, bombX, bombY int) {
	if h.isNearBomb(bombX, bombY, h.x, h.y) {
		h.isDead = true
		s.SetContent(h.x, h.y, DeadEmoji, nil, style)
		s.Show()
		return
	}
}
