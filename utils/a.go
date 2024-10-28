package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	json "github.com/json-iterator/go"
)

type APIRES struct {
	Status    int    `json:"status"`
	Name      string `json:"name"`
	ExpiredOn string `json:"expiredOn"`
	Message   string `json:"message"`
	Token     string `json:"token"`
}

func CheckJWTTOKEN(tokenstr string) {
	mikey := "9297519a9e99804dc27282fae5ef9edfa87907c9"
	secretKey := []byte(mikey)
	token, err := jwt.Parse(tokenstr, func(token *jwt.Token) (interface{}, error) {
		// Pastikan algoritma yang digunakan adalah HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		os.Exit(1)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		iat := int64(claims["iat"].(float64))
		exp := int64(claims["exp"].(float64))
		now := time.Now().Unix()
		if now < iat {
			fmt.Println("Token issued in the future!")
			os.Exit(0)
		} else if now > exp {
			os.Exit(0)
		}
	} else {
		os.Exit(0)
	}
}

func CheckWhitelistAddr(address string) {
	var apires APIRES
	resp, err := http.Get("https://liongfamily.net/check?address=" + address)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	json.Unmarshal(body, &apires)
	if apires.Status == 200 {
		fmt.Println(apires.Message)
		CheckJWTTOKEN(apires.Token)
	} else {
		os.Exit(0)
	}
}
