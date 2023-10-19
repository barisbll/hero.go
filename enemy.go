package main

import (
	"math/rand"
	"time"
)

var enemyEmojis = []rune{'👹', '👾', '👽', '👻', '🤡'}

func giveRandomEnemyEmoji() rune {
	return enemyEmojis[rand.Intn(len(enemyEmojis))]
}

func giveRandomLocation(maxWidth, maxHeight int) (int, int) {
	direction := rand.Intn(4)
	switch direction {
	case 0:
		return rand.Intn(maxWidth), 0
	case 1:
		return rand.Intn(maxWidth), maxHeight - 1
	case 2:
		return 0, rand.Intn(maxHeight)
	case 3:
		return maxWidth - 1, rand.Intn(maxHeight)
	}
	return 0, 0
}

type Enemy struct {
	id                int
	currentX          int
	currentY          int
	speed             uint8
	isDead            bool
	emoji             rune
	hero              *Hero
	lastDrawnPosition Coordinates
}

func NewEnemy(id, maxWidth, maxHeight int, hero *Hero) *Enemy {
	x, y := giveRandomLocation(maxWidth, maxHeight)
	return &Enemy{
		id:                id,
		currentX:          x,
		currentY:          y,
		speed:             1,
		isDead:            false,
		emoji:             giveRandomEnemyEmoji(),
		hero:              hero,
		lastDrawnPosition: Coordinates{x: 0, y: 0},
	}
}

func (e *Enemy) move() {
	if e.isDead || e.hero.isDead || e.hero.isPaused {
		return
	}

	isEnemyZombie := true
	for _, enemy := range e.hero.enemies {
		if enemy.id == e.id {
			isEnemyZombie = false
			break
		}
	}

	if isEnemyZombie {
		e.isDead = true
		return
	}

	if e.currentX == e.hero.x && e.currentY == e.hero.y {
		e.isDead = true
		e.hero.isDead = true
	}

	xDifference := makePositive(e.currentX - e.hero.x)
	yDifference := makePositive(e.currentY - e.hero.y)

	if xDifference >= yDifference {
		if e.currentX < e.hero.x {
			e.currentX += 1
		} else if e.currentX > e.hero.x {
			e.currentX -= 1
		}
		return
	} else if xDifference < yDifference {
		if e.currentY < e.hero.y {
			e.currentY += 1
		} else if e.currentY > e.hero.y {
			e.currentY -= 1
		}
	}

}

func (e *Enemy) draw(ticker *time.Ticker) {
	s := *e.hero.settings.screen
	style := *e.hero.settings.theme

	go func() {
		for {
			select {
			case <-ticker.C:
				s.SetContent(e.lastDrawnPosition.x, e.lastDrawnPosition.y, ' ', nil, style)
				if !e.isDead {
					e.lastDrawnPosition.x = e.currentX
					e.lastDrawnPosition.y = e.currentY
					s.SetContent(e.lastDrawnPosition.x, e.lastDrawnPosition.y, e.emoji, nil, style)
					e.move()
				} else {
					s.SetContent(e.currentX, e.currentY, DeadEmoji, nil, style)
					ticker.Stop()
					s.Show()
					// Make the dead animation stay for a second
					time.Sleep(1 * time.Second)
					s.SetContent(e.currentX, e.currentY, ' ', nil, style)
					s.Show()
					return
				}

				s.Show()

			}
		}
	}()
}
