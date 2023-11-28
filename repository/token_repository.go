package repository

import (
	"errors"
	"fmt"
	"store/config"

	"github.com/golang-jwt/jwt"
)

type TokenRepository interface {
	CreateToken(makeClaims func(string) jwt.Claims) (string, error)
	VerifyToken(tokenString string) (jwt.MapClaims, error)
}

type tokenRepository struct {
	tokenConfig config.TokenConfig
}

func (t *tokenRepository) CreateToken(makeClaims func(string) jwt.Claims) (string, error) {
	claims := makeClaims(t.tokenConfig.ApplicationName)
	token := jwt.NewWithClaims(
		t.tokenConfig.JwtSigningMethod,
		claims,
	)

	newToken, err := token.SignedString([]byte(t.tokenConfig.JwtSignatureKey))
	if err != nil {
		return "", err
	}
	return newToken, nil
}

func (t *tokenRepository) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("signing method invalid")
		} else if method != t.tokenConfig.JwtSigningMethod {
			return nil, fmt.Errorf("signing method invalid")
		}
		return []byte(t.tokenConfig.JwtSignatureKey), nil
	})
	if err != nil {
		fmt.Println("Parsing failed..")
		return nil, errors.New("tokenExpired")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid || claims["iss"] != t.tokenConfig.ApplicationName {
		fmt.Println("Token invalid..")
		return nil, err
	}
	return claims, nil
}

func NewTokenRepository(tokenConfig config.TokenConfig) TokenRepository {
	repo := new(tokenRepository)
	repo.tokenConfig = tokenConfig
	return repo
}
