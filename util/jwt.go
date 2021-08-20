package util

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

type UserClaims struct {
	jwt.StandardClaims
	tokenRole string
	userId    string
}

const EXPIRATION_TIME = 86400
const TOKEN = "token"
const REFRESH_TOKEN = "refresh_token"

func NewToken(userId string) (string, error) {
	userClaims := UserClaims{
		tokenRole: TOKEN,
		userId:    userId,
	}

	fmt.Printf("key %T\n", os.Getenv("JWT_KEY"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
