package features

import (
	"github.com/golang-jwt/jwt"
)

func CreateToken(id string) (string, error) {
	key := "private-key"

	var (
		t *jwt.Token
		s string
	)
	t = jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"ID": id,
	})
	var err error
	s, err = t.SignedString(key)
	return s, err
}
