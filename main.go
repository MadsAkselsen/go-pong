package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const PaddleHeight = 4
const PaddleSymbol = 0x2588

type Paddle struct {
	row, col, width, height int
}

var screen tcell.Screen
var player1 *Paddle
var player2 *Paddle
var debugLog string

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
		DrawState()
		time.Sleep(50 * time.Millisecond)

		key := ReadInput(inputChan)
		HandleUserInput(key)
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
	} else if key == "Rune[w]" && IsWithinBoundaries(player1, Top) {
		player1.row--
	} else if key == "Rune[s]" && IsWithinBoundaries(player1, Bottom) {
		player1.row++
	} else if key == "Up" && IsWithinBoundaries(player2, Top) {
		player2.row--
	} else if key == "Down" && IsWithinBoundaries(player2, Bottom) {
		player2.row++
	}
}

func IsWithinBoundaries(player *Paddle, boundary Boundary) bool {
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
				debugLog = ev.Name()
				inputChan <- ev.Name()
			}
		}
	}()

	return inputChan
}

func InitGameState() {
	width, height := screen.Size()
	paddleStart := height/2 - PaddleHeight/2

	player1 = &Paddle{
		row: paddleStart, col: 0, width: 1, height: PaddleHeight,
	}

	player2 = &Paddle{
		row: paddleStart, col: width-1, width: 1, height: PaddleHeight,
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
	screen.Clear()
	PrintString(0, 0, debugLog)
	Print(player1.row, player1.col, player1.width, player1.height, PaddleSymbol)
	Print(player2.row, player2.col, player2.width, player2.height, PaddleSymbol)
	screen.Show()
}