package password

import (
	"crypto/rand"
	"errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"math/big"
)

const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func CheckPass(passItem, passReq string, passChan chan bool) {
	err := bcrypt.CompareHashAndPassword([]byte(passItem), []byte(passReq))
	if err == nil {
		passChan <- true
	} else {
		passChan <- false
	}
}

func GenPass(passReq string, passChan chan string) {
	pass, err := bcrypt.GenerateFromPassword([]byte(passReq), 0)
	if err != nil {
		log.Error().Err(errors.New("generate from password: " + err.Error())).Msg("")
	}

	passChan <- string(pass)
}

func RandomPassword() string {
	passwordBytes := make([]byte, 8)

	for i := range passwordBytes {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err) // никогда не должно случиться
		}
		passwordBytes[i] = charset[num.Int64()]
	}

	return string(passwordBytes)
}
