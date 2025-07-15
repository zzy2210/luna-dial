package biz_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestValidUUID(t *testing.T) {
	// 550e8400e29b41d4a716446655440000
	if isValidUUID("550e8400e29b41d4a716446655440000") {
		fmt.Println("Valid UUID")
	} else {
		fmt.Println("Invalid UUID")
	}
}

func isValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
