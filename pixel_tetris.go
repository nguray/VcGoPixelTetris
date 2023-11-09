package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/golang/freetype/truetype"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"github.com/gopxl/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

const (
	LEFT       = 10
	TOP        = 10
	NB_ROWS    = 20
	NB_COLUMNS = 12
	WIN_WIDTH  = 480
	WIN_HEIGHT = 560
	TITLE      = "Go Pixel Tetris"
)

type GameMode int

const (
	STANDBY GameMode = iota
	PLAY
	GAMEPAUSE
	GAMEOVER
	HIGHSCORES
)

type HightScore struct {
	name  string
	score int
}

type Vector2i struct {
	x int32
	y int32
}

type Color struct {
	R, G, B, A uint8
}

var (
	tetrominos []Vector2i
	colors     []Color
)

type ProcessEvents_t func(win pixelgl.Window) bool

type IsOutLimit_t func() bool

type DrawMode_t func(win pixel.Target)

var (
	cellSize          int32
	myRand            *rand.Rand
	processEvents     ProcessEvents_t
	isOutLRBoardLimit IsOutLimit_t
	drawCurMode       DrawMode_t
	tt_font           font.Face
	atlas             *text.Atlas
	successBuffer     *beep.Buffer
	musicBuffer       *beep.Buffer
	musicCtrl         *beep.Ctrl
	musicVolume       *effects.Volume
	idtetrominosBag   int
	tetrominosBag     []int32
	game              *Game
	curTetromino      *Tetromino
	nextTetromino     *Tetromino
	startR            time.Time
)

func InitTetrominos() {

	tetrominos = []Vector2i{
		{0, 0}, {0, 0}, {0, 0}, {0, 0},
		{0, -1}, {0, 0}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 0}, {1, 0}, {1, 1},
		{0, -1}, {0, 0}, {0, 1}, {0, 2},
		{-1, 0}, {0, 0}, {1, 0}, {0, 1},
		{0, 0}, {1, 0}, {0, 1}, {1, 1},
		{-1, -1}, {0, -1}, {0, 0}, {0, 1},
		{1, -1}, {0, -1}, {0, 0}, {0, 1}}

	colors = []Color{
		{R: 0, G: 0, B: 0, A: 0xFF},
		{R: 0xFF, G: 0x60, B: 0x60, A: 0xFF},
		{R: 0x60, G: 0xFF, B: 0x60, A: 0xFF},
		{R: 0x60, G: 0x60, B: 0xFF, A: 0xFF},
		{R: 0xCC, G: 0xCC, B: 0x60, A: 0xFF},
		{R: 0xCC, G: 0x60, B: 0xCC, A: 0xFF},
		{R: 0x60, G: 0xCC, B: 0xCC, A: 0xFF},
		{R: 0xDA, G: 0xAA, B: 0x00, A: 0xFF}}

}

func TetrisRandomizer() int32 {

	var (
		iSrc int32
		ityp int32
	)

	if idtetrominosBag < 14 {
		ityp = tetrominosBag[idtetrominosBag]
		idtetrominosBag += 1
	} else {
		//-- Shuttle bag
		for i := 0; i < 14; i++ {
			iSrc = int32(myRand.Intn(14))
			ityp = tetrominosBag[iSrc]
			tetrominosBag[iSrc] = tetrominosBag[0]
			tetrominosBag[0] = ityp
		}
		ityp = tetrominosBag[0]
		idtetrominosBag = 1
	}

	return ityp
}

func NewTetromino() {
	//--------------------------------------------------
	curTetromino = nextTetromino
	curTetromino.x = 6 * cellSize
	//curTetromino.y = NB_ROWS*cellSize - cellSize
	curTetromino.y = WIN_HEIGHT - TOP + curTetromino.MaxY()*cellSize
	nextTetromino = TetrominoNew(TetrisRandomizer(), (NB_COLUMNS+3)*cellSize, 10*cellSize)

}

