package features

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(id string, role int) (string, error) {
	key := []byte("private-key")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":   id,
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(key)
}

func ParseToken(tokenString string) (string, string, error) {
	key := []byte("private-key")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return key, nil
	})

	if err != nil {
		return "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, ok := claims["ID"].(string)
		if !ok {
			return "", "", errors.New("ID claim missing or invalid")
		}
		role, ok := claims["role"].(string)
		if !ok {
			return "", "", errors.New("ROLE is missing")
		}

		return id, role, nil
	}
	return "", "", errors.New("invalid token")
}
