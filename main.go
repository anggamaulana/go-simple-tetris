// Simple Tetris Game
// created by angga maulana
// http://github.com/anggamaulana

package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
)

const (
	WIDTH  = 20
	HEIGHT = 20
)

var World [WIDTH][HEIGHT]int
var Objects []Shape
var ActiveObject Shape
var dropObject bool = true

type Point struct {
	X int
	Y int
}

type Shape struct {
	TYPE     int
	Position Point
	Body     []Point
	Angle    int
}

func (s *Shape) Create(typeobj int) {
	switch typeobj {
	case 0: //----
		s.TYPE = 1
		s.Position = Point{X: 10, Y: 0}
		s.Angle = 0
		s.Body = append(s.Body, s.Position)
		s.Body = append(s.Body, Point{X: 9, Y: 0})
		s.Body = append(s.Body, Point{X: 11, Y: 0})
		s.Body = append(s.Body, Point{X: 12, Y: 0})
		break
	case 1: // L
		s.TYPE = 2
		s.Position = Point{X: 10, Y: 1}
		s.Angle = 0
		s.Body = append(s.Body, s.Position)
		s.Body = append(s.Body, Point{X: 10, Y: 0})
		s.Body = append(s.Body, Point{X: 10, Y: 2})
		s.Body = append(s.Body, Point{X: 11, Y: 2})
		break
	case 2: // flip L
		s.TYPE = 3
		s.Position = Point{X: 10, Y: 1}
		s.Angle = 0
		s.Body = append(s.Body, s.Position)
		s.Body = append(s.Body, Point{X: 10, Y: 0})
		s.Body = append(s.Body, Point{X: 10, Y: 2})
		s.Body = append(s.Body, Point{X: 9, Y: 2})
		break
	case 3: //vertical zigzag
		s.TYPE = 4
		s.Position = Point{X: 10, Y: 1}
		s.Angle = 0
		s.Body = append(s.Body, s.Position)
		s.Body = append(s.Body, Point{X: 10, Y: 0})
		s.Body = append(s.Body, Point{X: 9, Y: 1})
		s.Body = append(s.Body, Point{X: 9, Y: 2})
		break
	case 4: //square
		s.TYPE = 5
		s.Position = Point{X: 10, Y: 1}
		s.Angle = 0
		s.Body = append(s.Body, s.Position)
		s.Body = append(s.Body, Point{X: 10, Y: 0})
		s.Body = append(s.Body, Point{X: 9, Y: 0})
		s.Body = append(s.Body, Point{X: 9, Y: 1})
		break
	}
}

func (s *Shape) isFalling() bool {

	s.MoveBody()

	for _, val := range s.Body {
		nextXPoint := val.X
		nextYPoint := val.Y + 1
		isOwnBody := false
		for _, p := range s.Body {
			if nextXPoint == p.X && nextYPoint == p.Y {
				isOwnBody = true
				break
			}
		}
		if World[nextYPoint][nextXPoint] == 1 && !isOwnBody {
			return false
		}
	}

	for i, _ := range s.Body {

		s.Body[i] = Point{X: s.Body[i].X, Y: s.Body[i].Y + 1}
	}
	return true
}

func (s *Shape) MoveBody() {
	for _, val := range s.Body {
		World[val.Y][val.X] = 0
	}
}

func (s *Shape) ShiftRight() {
	maxX := 0
	Ypos := 0
	for _, val := range s.Body {
		if val.X > maxX {
			maxX = val.X
			Ypos = val.Y
		}
	}
	if maxX+1 > WIDTH-1 {
		return
	}

	if World[Ypos][maxX+1] != 1 {
		s.MoveBody()
		for i, _ := range s.Body {
			s.Body[i] = Point{X: s.Body[i].X + 1, Y: s.Body[i].Y}
		}
	}

}

func (s *Shape) ShiftLeft() {
	minX := 100
	Ypos := 0
	for _, val := range s.Body {
		if val.X < minX {
			minX = val.X
			Ypos = val.Y
		}
	}
	if minX-1 < 0 {
		return
	}

	if World[Ypos][minX-1] != 1 {
		s.MoveBody()
		for i, _ := range s.Body {
			s.Body[i] = Point{X: s.Body[i].X - 1, Y: s.Body[i].Y}
		}
	}
}

