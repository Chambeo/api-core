package service

import (
	"chambeo-api-core/internal/auth/models"
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
	parsedToken, _ := authService.ParseToken(*token)

	assert.Equal(t, email, parsedToken.Claims.(*models.CustomClaims).Email)
	assert.NoError(t, err)
}

func TestParseTokenWithInvalidToken(t *testing.T) {
	authService := NewJWTService()

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSIsImVtYWlsIjoibWV6ZUBnbWFpbC5jb20iLCJpc3MiOiJjaG" +
		"FtYmVvLWNvIiwic3ViIjoiY2hhbWJlby1iZSIsImF1ZCI6WyJjaGFtYmVvLWZlIl0sImV4cCI6MTcwNTI3NjMyMiwibmJmIjoxNzA1MTg5OTI" +
		"yLCJpYXQiOjE3MDUxODk5MjIsImp0aSI6IjEifQ.p2jndX8Bn8q3mrJp4vv9nsGugZOZRcukrOBuMSIO4SAXX"
	parsedToken, err := authService.ParseToken(token)

	assert.Error(t, err)
	assert.Nil(t, parsedToken)
}
