package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const PaddleSymbol = 0x2588
const BallSymbol = 0x25CF

const PaddleHeight = 4
const InitialBallVelocityRow = 1
// a square on the screen is 2 times taller than its width
// so ball will be slower horizontally. To make even speed,
// we make horizontal ball velocity 2
const InitialBallVelocityCol = 2

type GameObject struct {
	row, col, width, height int
	velRow, velCol			int
	symbol					rune
}

var screen tcell.Screen
var player1Paddle *GameObject
var player2Paddle *GameObject
var ball *GameObject
var isGamePaused bool
var debugLog string

var gameObjects []*GameObject

type Boundary int

const (
	Top Boundary = iota
	Bottom
)

func main() {
	InitScreen()
	InitGameState()
	inputChan := InitUserInput()

	DrawState()

	for {
		HandleUserInput(ReadInput(inputChan))
		DrawState()
		UpdateState()

		time.Sleep(75 * time.Millisecond)
	}
}

func CollidesWithPaddle(ball *GameObject, paddle *GameObject) bool {
	var collidesOnColumn bool
	if ball.col < paddle.col {
		collidesOnColumn = ball.col + ball.velCol >= paddle.col
	} else {
		collidesOnColumn = ball.col + ball.velCol <= paddle.col
	}
	return 	collidesOnColumn &&
			ball.row >= paddle.row &&
			ball.row < paddle.row + paddle.height
}

func CollidesWithWall(obj *GameObject) bool {
	_, screenHeight := screen.Size()
	isWithinBoundaries := obj.row + obj.velRow >= 0 && obj.row + obj.velRow < screenHeight

	if isWithinBoundaries {

		return false
	}

	return true
}

func UpdateState() {
	if isGamePaused {
		return
	}

	// movement of paddles and ball
	for i := range gameObjects {
		gameObjects[i].row += gameObjects[i].velRow
		gameObjects[i].col += gameObjects[i].velCol
	}

	if CollidesWithPaddle(ball, player1Paddle) || CollidesWithPaddle(ball, player2Paddle) {
		ball.velCol = -ball.velCol
	}

	if CollidesWithWall(ball) {
		ball.velRow = -ball.velRow
	}
}

func InitScreen() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
}

func HandleUserInput(key string) {

	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	} else if key == "Rune[w]" && IsWithinBoundaries(player1Paddle, Top) {
		player1Paddle.row--
	} else if key == "Rune[s]" && IsWithinBoundaries(player1Paddle, Bottom) {
		player1Paddle.row++
	} else if key == "Up" && IsWithinBoundaries(player2Paddle, Top) {
		player2Paddle.row--
	} else if key == "Down" && IsWithinBoundaries(player2Paddle, Bottom) {
		player2Paddle.row++
	} else if key == "Rune[p]" {
		isGamePaused = !isGamePaused
	}
}

func IsWithinBoundaries(player *GameObject, boundary Boundary) bool {
	switch boundary {
	case Top:
		return player.row > 0
	case Bottom:
		_, height := screen.Size()
		return player.row < height - player.height
	default:
		fmt.Println("Failure checking boundaries")
		return false
	}
}

func InitUserInput() chan string {
	// creating a channel to be a communication channel between processes
	inputChan := make(chan string)
	go func() {
		for {
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventKey:
				inputChan <- ev.Name()
			}
		}
	}()

	return inputChan
}

func InitGameState() {
	width, height := screen.Size()

	paddleStart := height/2 - PaddleHeight/2

	player1Paddle = &GameObject{
		row: paddleStart, col: 0, width: 1, height: PaddleHeight,
		velRow: 0, velCol: 0,
		symbol: PaddleSymbol,
	}

	player2Paddle = &GameObject{
		row: paddleStart, col: width-1, width: 1, height: PaddleHeight,
		velRow: 0, velCol: 0,
		symbol: PaddleSymbol,
	}

	ball = &GameObject{
		row: height / 2, col: width / 2, width: 1, height: 1,
		velRow: InitialBallVelocityRow, velCol: InitialBallVelocityCol,
		symbol: BallSymbol,
	}

	gameObjects = []*GameObject{
		player1Paddle, player2Paddle, ball,
	}

	Print(paddleStart, 0, 1, PaddleHeight, PaddleSymbol)
	Print(paddleStart, width-1, 1, PaddleHeight, PaddleSymbol)
}

func ReadInput(inputChan chan string) string {
	var key string
	select {
		case key = <-inputChan:
		default:
			key = ""
	}
	return key
}

func PrintString(row, col int, str string) {
	for _, c := range str {
		
		screen.SetContent(col, row, c, nil, tcell.StyleDefault)
		col += 1
	}
}

func Print(row, col int, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
		
	}
}

func DrawState() {
	if isGamePaused {
		return
	}

	screen.Clear()
	// PrintString(0, 0, debugLog) //uncomment to enable debuglog
	for _, obj := range gameObjects {
		Print(obj.row, obj.col, obj.width, obj.height, obj.symbol)
	}
	// Print(player1Paddle.row, player1Paddle.col, player1Paddle.width, player1Paddle.height, PaddleSymbol)
	// Print(player2Paddle.row, player2Paddle.col, player2Paddle.width, player2Paddle.height, PaddleSymbol)
	// Print(ball.row, ball.col, ball.width, ball.height, BallSymbol)
	screen.Show()
}