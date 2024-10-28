package utils

import (
	"fmt"
	"time"
)

func LogWithTimestamp(message string) {
	currentTime := time.Now().Format("2006/01/02 15:04:05.000")
	fmt.Println(fmt.Sprintf("%s - %s", currentTime, message))
}
