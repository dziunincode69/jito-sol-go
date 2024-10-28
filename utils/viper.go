package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
)

var privateKey string

func InitializeVipers() {
	viper.SetConfigFile("config.yml")
	if err := viper.ReadInConfig(); err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			fmt.Println("Config file doesn't exist")

			fmt.Print("Private Key: ")
			_, err = fmt.Scanln(&privateKey)
			if err != nil {
				return
			}

			privateKeyEncrypt := EncryptString(privateKey)
			WriteConfig(privateKeyEncrypt)

			c := exec.Command("clear")
			fmt.Println("Config file created")
			c.Stdout = os.Stdout
			err = c.Run()
			if err != nil {
				return
			}

			os.Exit(1)
		}
	}
}

func DecryptString(encryptedString string) (decryptedString string) {
	bytes := []byte(viper.GetString("EncryptionKey"))
	keyString := hex.EncodeToString(bytes)
	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalln(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalln(err.Error())
	}

	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	var plaintext []byte
	plaintext, err = aesGCM.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		log.Fatalln(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}

func WriteConfig(privateKey string) {
	file, err := os.Create("config.yml")
	if err != nil {
		log.Fatalln(err.Error())
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}(file)
	configString := `PrivateKey: "` + privateKey + `"
Fee: 0.01
Https: "https://mainnet-beta.solana.com"
`
	_, err = file.WriteString(configString)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func EncryptString(stringToEncrypt string) (encryptedString string) {
	bytes := []byte(viper.GetString("EncryptionKey"))

	keyString := hex.EncodeToString(bytes)
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalln(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalln(err.Error())
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalln(err.Error())
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}
