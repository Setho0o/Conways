package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type Mode int

const (
  blocksize int = 10
  LayoutX int = ScreenX
  LayoutY int = ScreenY
  ScreenX int = 960
  ScreenY int = 540
Play Mode = iota
  Pause  
  Restart 
)

type Cell struct {
  alive bool 
  x int
  y int
  bx []int //bounds of the cell 
  by []int
}

func (c *Cell) Bounds() {
	for x := c.x; x < c.x+blocksize; x++ {
		c.bx = append(c.bx, x)
	}
	for y := c.y; y < c.y+blocksize; y++ {
		c.by = append(c.by, y)
	}
}

func (g *Game) FillBounds(c Cell) {
	for _, x := range c.bx {
		for _, y := range c.by {
			g.s.Set(x, y, color.White)
		}
	}
}

type Game struct{
  m Mode
  matrix [][]Cell
  s *ebiten.Image
}

func (g *Game) Update() error {
  g.Keys()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
  g.s = screen
  g.Life()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return LayoutX, LayoutY
}

func main() {
	ebiten.SetWindowSize(ScreenX,ScreenY)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(GameInit()); err != nil {
		log.Fatal(err)
	}
}

func GameInit() *Game {
  return &Game {
    m: Play,
    matrix: MatrixInit(),
  }
}
func MatrixInit() [][]Cell {
  m := make([][]Cell, LayoutY / blocksize)
  for i := range m {
    m[i] = make([]Cell,LayoutX / blocksize)
  }

  for i, row := range m {
		for j, _ := range row {
			m[i][j] = Cell {
        x: j * blocksize, 
        y: i * blocksize, 
        alive: false,
      }
      m[i][j].Bounds()
    }
	}
  return m
}
func (g *Game) Keys() {
  if ebiten.IsKeyPressed(ebiten.KeyQ) {
    os.Exit(0)
  } else if ebiten.IsKeyPressed(ebiten.KeyA) {
    if g.m == Play {
      g.m = Pause
    }
  } else if ebiten.IsKeyPressed(ebiten.KeyS) {
    if g.m == Pause {
      g.m = Play
    }
  } else if ebiten.IsKeyPressed(ebiten.KeyD) {
    g.m = Restart
    fmt.Println("restart")
  }
}
func (g *Game) Life() {
  for _, e := range g.matrix {
    for _, cell := range e {
      g.s.Set(cell.x,cell.y, color.White)
    }
  }
}







