/*
	w4 bundle build/cart.wasm --title "Snake" --html release/snake.html --windows release/snake.exe --mac release/snake-mac --linux release/snake-linux
*/

package main

import (
	"cart/w4"
	"math/rand"
	"strconv"
)

var (
	snake        = &Snake{}
	frameCount   = 0
	prevState    uint8
	fruit        = Point{X: 10, Y: 10}
	rnd          func(int) int
	fruitSprite  = [16]byte{0x00, 0xa0, 0x02, 0x00, 0x0e, 0xf0, 0x36, 0x5c, 0xd6, 0x57, 0xd5, 0x57, 0x35, 0x5c, 0x0f, 0xf0}
	speed        = 15 // 12, 10, 6
	score        = 0
	winningScore = 100
	level        = "easy"  // still ez, medium, hard
	mode         = "start" // playing, game over, win
	input_taken  = false
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
	frameCount++
	switch mode {
	case "start":
		startScreen()
	case "playing":
		playing()
	case "game over":
		gameOver()
	case "win":
		winScreen()
	}
}

func anyInput() {
	justPressed := *w4.GAMEPAD1 & (*w4.GAMEPAD1 ^ prevState)
	if (justPressed&w4.BUTTON_UP != 0 || justPressed&w4.BUTTON_DOWN != 0 ||
		justPressed&w4.BUTTON_LEFT != 0 || justPressed&w4.BUTTON_RIGHT != 0) &&
		(mode == "start" || mode == "game over") {
		mode = "playing"
	}

	prevState = *w4.GAMEPAD1
}

func input() {
	justPressed := *w4.GAMEPAD1 & (*w4.GAMEPAD1 ^ prevState)

	if *w4.GAMEPAD1 != 0 && rnd == nil {
		rnd = rand.New(rand.NewSource(int64(frameCount))).Intn
	}

	if justPressed&w4.BUTTON_UP != 0 {
		input_taken = true
		snake.Up()
	} else if justPressed&w4.BUTTON_DOWN != 0 {
		input_taken = true
		snake.Down()
	} else if justPressed&w4.BUTTON_LEFT != 0 {
		input_taken = true
		snake.Left()
	} else if justPressed&w4.BUTTON_RIGHT != 0 {
		input_taken = true
		snake.Right()
	}
	prevState = *w4.GAMEPAD1
}

func winScreen() {
	w4.Text("Wow. I am impressed", 0, 70)
	w4.Text("Well done", 0, 90)
	w4.Text("Take a break :)", 0, 110)
}

func gameOver() {
	w4.Text("GAME OVER", 30, 70)
	w4.Text("good luck next time", 5, 90)
	w4.Text(":)", 75, 120)
	anyInput()
}

func startScreen() {
	anyInput()
	w4.Text("SNAKE", 30, 30)
	w4.Text("Made by me", 20, 90)
	w4.Text("github.com", 0, 110)
	w4.Text("/", 75, 120)
	w4.Text("0riginaln0", 80, 130)
}

func playing() {
	w4.Text(
		"Score:"+strconv.FormatInt(int64(score), 10)+
			"/"+strconv.FormatInt(int64(winningScore), 10),
		4, 4)
	w4.Text(level, 92, 148)

	if !input_taken {
		input()
	}

	if frameCount%speed == 0 {
		input_taken = false
		snake.Update()

		if snake.IsDead() {
			level = "easy"
			snake.Reset()
			speed = 15
			score = 0
			frameCount = 0
			mode = "game over"
		}

		if snake.Body[0] == fruit {
			score += 1
			snake.Body = append(snake.Body, snake.Body[len(snake.Body)-1])

			generateNewApple()

			switch score {
			case 20:
				speed = 12
				level = "still ez"
			case 40:
				speed = 10
				level = "medium"
			case 60:
				speed = 6
				level = "hard"
			case winningScore:
				mode = "win"
			}
		}
	}
	snake.Draw()

	*w4.DRAW_COLORS = 0x4321
	w4.Blit(&fruitSprite[0], fruit.X*8, fruit.Y*8, 8, 8, w4.BLIT_2BPP)
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
