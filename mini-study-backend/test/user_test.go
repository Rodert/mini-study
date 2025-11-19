package test

import (
	"testing"

	"github.com/javapub/mini-study/mini-study-backend/internal/utils"
)

func TestHashPassword(t *testing.T) {
	pwd := "secret123"
	hash, err := utils.HashPassword(pwd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := utils.CheckPassword(hash, pwd); err != nil {
		t.Fatalf("password validation failed: %v", err)
	}
}
