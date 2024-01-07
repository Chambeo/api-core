package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	result := GenerateToken("meze@gmail.com", "1")
	issuer, _ := result.Claims.GetIssuer()
	assert.NotNil(t, result)
	assert.Equal(t, "chambeo-co", issuer)
}

func TestGetSignedJWT(t *testing.T) {
	token := GenerateToken("email@email.com", "1")

	getSignedToken := GetSignedJWT(token)

	assert.NotNil(t, getSignedToken)
}
