package main

import (
	"image/color"
	"log"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type Mode int
type Direction int

const (
	randamt   int  = 20
	blocksize int  = 6
	LayoutX   int  = ScreenX
	LayoutY   int  = ScreenY
	ScreenX   int  = 1920
	ScreenY   int  = 1080
	Play      Mode = iota
	Pause
	TopLeft Direction = iota
	TopMiddle
	TopRight
	Left
	Right
	BottomLeft
	BottomMiddle
	BottomRight
)

var Dir = []Direction{
	TopLeft, TopMiddle, TopRight, Left, Right, BottomLeft, BottomMiddle, BottomRight,
}

type Cell struct {
	alive bool
	x     int
	y     int
	ix    int //index
	iy    int
}

func (g *Game) Fill(c Cell, color color.Color) {
	for i := range blocksize {
		for j := range blocksize {
			g.s.Set(c.x+i, c.y+j, color)
		}
	}
}

type Game struct {
	m      Mode
	matrix [][]*Cell
	s      *ebiten.Image
}

func (g *Game) Update() error {
	g.Keys()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.s = screen
	if g.m == Play {
		g.Life()
	} else {
		for _, e := range g.matrix {
			for _, cell := range e {
				if cell.alive {
					g.Fill(*cell, color.White)
				} else {
					g.Fill(*cell, color.Black)
				}
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return LayoutX, LayoutY
}

func main() {
	ebiten.SetWindowSize(ScreenX, ScreenY)
	game := GameInit()
	game.LifeRoulete()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func GameInit() *Game {
	return &Game{
		m:      Play,
		matrix: MatrixInit(),
	}
}

func MatrixInit() [][]*Cell {
	m := make([][]*Cell, LayoutY/blocksize)
	for i := range m {
		m[i] = make([]*Cell, LayoutX/blocksize)
	}

	for i, row := range m {
		for j, _ := range row {
			m[i][j] = &Cell{
				x:     j * blocksize,
				y:     i * blocksize,
				ix:    j,
				iy:    i,
				alive: false,
			}
		}
	}
	return m
}

func (g *Game) LifeRoulete() {
	for _, e := range g.matrix {
		for _, cell := range e {
			if rand.Intn(randamt) == 0 {
				cell.alive = true
			}
		}
	}
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
		for _, e := range g.matrix {
			for _, cell := range e {
				cell.alive = false
			}
		}
		g.LifeRoulete()
	}
}
func (g *Game) Life() {
	for _, e := range g.matrix {
		for _, cell := range e {
			if OutOfBounds(*cell) {
				continue
			}
			g.ConwaysRules(*cell)
			if cell.alive {
				g.Fill(*cell, color.White)
			} else {
				g.Fill(*cell, color.Black)
			}
		}
	}
}

func (g *Game) ConwaysRules(c Cell) {
	var life int
	for _, e := range Dir {
		if g.GetCellLife(c, e) {
			life++
		}
	}
	if c.alive {
		if life < 2 { // rule 1: any cell with fewer than two alive neighbors dies
			g.matrix[c.iy][c.ix].alive = false
		} else if life == 2 || life == 3 { // rule 2: any cell with two or three neighbors lives on
			g.matrix[c.iy][c.ix].alive = true
		} else if life > 3 { // rule 3: any cell with more than three live neighbors dies
			g.matrix[c.iy][c.ix].alive = false
		}
	} else if life == 3 { // rule 4: any dead cell with exactly three live neighbors becomes a live cell
		g.matrix[c.iy][c.ix].alive = true
	}
}

func (g *Game) GetCellLife(c Cell, d Direction) bool {
	var life bool
	switch d {
	case TopLeft:
		life = g.matrix[c.iy-1][c.ix-1].alive
	case TopMiddle:
		life = g.matrix[c.iy-1][c.ix].alive
	case TopRight:
		life = g.matrix[c.iy-1][c.ix+1].alive
	case Left:
		life = g.matrix[c.iy][c.ix-1].alive
	case Right:
		life = g.matrix[c.iy][c.ix+1].alive
	case BottomLeft:
		life = g.matrix[c.iy+1][c.ix-1].alive
	case BottomMiddle:
		life = g.matrix[c.iy+1][c.ix].alive
	case BottomRight:
		life = g.matrix[c.iy+1][c.ix+1].alive
	}
	return life
}

func OutOfBounds(c Cell) bool {
	if c.y < 1 || c.y > LayoutY {
		return true
	}
	if c.x < 1 || c.x > LayoutX {
		return true
	}
	if c.iy < 0 || c.iy > LayoutY/blocksize-2 {
		return true
	}
	if c.ix < 0 || c.ix > LayoutX/blocksize-2 {
		return true
	}
	return false
}
