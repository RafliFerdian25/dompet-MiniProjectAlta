package helper

import (
	"dompet-miniprojectalta/config"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(userId uint, name string) (string, error) {
	// Set some claims
	claims := jwt.MapClaims{}
	claims["userId"] = userId
	claims["name"] = name
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	// Generate encoded token and send it as response.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.ConfigValue("TOKEN_SECRET")))
	if err != nil {
		return "", err
	}
	return t, nil
}