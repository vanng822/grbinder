package grbinder

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("Test starting")
	InitDefaultLocker()
	retCode := m.Run()
	fmt.Println("Test ending")
	os.Exit(retCode)
}
