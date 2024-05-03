package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

const (
	up    = 1
	down  = -1
	left  = 2
	right = -2
)

type Game struct {
	width, height int
	snake         *Snake
	food          Point
	direction     int
	over          bool
}

func (g *Game) tick() {
	if g.over {
		return
	}

	g.snake.move(g.direction)

	if g.snake.isCollidingWith(g.food) {
		g.placeFood()
	} else {
		g.snake.popTail()
	}

	g.checkOver()
}

func (g *Game) listenKeys(s tcell.Screen) {
	direction := g.direction

	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			// arrow keys
			switch ev.Key() {
			case tcell.KeyUp:
				direction = up
			case tcell.KeyDown:
				direction = down
			case tcell.KeyLeft:
				direction = left
			case tcell.KeyRight:
				direction = right

			case tcell.KeyCtrlC:
				s.Fini()
				log.Printf("Screen size: %d x %d", g.width, g.height)
				os.Exit(0)
				return
			}

			// vim keys
			switch ev.Rune() {
			case 'h':
				direction = left
			case 'j':
				direction = down
			case 'k':
				direction = up
			case 'l':
				direction = right
			}
		}

		if (direction + g.direction) != 0 {
			g.direction = direction
		}
	}
}

func newGame(width, height int) Game {
	snake := Snake{body: []Point{{x: 10, y: 10}, {x: 11, y: 10}, {x: 12, y: 10}}}
	game := Game{snake: &snake, direction: right, width: width, height: height}
	game.placeFood()

	return game
}

func (g *Game) placeFood() {
	for {
		x := rand.Intn(g.width) + 1
		y := rand.Intn(g.height) + 2
		newFood := Point{x: x, y: y}

		if !g.snake.isCollidingWith(newFood) {
			g.food = newFood
			break
		}
	}
}

func (g *Game) checkOver() {
	head := g.snake.body[0]
	if head.x < 0 || head.x > (g.width-1) || head.y < 0 || head.y > (g.height-1) {
		g.over = true
	}
}

type Snake struct {
	body []Point
}

func (s *Snake) move(direction int) {
	var newHead Point = s.body[0]

	switch direction {
	case up:
		newHead.y--
	case down:
		newHead.y++
	case left:
		newHead.x--
	case right:
		newHead.x++
	}

	s.body = append([]Point{newHead}, s.body...)
}

func (s *Snake) popTail() {
	s.body = s.body[:len(s.body)-1]
}

func (s *Snake) isCollidingWith(p Point) bool {
	for _, b := range s.body {
		if b.x == p.x && b.y == p.y {
			return true
		}
	}

	return false
}

type Point struct {
	x, y int
}

const delay = 80 * time.Millisecond

func main() {
	screen := setupScreen()
	screenWidth, screenHeight := screen.Size()
	game := newGame(screenWidth, screenHeight)

	go game.listenKeys(screen)

	for {
		screen.Clear()
		game.tick()

		drawSnake(screen, game.snake)
		drawFood(screen, game.food)
		screen.Show()
		time.Sleep(delay)

		if game.over {
			fmt.Println("You loose !")
			os.Exit(0)
		}
	}
}

func setupScreen() tcell.Screen {
	screen, err := tcell.NewScreen()

	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	return screen
}

func drawSnake(s tcell.Screen, snake *Snake) {
	bodyStyle := tcell.StyleDefault.Foreground(tcell.Color105).Background(tcell.Color105)
	for _, p := range snake.body {
		s.SetContent(p.x, p.y, '.', nil, bodyStyle)
	}
}

func drawFood(s tcell.Screen, p Point) {
	style := tcell.StyleDefault.Foreground(tcell.Color126).Background(tcell.Color126)
	s.SetContent(p.x, p.y, 'o', nil, style)
}
