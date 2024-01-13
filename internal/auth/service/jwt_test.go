package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateToken(t *testing.T) {

	authService := NewJWTService()

	result, err := authService.GenerateToken("meze@gmail.com", "1")
	assert.NotNil(t, result)
	assert.NoError(t, err)
}

func TestParseToken(t *testing.T) {
	authService := NewJWTService()

	email := "email@email.com"
	userID := "1"

	token, err := authService.GenerateToken(email, userID)
	parsedToken := authService.ParseToken(*token)

	assert.Equal(t, email, parsedToken.Claims.(*CustomClaims).Email)
	assert.NoError(t, err)
}
