package util

import (
	"errors"
	"time"

	"gubuk-service/config"

	"github.com/dgrijalva/jwt-go"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type UserPayload struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	UserID    string    `json:"user_id"`
	UserRole  string    `json:"user_role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// Valid checks if the token payload is valid or not
func (payload *UserPayload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}

	return nil
}

// CreateToken creates a new token with payload
func CreateToken(payload *UserPayload) (string, *UserPayload, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(config.SecretKey))

	return token, payload, err
}

// VerifyToken checks if the token is valid or not
func VerifyToken(token string) (*UserPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(config.SecretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &UserPayload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}

		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*UserPayload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
