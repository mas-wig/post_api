package utils

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(ttl time.Duration, payload interface{}, privateKey string) (string, error) {
	decodePrivateKey, err := base64.StdEncoding.DecodeString(privateKey) // decode key
	if err != nil {
		return "", fmt.Errorf("tidak bisa melakukan decode key : %w ", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodePrivateKey) // paser key rsa
	if err != nil {
		return "", fmt.Errorf("tidak dapat melakukan parse key : %w", err)
	}

	now := time.Now().UTC()
	claims := make(jwt.MapClaims) // create claims token

	claims["sub"] = payload
	claims["exp"] = now.Add(ttl).Unix() // expired jwt cliams
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key) // sign ke rs256
	if err != nil {
		return "", fmt.Errorf("tidak dapat melakukan sign token : %w", err)
	}
	return token, nil
}

func ValidateToken(token string, publicKey string) (interface{}, error) {
	decodePublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("tidak dapat melakukan decode public key : %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodePublicKey)
	if err != nil {
		return nil, fmt.Errorf("validate : parse key => %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, err := t.Method.(*jwt.SigningMethodRSA); !err {
			return nil, fmt.Errorf("unexpected method => %s ", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate => %w ", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate => invalid token ")
	}

	return claims["sub"], nil
}