func ProcessEventsPlay(win pixelgl.Window) bool {

	if win.JustPressed(pixelgl.KeyP) {
		game.fPause = !game.fPause
	} else if win.JustPressed(pixelgl.KeyLeft) {
		game.velX = -1
		isOutLRBoardLimit = curTetromino.IsOutLeftBoardLimit
	} else if win.JustPressed(pixelgl.KeyRight) {
		game.velX = 1
		isOutLRBoardLimit = curTetromino.IsOutRightBoardLimit
	} else if win.JustPressed(pixelgl.KeyUp) {
		if curTetromino != nil {
			curTetromino.RotateLeft()
		}
	} else if win.JustPressed(pixelgl.KeyDown) {
		game.fFastDown = true
	} else if win.JustPressed(pixelgl.KeySpace) {
		if curTetromino != nil {
			//-- Drop current Tetromino
			game.fDrop = true
		}
	} else if win.JustPressed(pixelgl.KeyPause) {
		speaker.Lock()
		musicCtrl.Paused = !musicCtrl.Paused
		speaker.Unlock()
	} else if win.JustPressed(pixelgl.KeyKPAdd) {
		speaker.Lock()
		musicVolume.Volume += 0.5
		speaker.Unlock()
	} else if win.JustPressed(pixelgl.KeyKPSubtract) {
		speaker.Lock()
		musicVolume.Volume -= 0.5
		speaker.Unlock()
	} else if win.JustPressed(pixelgl.KeyEscape) {
		game.curMode = STANDBY
		processEvents = ProcessEventsStandBy
		drawCurMode = DrawStandByMode
		curTetromino = nil
		game.ClearBoard()
		return false
	}

	if win.JustReleased(pixelgl.KeyLeft) || win.JustReleased(pixelgl.KeyRight) {
		game.velX = 0
		isOutLRBoardLimit = curTetromino.IsAlwaysOutBoardLimit
	} else if win.JustReleased(pixelgl.KeyDown) {
		game.fFastDown = false
	}

	return true
}

func ProcessEventsStandBy(win pixelgl.Window) bool {

	if win.JustPressed(pixelgl.KeySpace) {
		game.curMode = PLAY
		processEvents = ProcessEventsPlay
		drawCurMode = DrawPlayMode
		NewTetromino()
		game.curScore = 0
	} else if win.JustPressed(pixelgl.KeyPause) {
		speaker.Lock()
		musicCtrl.Paused = !musicCtrl.Paused
		speaker.Unlock()
	} else if win.JustPressed(pixelgl.KeyKPAdd) {
		speaker.Lock()
		musicVolume.Volume += 0.5
		speaker.Unlock()
	} else if win.JustPressed(pixelgl.KeyKPSubtract) {
		speaker.Lock()
		musicVolume.Volume -= 0.5
		speaker.Unlock()
	} else if win.JustPressed(pixelgl.KeyEscape) {
		game.fQuitGame = true
	}
	return true
}

func ProcessEventsGameOver(win pixelgl.Window) bool {

	if win.JustPressed(pixelgl.KeySpace) {
		game.curMode = STANDBY
		processEvents = ProcessEventsStandBy
		drawCurMode = DrawStandByMode
		curTetromino = nil
		game.ClearBoard()
	} else if win.JustPressed(pixelgl.KeyPause) {
		speaker.Lock()
		musicCtrl.Paused = !musicCtrl.Paused
		speaker.Unlock()
	} else if win.JustPressed(pixelgl.KeyKPAdd) {
		speaker.Lock()
		musicVolume.Volume += 0.5
		speaker.Unlock()
	} else if win.JustPressed(pixelgl.KeyKPSubtract) {
		speaker.Lock()
		musicVolume.Volume -= 0.5
		speaker.Unlock()

	}
	return true
}

func ProcessEventsHightScores(win pixelgl.Window) bool {

	if win.JustPressed(pixelgl.KeyEnter) || win.JustPressed(pixelgl.KeyKPEnter) {
		game.SaveHighScores("HighScores.txt")
		game.curMode = STANDBY
		processEvents = ProcessEventsStandBy
		drawCurMode = DrawStandByMode
	} else if win.JustPressed(pixelgl.KeyPause) {
		speaker.Lock()
		musicCtrl.Paused = !musicCtrl.Paused
		speaker.Unlock()
	} else if win.JustPressed(pixelgl.KeyKPAdd) {
		speaker.Lock()
		musicVolume.Volume += 0.5
		speaker.Unlock()
	} else if win.JustPressed(pixelgl.KeyKPSubtract) {
		speaker.Lock()
		musicVolume.Volume -= 0.5
		speaker.Unlock()
	} else if win.JustPressed(pixelgl.KeyBackspace) {
		sz := len(game.userName)
		if sz > 0 {
			game.userName = game.userName[:sz-1]
			game.highScores[game.idHighScore].name = game.userName
		}
	} else if win.JustPressed(pixelgl.KeyEscape) {
		if len(game.userName) == 0 && game.idHighScore >= 0 {
			game.highScores[game.idHighScore].name = "XXXXXX"
		}
		game.SaveHighScores("HighScores.txt")
		game.curMode = STANDBY
		processEvents = ProcessEventsStandBy
		drawCurMode = DrawStandByMode
	} else {
		for _, k := range game.tblKeyChars {
			if win.JustPressed(k.keycode) {
				if game.idHighScore >= 0 {
					if len(game.userName) < 10 {
						game.userName += k.c
						game.highScores[game.idHighScore].name = game.userName
					}
				}
			}
		}
	}

	return true
}

func loadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

func DrawPlayMode(win pixel.Target) {

	if curTetromino != nil {
		curTetromino.Draw(win)
	}

}

func DrawStandByMode(win pixel.Target) {

	ox := float64(LEFT + (NB_COLUMNS/2)*cellSize)
	oy := float64(WIN_HEIGHT - TOP - 7*cellSize)
	txt := text.New(pixel.V(ox, oy), atlas)
	txt.Color = colornames.Gold
	rect := pixel.R(LEFT, oy, float64(LEFT+NB_COLUMNS*cellSize), oy+float64(cellSize))
	fmt.Fprintf(txt, "TETRIS in Golang")
	txt.Draw(win, pixel.IM.Moved(rect.Bounds().Center().Sub(txt.Bounds().Center())))

	oy -= float64(2*cellSize + 4)
	txt = text.New(pixel.V(ox, oy), atlas)
	txt.Color = colornames.Gold
	rect = pixel.R(LEFT, oy, float64(LEFT+NB_COLUMNS*cellSize), oy+float64(cellSize))
	fmt.Fprintf(txt, "Press SPACE to PLAY")
	txt.Draw(win, pixel.IM.Moved(rect.Bounds().Center().Sub(txt.Bounds().Center())))

}

func DrawGameOverMode(win pixel.Target) {

	oy := float64(WIN_HEIGHT - TOP - 7*cellSize)
	ox := float64(LEFT + (NB_COLUMNS/2)*cellSize)
	txt := text.New(pixel.V(ox, oy), atlas)
	txt.Color = colornames.Gold
	rect := pixel.R(LEFT, oy, float64(LEFT+NB_COLUMNS*cellSize), oy+float64(cellSize))
	fmt.Fprintf(txt, "GAME OVER")
	txt.Draw(win, pixel.IM.Moved(rect.Bounds().Center().Sub(txt.Bounds().Center())))

	oy -= float64(2*cellSize + 4)
	txt = text.New(pixel.V(ox, oy), atlas)
	txt.Color = colornames.Gold
	rect = pixel.R(LEFT, oy, float64(LEFT+NB_COLUMNS*cellSize), oy+float64(cellSize))
	fmt.Fprintf(txt, "Press SPACE to Continue")
	txt.Draw(win, pixel.IM.Moved(rect.Bounds().Center().Sub(txt.Bounds().Center())))

}

func DrawHighScoresMode(win pixel.Target) {

	oy := float64(WIN_HEIGHT - TOP - 2*cellSize)
	ox := float64(LEFT + (NB_COLUMNS/2)*cellSize)
	txt := text.New(pixel.V(ox, oy), atlas)
	txt.Color = colornames.Gold
	rect := pixel.R(LEFT, oy, float64(LEFT+NB_COLUMNS*cellSize), oy+float64(cellSize))
	fmt.Fprintf(txt, "HIGH SCORES")
	txt.Draw(win, pixel.IM.Moved(rect.Bounds().Center().Sub(txt.Bounds().Center())))

	x1 := float64(LEFT + 4)
	x2 := float64(LEFT + NB_COLUMNS*cellSize/2)
	x3 := float64(LEFT + NB_COLUMNS*cellSize - 4)
	for i, h := range game.highScores {
		lineColor := colornames.Gold
		if game.idHighScore == i && game.iColorHighScore%2 != 0 {
			lineColor = colornames.Orange
		}
		oy -= float64(cellSize + 4)
		txt = text.New(pixel.V(ox, oy), atlas)
		txt.Color = lineColor
		rect = pixel.R(x1, oy, x2-4, oy)
		fmt.Fprintf(txt, "%s", h.name)
		txt.Draw(win, pixel.IM.Moved(rect.Bounds().Center().Sub(txt.Bounds().Center())))
		txt = text.New(pixel.V(ox, oy), atlas)
		txt.Color = lineColor
		rect = pixel.R(x2+4, oy, x3, oy)
		fmt.Fprintf(txt, "%06d", h.score)
		txt.Draw(win, pixel.IM.Moved(rect.Bounds().Center().Sub(txt.Bounds().Center())))

	}

	elapsedR := time.Since(startR)
	if elapsedR.Milliseconds() > 500 {
		startR = time.Now()
		game.iColorHighScore += 1
	}

}

