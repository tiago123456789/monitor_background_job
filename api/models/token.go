package models

import (
	"errors"
	"os"

	"github.com/dgrijalva/jwt-go"
)

type TokenInterface interface {
	Get() (string, error)
	IsValid(tokenString string) (string, error)
}

type Token struct {
}

func NewToken() *Token {
	return &Token{}
}

func (t *Token) IsValid(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Token invalid")
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if err != nil {
		return "", err
	}
	return token.Raw, nil
}

func (t *Token) Get(payload jwt.Claims) (string, error) {
	atClaims := payload
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
