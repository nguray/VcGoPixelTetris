package main

import (
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
)

type Tetromino struct {
	typ int32
	x   int32
	y   int32
	v   [4]Vector2i
}

func TetrominoNew(typ, x, y int32) *Tetromino {

	//--
	t := &Tetromino{typ: typ, x: x, y: y}
	t.InitGfx()
	return t
}

func (te *Tetromino) Draw(win pixel.Target) {

	var (
		x float64
		y float64
	)

	imd1 := imdraw.New(nil)
	c := colors[te.typ]
	imd1.Color = pixel.RGB(float64(c.R)/255.0, float64(c.G)/255.0, float64(c.B)/255.0)
	offsetV := WIN_HEIGHT - TOP - NB_ROWS*cellSize
	d := float64(cellSize - 2)
	for _, v := range te.v {
		x = float64(v.x*cellSize+te.x) + LEFT + 1
		y = float64(v.y*cellSize+te.y) + float64(offsetV) - 1
		imd1.Push(pixel.V(x, y))
		imd1.Push(pixel.V(x+d, y))
		imd1.Push(pixel.V(x+d, y-d))
		imd1.Push(pixel.V(x, y-d))
		imd1.Polygon(0)
	}

	// y = float64(te.MinY()*cellSize + te.y - cellSize + offsetV)
	// imd1.Push(pixel.V(10, y))
	// imd1.Push(pixel.V(float64(WIN_WIDTH-10), y))
	// imd1.Line(1)

	imd1.Draw(win)

}

func (te *Tetromino) InitGfx() {

	offSet := int(te.typ) * len(te.v)
	for i := 0; i < len(te.v); i++ {
		te.v[i] = tetrominos[i+offSet]
	}

}

func (te *Tetromino) RotateLeft() {
	if te.typ != 5 {
		var x, y int32
		for i := 0; i < len(te.v); i++ {
			x = -te.v[i].y
			y = te.v[i].x
			te.v[i].x = x
			te.v[i].y = y
		}
	}
}

func (te *Tetromino) RotateRight() {
	if te.typ != 5 {
		var x, y int32
		for i := 0; i < len(te.v); i++ {
			x = te.v[i].y
			y = -te.v[i].x
			te.v[i].x = x
			te.v[i].y = y
		}
	}
}

func (te *Tetromino) MinX() int32 {
	var (
		x    int32
		minX int32
	)
	minX = te.v[0].x
	for i := 1; i < len(te.v); i++ {
		x = te.v[i].x
		if x < minX {
			minX = x
		}
	}
	return minX
}

func (te *Tetromino) MaxX() int32 {
	var (
		x    int32
		maxX int32
	)
	maxX = te.v[0].x
	for i := 1; i < len(te.v); i++ {
		x = te.v[i].x
		if x > maxX {
			maxX = x
		}
	}
	return maxX
}

func (te *Tetromino) MaxY() int32 {
	var (
		y int32
	)
	maxY := te.v[0].y
	for i := 1; i < len(te.v); i++ {
		y = te.v[i].y
		if y > maxY {
			maxY = y
		}
	}
	return maxY
}

func (te *Tetromino) MinY() int32 {
	var (
		y int32
	)
	minY := te.v[0].y
	for i := 1; i < len(te.v); i++ {
		y = te.v[i].y
		if y < minY {
			minY = y
		}
	}
	return minY
}

func (te *Tetromino) Column() int32 {
	return int32(te.x / cellSize)
}

func (te *Tetromino) IsOutLeftBoardLimit() bool {
	l := te.MinX()*cellSize + te.x
	return l < 0
}

func (te *Tetromino) IsOutRightBoardLimit() bool {
	r := te.MaxX()*cellSize + cellSize + te.x
	return r > NB_COLUMNS*cellSize
}

func (te *Tetromino) IsAlwaysOutBoardLimit() bool {
	return true
}

func (te *Tetromino) IsOutBottomLimit() bool {
	//--------------------------------------------------
	y := float64(te.MinY()*cellSize + te.y - cellSize)
	return y <= 0
}

func (te *Tetromino) HitGround(board []int) bool {

	//--------------------------------------------------
	Hit := func(x int32, y int32) bool {
		ix := int32(x / cellSize)
		iy := int32(y / cellSize)
		//fmt.Println(ix, iy)
		if (ix >= 0) && ix < NB_COLUMNS && (iy >= 0) && (iy < NB_ROWS) {
			v := board[iy*NB_COLUMNS+ix]
			if v != 0 {
				return true
			}
		}
		return false
	}

	offSet := NB_ROWS * cellSize

	for _, v := range te.v {

		x := v.x*cellSize + te.x + 1
		y := offSet - (v.y*cellSize + te.y) + 1
		if Hit(x, y) {
			return true
		}

		x = v.x*cellSize + cellSize - 1 + te.x
		y = offSet - (v.y*cellSize + te.y) + 1
		if Hit(x, y) {
			return true
		}

		x = v.x*cellSize + cellSize - 1 + te.x
		y = offSet - (v.y*cellSize + te.y - cellSize + 1)
		if Hit(x, y) {
			return true
		}

		x = v.x*cellSize + te.x + 1
		y = offSet - (v.y*cellSize + te.y - cellSize + 1)
		if Hit(x, y) {
			return true
		}

	}

	return false
}
