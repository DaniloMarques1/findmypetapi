package test

import (
	"log"
	"testing"

	"github.com/danilomarques1/findmypetapi/util"
)

func TestNewToken(t *testing.T) {
	userId := "123"
	token, err := util.NewToken(userId)

	log.Printf("TOKEN RETURNED = %v\n", token)

	assertNil(t, err)
	assertNotEqual(t, "", token)
}
