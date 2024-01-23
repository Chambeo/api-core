package service

import (
	"chambeo-api-core/internal/auth/models"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

type AuthService struct{}

func NewJWTService() AuthService {
	return AuthService{}
}

func (a *AuthService) GenerateToken(email string, userId string) (*string, error) {

	mySigningKey := []byte("secretPassword")

	claims := a.generateClaims(email, userId)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Println("error trying to generate token")
		return nil, errors.New("error al intentar generar el token")
	}

	return &ss, err

}

func (a *AuthService) ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secretPassword"), nil
	})
	if err != nil {
		log.Println("ocurrio un error al intentar parsear el token")
		return nil, err
	}
	if claims, ok := token.Claims.(*models.CustomClaims); ok {
		log.Println(fmt.Sprintf("Retrieved userID claim is %s", claims.UserID))
		return token, nil
	} else {
		log.Println("ocurrio un error al intentar parsear los claims del token")
		return nil, errors.New("unknown error occurred trying to parse token claims")
	}

}

func (a *AuthService) generateClaims(email, userId string) models.CustomClaims {
	return models.CustomClaims{
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
