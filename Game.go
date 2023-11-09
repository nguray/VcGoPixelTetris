package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"
)

type KeyChar struct {
	keycode pixelgl.Button
	c       string
}

type Game struct {
	velX                  int32
	fDrop                 bool
	fFastDown             bool
	curMode               GameMode
	curScore              int
	board                 []int
	highScores            []HightScore
	idHighScore           int
	userName              string
	tblKeyChars           []KeyChar
	fQuitGame             bool
	horizontalMove        int32
	horizontalStartColumn int32
	fPause                bool
	nbCompledLines        int
	iColorHighScore       int
}

func GameNew() *Game { //int32(myRand.Intn(7)+1)
	game := &Game{0, false, false, STANDBY, 0, make([]int, NB_ROWS*NB_COLUMNS),
		make([]HightScore, 10), -1, "", make([]KeyChar, 1), false, 0, 0, false, 0, 0}
	for i := 0; i < len(game.highScores); i++ {
		game.highScores[i] = HightScore{"--------", 0}
	}

	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyA, c: "A"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyB, c: "B"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyC, c: "C"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyD, c: "D"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyE, c: "E"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyF, c: "F"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyG, c: "G"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyH, c: "H"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyI, c: "I"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyJ, c: "J"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyK, c: "K"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyL, c: "L"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyM, c: "M"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyN, c: "N"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyO, c: "O"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyP, c: "P"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyQ, c: "Q"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyR, c: "R"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyS, c: "S"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyT, c: "T"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyU, c: "U"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyV, c: "V"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyW, c: "W"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyX, c: "X"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyY, c: "Y"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyZ, c: "Z"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.Key0, c: "0"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.Key1, c: "1"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.Key2, c: "2"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.Key3, c: "3"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.Key4, c: "4"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.Key5, c: "5"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.Key6, c: "6"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.Key7, c: "7"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.Key8, c: "8"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.Key9, c: "9"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyKP0, c: "0"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyKP1, c: "1"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyKP2, c: "2"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyKP3, c: "3"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyKP4, c: "4"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyKP5, c: "5"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyKP6, c: "6"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyKP7, c: "7"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyKP8, c: "8"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: pixelgl.KeyKP9, c: "9"})

	// iy := 15 * NB_COLUMNS
	// for i := 0; i < NB_COLUMNS; i++ {
	// 	game.board[i+iy] = 2
	// 	game.board[i+iy+NB_COLUMNS] = 4

	// }

	// iy := 2 * NB_COLUMNS
	// game.board[iy+5] = 3
	// iy += NB_COLUMNS
	// game.board[iy+5] = 3
	// iy += 2*NB_COLUMNS
	// game.board[iy+5] = 3
	// iy += NB_COLUMNS
	// game.board[iy+5] = 3

	return game
}

func (ga *Game) SaveHighScores(fileName string) {
	//------------------------------------------------------

	var (
		str1 string
	)

	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for i := 0; i < 10; i++ {
		h := ga.highScores[i]
		if h.name == "" {
			h.name = "XXXX"
		}
		str1 = fmt.Sprintf("%s %d\n", h.name, h.score)
		_, _ = f.WriteString(str1)

	}

}

