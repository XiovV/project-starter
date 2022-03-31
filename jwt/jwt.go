package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

func New(userId int, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userId,
		"username": username,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func Validate(tok string) (*jwt.Token, error) {
	token, err := jwt.Parse(tok, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("couldn't verify signing method")
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func GetClaimString(token *jwt.Token, claim string) string {
	return token.Claims.(jwt.MapClaims)[claim].(string)
}

func GetClaimInt(token *jwt.Token, claim string) int {
	return int(token.Claims.(jwt.MapClaims)[claim].(float64))
}
