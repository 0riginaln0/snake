/*
	w4 bundle build/cart.wasm --title "Snake" --html release/snake.html --windows release/snake.exe --mac release/snake-mac --linux release/snake-linux
*/

package main

import (
	"cart/difficulty"
	"cart/state"
	"cart/w4"
	"math/rand"
	"strconv"
)

const winningScore = 100

type Direction uint8

const (
	Up Direction = iota
	Down
	Left
	Right
)

var (
	snake           = &Snake{}
	frameCounter    = 0
	timeoutCounter  = 0 // Used for menu animations
	prevState       uint8
	prevDir         = Right
	fruit           = Point{X: 10, Y: 10}
	fruitSprite     = [16]byte{0x00, 0xa0, 0x02, 0x00, 0x0e, 0xf0, 0x36, 0x5c, 0xd6, 0x57, 0xd5, 0x57, 0x35, 0x5c, 0x0f, 0xf0}
	rnd             func(int) int
	speed           = 15 // 12, 10, 6
	score           = 0
	difficultyLevel = difficulty.Easy
	gameState       = state.Start
	inputBuffer     = make([]Direction, 0, 60)
)

//go:export start
func start() {
	w4.PALETTE[0] = 0xfbf7f3
	w4.PALETTE[1] = 0xe5b083
	w4.PALETTE[2] = 0x426e5d
	w4.PALETTE[3] = 0x20283d

	snake.Reset()
}

//go:export update
func update() {
	frameCounter++

	switch gameState {
	case state.Win:
		winScreen()
	case state.Start:
		startScreen()
	case state.GameOver:
		gameOver()
	case state.Playing:
		playing()
	}
}

func winScreen() {
	w4.Text("Wow. I am impressed", 0, 70)
	w4.Text("Well done", 0, 90)
	w4.Text("Take a break :)", 0, 110)
}

func startScreen() {
	w4.Text("SNAKE", 30, 30)

	timeoutCounter += 1
	if timeoutCounter >= 90 {
		w4.Text("Made by me", 20, 90)
		w4.Text("github.com", 0, 110)
		w4.Text("/", 75, 120)
		w4.Text("0riginaln0", 80, 130)

		startGameOnInput()
	}
}

func gameOver() {
	w4.Text("GAME OVER", 30, 70)

	timeoutCounter += 1
	if timeoutCounter >= 60 {
		w4.Text("good luck next time", 5, 90)
	}
	if timeoutCounter >= 120 {
		w4.Text(":)", 75, 120)
	}
	if timeoutCounter >= 180 {
		startGameOnInput()
	}
}

func startGameOnInput() {
	justPressed := *w4.GAMEPAD1 & (*w4.GAMEPAD1 ^ prevState)
	anyButtonPressed := justPressed&w4.BUTTON_UP != 0 || justPressed&w4.BUTTON_DOWN != 0 ||
		justPressed&w4.BUTTON_LEFT != 0 || justPressed&w4.BUTTON_RIGHT != 0

	if anyButtonPressed && (gameState == state.Start || gameState == state.GameOver) {
		timeoutCounter = 0

		gameState = state.Playing
	}
	prevState = *w4.GAMEPAD1
}

func playing() {
	takeInput()

	w4.Text(
		"Score:"+
			strconv.FormatInt(int64(score), 10)+
			"/"+
			strconv.FormatInt(int64(winningScore), 10),
		4, 4)
	w4.Text(difficultyLevel.String(), 92, 148)

	if frameCounter%speed == 0 {
		if len(inputBuffer) != 0 {
			switch inputBuffer[0] {
			case Up:
				snake.Up()
			case Down:
				snake.Down()
			case Left:
				snake.Left()
			case Right:
				snake.Right()
			}
			inputBuffer = inputBuffer[1:]
		}
		snake.Update()

		if snake.IsDead() {
			difficultyLevel.Reset()
			snake.Reset()
			speed = 15
			score = 0
			frameCounter = 0
			gameState = state.GameOver
		}

		if snake.Body[0] == fruit {
			//w4.Trace("am nyam")
			score += 1
			snake.Body = append(snake.Body, snake.Body[len(snake.Body)-1])

			generateNewApple()

			switch score {
			case 20:
				speed = 12
				difficultyLevel = difficulty.StillEz
			case 40:
				speed = 10
				difficultyLevel = difficulty.Medium
			case 60:
				speed = 6
				difficultyLevel = difficulty.Hard
			case winningScore:
				gameState = state.Win
			}
		}
	}
	snake.Draw()

	*w4.DRAW_COLORS = 0x4320
	w4.Blit(&fruitSprite[0], fruit.X*8, fruit.Y*8, 8, 8, w4.BLIT_2BPP)
}

func takeInput() {
	justPressed := *w4.GAMEPAD1 & (*w4.GAMEPAD1 ^ prevState)

	if *w4.GAMEPAD1 != 0 && rnd == nil {
		rnd = rand.New(rand.NewSource(int64(frameCounter))).Intn
	}

	if justPressed&w4.BUTTON_UP != 0 && prevDir != Down && prevDir != Up {
		inputBuffer = append(inputBuffer, Up)
		prevDir = Up
	} else if justPressed&w4.BUTTON_DOWN != 0 && prevDir != Down && prevDir != Up {
		inputBuffer = append(inputBuffer, Down)
		prevDir = Down
	} else if justPressed&w4.BUTTON_LEFT != 0 && prevDir != Left && prevDir != Right {
		inputBuffer = append(inputBuffer, Left)
		prevDir = Left
	} else if justPressed&w4.BUTTON_RIGHT != 0 && prevDir != Left && prevDir != Right {
		inputBuffer = append(inputBuffer, Right)
		prevDir = Right
	}

	prevState = *w4.GAMEPAD1
}

func generateNewApple() {
	foundNewPlace := false
	p := Point{}
	for !foundNewPlace {
		foundNewPlace = true
		p.X = rnd(20)
		p.Y = rnd(20)
		for _, body_part := range snake.Body {
			if body_part == p {
				foundNewPlace = false
				break
			}
		}
	}
	fruit = p
}
