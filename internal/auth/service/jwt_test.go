package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateToken(t *testing.T) {

	authService := New()

	result := authService.GenerateToken("meze@gmail.com", "1")
	assert.NotNil(t, result)
}

func TestParseToken(t *testing.T) {
	authService := New()

	email := "email@email.com"
	userID := "1"

	token := authService.GenerateToken(email, userID)
	parsedToken := authService.ParseToken(token)

	assert.Equal(t, email, parsedToken.Claims.(*CustomClaims).Email)
}
