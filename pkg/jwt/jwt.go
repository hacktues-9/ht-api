package jwt

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(ttl time.Duration, payload interface{}, privateKey string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("could not decode key: %w", err)
	}

	jwtkey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := jwt.MapClaims{
		"sub" : payload,
		"exp" : now.Add(ttl).Unix(),
		"iat" : now.Unix(),
		"nbf" : now.Unix()
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(jwtkey)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func ValidateToken(token string, publicKey string) (interface{}, error) {
	key, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode key: %w", err)
	}

	jwtkey, err := jwt.ParseRSAPublicKeyFromPEM(key)
	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}

	token, err = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("invalid signing method: %s", t.Header["alg"])
		}
		return jwtkey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate: parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["sub"], nil
	}

	return nil, fmt.Errorf("invalid token")
}