package helper

import (
	"dompet-miniprojectalta/config"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type JWTCustomClaims struct {
	UserID uint   `json:"userId"`
	Name  string `json:"name"`
	jwt.StandardClaims
}

func CreateToken(userId uint, name string) (string, error) {
	config.InitConfig()
	claims := &JWTCustomClaims{
		userId,
		name,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
		},
	}

	// Generate encoded token and send it as response.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.Cfg.TOKEN_SECRET))
	if err != nil {
		return "", err
	}
	return t, nil
}

func GetJwt(c echo.Context) (userId uint, name string) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTCustomClaims)
	userId = claims.UserID
	name = claims.Name
	return 
}