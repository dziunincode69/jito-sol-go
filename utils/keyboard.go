package utils

import (
	"fmt"

	"github.com/eiannone/keyboard"
)

func Keyboard(keychan chan int) {
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	var isSpressed bool
	var isKeyPressed bool
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}

		if event.Key == keyboard.KeyEsc {
			break
		}
		if event.Key == keyboard.KeyCtrlC {
			break
		}
		if event.Rune == 115 {
			if !isSpressed {
				isSpressed = true
			} else {
				isSpressed = false
			}
		}
		if isSpressed {
			switch event.Rune {
			case 48:
				fmt.Println("0 pressed, Selling 100% your token")
				keychan <- 100
				isKeyPressed = true
			case 49:
				fmt.Println("1 pressed, Selling 10% your token")
				keychan <- 10
				isKeyPressed = true

			case 50:
				fmt.Println("2 pressed, Selling 20% your token")
				keychan <- 20
				isKeyPressed = true

			case 51:
				fmt.Println("3 pressed, Selling 30% your token")
				isKeyPressed = true
				keychan <- 30
			case 52:
				fmt.Println("4 pressed, Selling 40% your token")
				isKeyPressed = true
				keychan <- 40
			case 53:
				fmt.Println("5 pressed, Selling 50% your token")
				isKeyPressed = true
				keychan <- 50
				// break
			case 54:
				fmt.Println("6 pressed, Selling 60% your token")
				keychan <- 60
				isKeyPressed = true
				// break
			case 55:
				fmt.Println("7 pressed, Selling 70% your token")
				isKeyPressed = true
				keychan <- 70

			case 56:
				fmt.Println("8 pressed, Selling 80% your token")
				isKeyPressed = true
				keychan <- 80
			case 57:
				fmt.Println("9 pressed, Selling 90% your token")
				isKeyPressed = true
				keychan <- 90
			}
		}
		if isKeyPressed {
			break
		}
	}
}