func PlaySuccesSound() {
	//-----------------------------------------
	shot := successBuffer.Streamer(0, successBuffer.Len())
	volume := &effects.Volume{Streamer: shot, Base: 2, Volume: -2}
	speaker.Play(volume)

}

func run() {

	cfg := pixelgl.WindowConfig{
		Title:  TITLE,
		Bounds: pixel.R(0, 0, WIN_WIDTH, WIN_HEIGHT),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	successfile, err := os.Open("109662__grunz__success.wav")
	if err != nil {
		log.Fatal(err)
	}
	successStreamer, successFormat, err := wav.Decode(successfile)
	if err != nil {
		log.Fatal(err)
	}
	defer successStreamer.Close()

	musicfile, err := os.Open("Tetris.wav")
	if err != nil {
		log.Fatal(err)
	}

	musicStreamer, musicFormat, err := wav.Decode(musicfile)
	if err != nil {
		log.Fatal(err)
	}
	defer musicStreamer.Close()

	successBuffer = beep.NewBuffer(successFormat)
	successBuffer.Append(successStreamer)
	successStreamer.Close()

	musicBuffer = beep.NewBuffer(musicFormat)
	musicBuffer.Append(musicStreamer)
	musicStreamer.Close()

	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))

	musicTetris := musicBuffer.Streamer(0, musicBuffer.Len())
	musicCtrl = &beep.Ctrl{Streamer: beep.Loop(-1, musicTetris), Paused: false}
	musicVolume = &effects.Volume{Streamer: musicCtrl, Base: 2, Volume: -3}
	speaker.Play(musicVolume)

	tt_font, err = loadTTF("Sansation_Bold.ttf", 20)
	if err != nil {
		panic(err)
	}
	//--
	tetrominosBag = make([]int32, 14)
	tetrominosBag = []int32{
		1, 2, 3, 4, 5, 6, 7, 1, 2, 3, 4, 5, 6, 7,
	}
	idtetrominosBag = 14

	cellSize = int32(WIN_WIDTH / (NB_COLUMNS + 7))

	InitTetrominos()
	myRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	game = GameNew()
	curTetromino = nil
	nextTetromino = TetrominoNew(TetrisRandomizer(), (NB_COLUMNS+3)*cellSize, 10*cellSize)
	game.LoadHighScores("HighScores.txt")

	atlas = text.NewAtlas(tt_font, text.ASCII)
	txt := text.New(pixel.V(10, 20), atlas)
	txt.Color = colornames.Gold
	fmt.Fprintf(txt, "SCORE : %06d", 100)

	startV := time.Now()
	startH := startV
	startR = startV

	game.curMode = STANDBY
	processEvents = ProcessEventsStandBy
	drawCurMode = DrawStandByMode

	isOutLRBoardLimit = curTetromino.IsAlwaysOutBoardLimit

	for !win.Closed() {

		if !processEvents(*win) {
			//-- Manage Escape from PLAY mode
			if game.curScore != 0 {
				id := game.IsHightScore(game.curScore)
				//-- Manage Game Over and User Escape
				if id >= 0 {
					//--
					game.InsertHightScore(id, game.userName, game.curScore)
					game.curMode = HIGHSCORES
					processEvents = ProcessEventsHightScores
					drawCurMode = DrawHighScoresMode
					game.ClearBoard()
					curTetromino = nil
				} else {
					//--
					game.ClearBoard()
					curTetromino = nil
					game.curMode = STANDBY
					processEvents = ProcessEventsStandBy
					drawCurMode = DrawStandByMode
				}

			}

		}

		if game.fQuitGame {
			break
		}

		if game.curMode == PLAY {
			//-- Update game state

			elapsedV := time.Since(startV)
			elapsedR := time.Since(startR)

			if game.nbCompledLines > 0 {
				//-- Remove Completed lines
				if elapsedV.Milliseconds() > 250 {
					startV = time.Now()
					game.nbCompledLines--
					game.EraseFirstCompletedLine()
					PlaySuccesSound()
				}
			} else if game.horizontalMove != 0 {
				//-- Move to the next slot
				elapsed := time.Since(startH)
				if elapsed.Milliseconds() > 20 {
					startH = time.Now()

					for iOffSet := 0; iOffSet < int(4); iOffSet++ {

						backupX := curTetromino.x
						curTetromino.x += game.horizontalMove

						if game.horizontalMove < 0 {
							isOutLRBoardLimit = curTetromino.IsOutLeftBoardLimit
						} else {
							// game.horizontalMove > 0
							isOutLRBoardLimit = curTetromino.IsOutRightBoardLimit
						}

						if isOutLRBoardLimit() {
							curTetromino.x = backupX
							game.horizontalMove = 0
							break
						} else {
							if curTetromino.HitGround(game.board) {
								curTetromino.x = backupX
								game.horizontalMove = 0
								break
							}
						}

						if game.horizontalMove != 0 {
							if game.horizontalStartColumn != curTetromino.Column() {
								curTetromino.x = backupX
								game.horizontalMove = 0
								startH = time.Now()
								break
							}
						}

					}
				}

			} else if game.fDrop {
				//-- Drop Tetromino
				if elapsedV.Milliseconds() > 10 {
					startV = time.Now()
					for iOffSet := 0; iOffSet < 6; iOffSet++ {
						//-- Move down to check
						curTetromino.y--
						if curTetromino.HitGround(game.board) {
							curTetromino.y++
							game.FreezeTetromino(curTetromino)
							NewTetromino()
							game.fDrop = false
						} else if curTetromino.IsOutBottomLimit() {
							curTetromino.y++
							game.FreezeTetromino(curTetromino)
							NewTetromino()
							game.fDrop = false
						}
						if game.fDrop {
							if game.velX != 0 {
								elapsed := time.Since(startH)

								if elapsed.Milliseconds() > 20 {

									backupX := curTetromino.x
									curTetromino.x += game.velX

									if isOutLRBoardLimit() {
										curTetromino.x = backupX
									} else {
										if curTetromino.HitGround(game.board) {
											curTetromino.x = backupX
										} else {
											startH = time.Now()
											game.horizontalMove = game.velX
											game.horizontalStartColumn = curTetromino.Column()
											break
										}
									}

								}

							}

						}

					}
				}

			} else {
				//-- Move down Tetromino
				var limitElapse int64 = 25
				if game.fFastDown {
					limitElapse = 10
				}

				if elapsedV.Milliseconds() > limitElapse {
					startV = time.Now()
					if curTetromino.HitGround(game.board) {
						fmt.Println("-----")
					}
					for iOffSet := 0; iOffSet < 3; iOffSet++ {
						//-- Move down to check
						curTetromino.y--
						fMove := true
						if curTetromino.HitGround(game.board) {
							curTetromino.y++
							game.FreezeTetromino(curTetromino)
							NewTetromino()
							fMove = false

						} else if curTetromino.IsOutBottomLimit() {
							curTetromino.y++
							game.FreezeTetromino(curTetromino)
							NewTetromino()
							fMove = false
						}
						if fMove {
							if game.velX != 0 {
								elapsed := time.Since(startH)
								if elapsed.Milliseconds() > 15 {

									backupX := curTetromino.x
									curTetromino.x += game.velX

									if isOutLRBoardLimit() {
										curTetromino.x = backupX
									} else {
										if curTetromino.HitGround(game.board) {
											curTetromino.x = backupX
										} else {
											startH = time.Now()
											game.horizontalMove = game.velX
											game.horizontalStartColumn = curTetromino.Column()
											break
										}
									}

								}
							}
						}
					}
				}
			}

			//-- Check Game Over
			if game.IsGameOver() {

				//--
				id := game.IsHightScore(game.curScore)

				if id >= 0 {
					//--
					game.InsertHightScore(id, game.userName, game.curScore)
					game.curMode = HIGHSCORES
					processEvents = ProcessEventsHightScores
					drawCurMode = DrawHighScoresMode
					game.ClearBoard()
					curTetromino = nil
				} else {
					//--
					game.curMode = GAMEOVER
					processEvents = ProcessEventsGameOver
					drawCurMode = DrawGameOverMode
					game.ClearBoard()
					curTetromino = nil

				}

			}

			if elapsedR.Milliseconds() > 500 {
				startR = time.Now()
				nextTetromino.RotateRight()
			}

		}

		win.Clear(colornames.Darkblue)

		game.DrawBoard(win)

		if nextTetromino != nil {
			nextTetromino.Draw(win)
		}

		drawCurMode(win)

		//-- Draw current score
		txt := text.New(pixel.V(10, 20), atlas)
		txt.Color = colornames.Gold
		fmt.Fprintf(txt, "SCORE : %06d", game.curScore)
		txt.Draw(win, pixel.IM)

		win.Update()

	}

}

func main() {
	pixelgl.Run(run)
}
