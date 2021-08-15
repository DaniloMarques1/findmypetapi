package test

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func executeRequest(request http.Request) {
	// TODO actually do the request
}

func assertEqual(t *testing.T, expect, actual interface{}) {
	if expect != actual {
		t.Fatalf(fmt.Sprintf("Expected value %v\nActual value: %v\n", expect, actual))
	}
}
