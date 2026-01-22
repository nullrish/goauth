// Package auth is used to sign and unsign jwt
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nullrish/goauth/internal/keys"
	"github.com/nullrish/goauth/model"
)

func SignJWT(u *model.User) (string, error) {
	claims := jwt.MapClaims{
		"id":       u.ID,
		"username": u.Username,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(keys.PrivateKey)
}

func VerifyJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return keys.PublicKey, nil
	})
}
