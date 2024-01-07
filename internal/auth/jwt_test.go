package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	result := GenerateToken("meze@gmail.com", "1")
	assert.NotNil(t, result)
}

func TestParseToken(t *testing.T) {
	email := "email@email.com"
	userID := "1"
	token := GenerateToken(email, userID)
	parsedToken := ParseToken(token)

	assert.Equal(t, email, parsedToken.Claims.(*CustomClaims).Email)
}
