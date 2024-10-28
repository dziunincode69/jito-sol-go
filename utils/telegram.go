package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendingTelegramNotification(chatId string, text string) {
	body, _ := json.Marshal(map[string]string{
		"chat_id": chatId,
		"text":    text,
	})
	url := "https://api.telegram.org/bot5700775936:AAHCXFUST_vaXjzdXEXxGLwzDqWx-IV0bSs/sendMessage"
	resp, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		fmt.Println("error sending telegram notification")
	}
	defer resp.Body.Close()
}
