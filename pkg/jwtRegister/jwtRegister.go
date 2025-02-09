package jwtRegister

import (
	"diplom_chat_ssh/internal/model"
	custom_errors "diplom_chat_ssh/pkg/custom-errors"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"os"
)

func GenerateToken(claims *model.JWTCustomClaims) (string, *custom_errors.ErrHttp) {
	data, err := os.ReadFile("config/config.json")
	if err != nil {
		log.Error().Err(errors.New("file open: " + err.Error()))
	}

	var config model.Config

	if err = json.Unmarshal(data, &config); err != nil {
		log.Error().Err(errors.New("unmarshal config: " + err.Error()))
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte(config.JWT.Key))
	if err != nil {
		return "", custom_errors.New(fasthttp.StatusInternalServerError, "new with claims: "+err.Error())
	}

	return token, nil
}
