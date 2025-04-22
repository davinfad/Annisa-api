package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserAuthService interface {
	GenerateToken(username string) (string, error)
	ValidasiToken(token string) (*jwt.Token, error)
}

var SecretKey []byte

type jwtService struct {
}

func NewUserAuthService() *jwtService {
	return &jwtService{}
}

func (s *jwtService) SetSecretKey(key string) {
	SecretKey = []byte(key)
}

func (s *jwtService) GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(10 * time.Hour).Unix()

	claims := jwt.MapClaims{
		"username": username,
		"exp":      expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (s *jwtService) ValidasiToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(SecretKey), nil
	})
	if err != nil {
		return token, err
	}

	return token, nil
}