func (ga *Game) LoadHighScores(fileName string) {
	//------------------------------------------------------

	f, err := os.Open(fileName)

	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for nbL := 0; nbL < 10 && scanner.Scan(); nbL++ {

		//--
		strLineVal := scanner.Text()

		wordBreakDown := strings.Fields(strLineVal)

		ga.highScores[nbL].name = wordBreakDown[0]
		val, _ := strconv.ParseInt(wordBreakDown[1], 10, 32)
		ga.highScores[nbL].score = int(val)

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func (ga *Game) DrawBoard(win pixel.Target) {
	//----------------------------------------------------------------
	var (
		x    float64
		y    float64
		l, c int32
	)
	ga.DrawBackground(win)
	a := float64(cellSize - 2)
	imd1 := imdraw.New(nil)
	offsetV := float64(WIN_HEIGHT - TOP)
	for l = 0; l < NB_ROWS; l++ {
		for c = 0; c < NB_COLUMNS; c++ {
			v := ga.board[l*NB_COLUMNS+c]
			if v != 0 {
				x = float64(c*cellSize) + float64(LEFT) + 1
				y = -float64(cellSize*l) + offsetV - 1
				c := colors[v]
				imd1.Color = pixel.RGB(float64(c.R)/255.0, float64(c.G)/255.0, float64(c.B)/255.0)
				imd1.Push(pixel.V(x, y))
				imd1.Push(pixel.V(x+a, y))
				imd1.Push(pixel.V(x+a, y-a))
				imd1.Push(pixel.V(x, y-a))
				imd1.Polygon(0)
			}
		}
	}
	imd1.Draw(win)

}

func (ga *Game) DrawBackground(win pixel.Target) {

	var (
		left, top, right, bottom float64
	)

	imd := imdraw.New(nil)

	left = LEFT
	top = WIN_HEIGHT - TOP
	right = left + float64(NB_COLUMNS*cellSize)
	bottom = top - float64(NB_ROWS*cellSize)
	imd.Color = pixel.RGB(10.0/255.0, 10.0/255.0, 100.0/255.0)
	imd.Push(pixel.V(left, top))
	imd.Push(pixel.V(right, top))
	imd.Push(pixel.V(right, bottom))
	imd.Push(pixel.V(left, bottom))
	imd.Polygon(0)

	imd.Draw(win)
}

func (ga *Game) FreezeTetromino(tetro *Tetromino) {
	//--------------------------------------------------
	if tetro != nil {
		offSet := NB_ROWS * cellSize
		ix := int32((tetro.x + 1) / cellSize)
		iy := int32((offSet - tetro.y + 1) / cellSize)
		for _, v := range tetro.v {
			x := v.x + ix
			y := iy - v.y
			if x >= 0 && x < NB_COLUMNS && y >= 0 && y < NB_ROWS {
				ga.board[y*NB_COLUMNS+x] = int(tetro.typ)
			}
		}
		//--
		ga.nbCompledLines = ga.ComputeCompletedLines()
		if ga.nbCompledLines > 0 {
			ga.curScore += ga.ComputeScore(ga.nbCompledLines)
		}

	}
}

func (ga *Game) ComputeCompletedLines() int {
	//--------------------------------------------------
	nbLines := 0
	fCompleted := false
	for r := 0; r < NB_ROWS; r++ {
		fCompleted = true
		for c := 0; c < NB_COLUMNS; c++ {
			if ga.board[r*NB_COLUMNS+c] == 0 {
				fCompleted = false
				break
			}
		}
		if fCompleted {
			nbLines++
		}
	}
	//fmt.Println("Nbre Erased Lines ", nbLines)
	return nbLines
}

func (ga *Game) EraseFirstCompletedLine() {
	//--------------------------------------------------
	fCompleted := false
	for r := 0; r < NB_ROWS; r++ {
		fCompleted = true
		for c := 0; c < NB_COLUMNS; c++ {
			if ga.board[r*NB_COLUMNS+c] == 0 {
				fCompleted = false
				break
			}
		}
		if fCompleted {
			//-- Décaler d'une ligne le plateau
			for r1 := r; r1 > 0; r1-- {
				for c1 := 0; c1 < NB_COLUMNS; c1++ {
					ga.board[r1*NB_COLUMNS+c1] = ga.board[(r1-1)*NB_COLUMNS+c1]
				}
			}
			return
		}
	}
	//fmt.Println("Nbre Erased Lines ", nbLines)
}

func (ga *Game) ClearBoard() {
	//--------------------------------------------------
	for i := 0; i < NB_ROWS*NB_COLUMNS; i++ {
		ga.board[i] = 0
	}

}

func (ga *Game) IsGameOver() bool {
	//------------------------------------------------------
	for i := 0; i < NB_COLUMNS; i++ {
		if ga.board[i] != 0 {
			return true
		}
	}
	return false
}

func (ga *Game) IsHightScore(newscore int) int {
	//--------------------------------------------------
	for i, v := range ga.highScores {
		if newscore > v.score {
			return i
		}
	}
	return -1
}

func (ga *Game) InsertHightScore(id int, name string, score int) {
	//--------------------------------------------------
	ga.highScores = append(ga.highScores[:id+1], ga.highScores[id:]...)
	ga.highScores[id] = HightScore{name: name, score: score}
	ga.idHighScore = id
	ga.userName = name

}

func (ga *Game) ComputeScore(nbLines int) int {
	var score int
	switch nbLines {
	case 0:
		score = 0
	case 1:
		score = 40
	case 2:
		score = 100
	case 3:
		score = 300
	case 4:
		score = 1400
	default:
		score = 3000
	}
	return score
}
