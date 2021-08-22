package util

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserClaims struct {
	jwt.StandardClaims
	TokenRole string
	UserId    string
}

const (
	EXPIRATION_TIME               = 86400
	REFRESH_TOKEN_EXPIRATION_TIME = EXPIRATION_TIME * 3
	TOKEN                         = "token"
	REFRESH_TOKEN                 = "refresh_token"
)

// return a token, refresh token and a possible error
func NewToken(userId string) (string, string, error) {
	tokenClaims := UserClaims{
		TokenRole: TOKEN,
		UserId:    userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: EXPIRATION_TIME + time.Now().Unix(),
		},
	}
	token, err := generateToken(tokenClaims)
	if err != nil {
		return "", "", err
	}

	refreshTokenClaims := UserClaims{
		TokenRole: REFRESH_TOKEN,
		UserId:    userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: REFRESH_TOKEN_EXPIRATION_TIME + time.Now().Unix(),
		},
	}
	refreshToken, err := generateToken(refreshTokenClaims)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func generateToken(userClaims UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func VerifyToken(tokenString string) (*UserClaims, error) {
	var userClaims UserClaims
	_, err := jwt.ParseWithClaims(tokenString, &userClaims, func(t *jwt.Token) (
		interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	return &userClaims, nil
}

func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		slice := strings.Split(bearer, " ")
		if len(slice) < 2 {
			HandleError(w, NewApiError("Invalid token", http.StatusUnauthorized))
			return
		} else {
			token := slice[1]
			userClaims, err := VerifyToken(token)
			if err != nil {
				HandleError(w, NewApiError("Invalid token", http.StatusUnauthorized))
				return
			}

			r.Header.Add("user_id", userClaims.UserId)
			next.ServeHTTP(w, r)
		}
	})
}
