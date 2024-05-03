package main

import (
	"github.com/gdamore/tcell"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	up = iota
	down
	left
	right
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
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyUp:
				g.direction = up
			case tcell.KeyDown:
				g.direction = down
			case tcell.KeyLeft:
				g.direction = left
			case tcell.KeyRight:
				g.direction = right
			case tcell.KeyCtrlC:
				s.Fini()
				log.Printf("Screen size: %d x %d", g.width, g.height)
				os.Exit(0)
				return
			}
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
const looseMessage = "You loose"

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
			screen.Clear()
			drawLoose(screen)
			screen.Show()
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
	style := tcell.StyleDefault.Foreground(tcell.Color105).Background(tcell.Color105)
	for _, p := range snake.body {
		s.SetContent(p.x, p.y, '.', nil, style)
	}
}

func drawFood(s tcell.Screen, p Point) {
	style := tcell.StyleDefault.Foreground(tcell.Color126).Background(tcell.Color126)
	s.SetContent(p.x, p.y, 'o', nil, style)
}

func drawLoose(s tcell.Screen) {
	style := tcell.StyleDefault.Foreground(tcell.Color119).Background(tcell.Color90)
	startPoint := Point{x: 10, y: 10}

	for i, char := range []rune(looseMessage) {
		s.SetContent(startPoint.x+i, startPoint.y, char, nil, style)
	}
}