func (s *Shape) RotateCW() {
	// p'x = cos(theta) * (px-ox) - sin(theta) * (py-oy) + ox
	// p'y = sin(theta) * (px-ox) + cos(theta) * (py-oy) + oy

	tmpPoint := []Point{}
	for _, val := range s.Body {
		lenX := int(math.Abs(float64(val.X - s.Position.X)))
		lenY := int(math.Abs(float64(val.Y - s.Position.Y)))
		newX := (0 * lenX) - (1 * lenY) + s.Position.X
		newY := (1 * lenX) + (0 * lenY) + s.Position.Y

		if newX < 0 || newX > WIDTH-1 || newY < 0 || newY > HEIGHT || World[val.Y][val.X] == 0 {
			return
		}

		newPoint := Point{X: newX, Y: newY}
		tmpPoint = append(tmpPoint, newPoint)
	}

	s.MoveBody()
	for j, _ := range tmpPoint {
		if s.Body[j].X == s.Position.X && s.Body[j].Y == s.Position.Y {
			s.Position = tmpPoint[j]
		}

		s.Body[j] = tmpPoint[j]
	}
}

func (s *Shape) RotateCCW() {
	tmpPoint := []Point{}
	for _, val := range s.Body {
		lenX := int(math.Abs(float64(val.X - s.Position.X)))
		lenY := int(math.Abs(float64(val.Y - s.Position.Y)))
		newX := (0 * lenX) + (1 * lenY) + s.Position.X
		newY := (1 * lenX) - (0 * lenY) + s.Position.Y

		if newX < 0 || newX > WIDTH-1 || newY < 0 || newY > HEIGHT || World[val.Y][val.X] == 0 {

			return
		}

		newPoint := Point{X: newX, Y: newY}
		tmpPoint = append(tmpPoint, newPoint)

	}

	s.MoveBody()
	for j, _ := range tmpPoint {
		if s.Body[j].X == s.Position.X && s.Body[j].Y == s.Position.Y {
			s.Position = tmpPoint[j]
		}

		s.Body[j] = tmpPoint[j]
	}
}

func initWorld() {
	for i := 0; i < HEIGHT; i++ {
		World[i][0] = 1
		World[i][WIDTH-1] = 1
	}
	for j := 0; j < WIDTH; j++ {
		World[HEIGHT-1][j] = 1
	}
}

func drawWorld(command string) {
	command = strings.Replace(command, "\r", "", -1)
	command = strings.Replace(command, "\n", "", -1)
	switch command {
	case "a":
		fmt.Println("move left")
		ActiveObject.ShiftLeft()
		break
	case "d":
		fmt.Println("move right")
		ActiveObject.ShiftRight()
		break
	case "w":
		fmt.Println("rotate cw")
		ActiveObject.RotateCW()
		break
	case "s":
		fmt.Println("rotate ccw")
		ActiveObject.RotateCCW()
		break

	}

	if ActiveObject.isFalling() {
		dropObject = false
	} else {
		dropObject = true

	}

	for _, val := range Objects {
		//draw all objects
		for _, v := range val.Body {
			World[v.Y][v.X] = 1

		}
	}

	for i := 0; i < WIDTH; i++ {
		res := ""
		for j := 0; j < HEIGHT; j++ {
			if World[i][j] == 1 {
				res += "*"
			} else {
				res += " "
			}
		}
		fmt.Println(res)
	}
}

func isGameOver() bool {
	for i := 1; i < WIDTH-1; i++ {
		if World[0][i] == 1 {
			return true
		}
	}
	return false
}

func main() {

	initWorld()

	i := 0

	for true {
		if dropObject {
			sp := Shape{}
			randomobj := rand.Intn(4)
			sp.Create(randomobj)
			Objects = append(Objects, sp)
			ActiveObject = sp
		}
		reader := bufio.NewReader(os.Stdin)
		if i == 0 {
			fmt.Print("press enter to start the game")
			i++
		} else {
			fmt.Print("Enter command: ")
		}

		cmd, _ := reader.ReadString('\n')

		if isGameOver() {
			fmt.Println("Game over")
			break

		}

		drawWorld(cmd)

	}
}
