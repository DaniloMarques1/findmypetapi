package test

import (
	"log"
	"testing"

	"github.com/danilomarques1/findmypetapi/util"
)

func TestNewToken(t *testing.T) {
	userId := "123"
	token, refreshToken, err := util.NewToken(userId)

	log.Printf("TOKEN RETURNED = %v\n", token)

	assertNil(t, err)
	assertNotEqual(t, "", token)
	assertNotEqual(t, "", refreshToken)
}

func TestNewApiError(t *testing.T) {
	msg := "This is a message"
	err := util.NewApiError(msg, 200)
	assertEqual(t, msg, err.Error())
}

func TestVerifyToken(t *testing.T) {
	userId := "123"
	token, refreshToken, err := util.NewToken(userId)
	assertNil(t, err)
	assertNotEqual(t, "", token)
	assertNotEqual(t, "", refreshToken)

	userClaims, err := util.VerifyToken(token)
	assertNil(t, err)
	assertEqual(t, util.TOKEN, userClaims.TokenRole)

	userClaims, err = util.VerifyToken(refreshToken)
	assertNil(t, err)
	assertEqual(t, util.REFRESH_TOKEN, userClaims.TokenRole)
}
