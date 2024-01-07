package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type CustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(email string, userId string) *jwt.Token {

	mySigningKey := []byte("secretPassword")

	claims := generateClaims(email, userId)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	fmt.Println(ss, err)

	return token

}

func generateClaims(email, userId string) CustomClaims {
	return CustomClaims{
		UserID: userId,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "chambeo-co",
			Subject:   "chambeo-be",
			ID:        userId,
			Audience:  []string{"chambeo-fe"},
		},
	}
}

func GetSignedJWT(jwt *jwt.Token) string {
	signedToken, _ := jwt.SigningString()
	return signedToken
}
