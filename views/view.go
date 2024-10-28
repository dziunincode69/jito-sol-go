package views

import (
	"fmt"

	"github.com/fatih/color"
)

func SellMonitorViews(token string) {
	// fmt.Print("\033[H\033[2J")
	color.HiGreen("Monitor Price + Keyboard listener Started")
	color.Green("Hotkeys: ")
	color.Cyan("S1-9 sell 10-90% token ")
	color.Cyan("S0 sell 100% token ")
	color.Magenta("ex: S5 for sell 50% token ")
	fmt.Println("\n\n")
	color.HiGreen("Token: " + token)
}
