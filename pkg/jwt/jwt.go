package jwt

import (
	b64 "encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hacktues-9/API/pkg/models"
	"gorm.io/gorm"
)

func CreateToken(ttl time.Duration, payload interface{}, privateKey string, publicKey string) (string, error) {
	privateKeyData, err := b64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("createToken: decode: private key: %w", err)
	}
	publicKeyData, err := b64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", fmt.Errorf("createToken: decode: public key: %w", err)
	}

	parsePrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return "", fmt.Errorf("createToken: parse: private key: %w", err)
	}

	parsePublicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return "", fmt.Errorf("createToken: parse: public key: %w", err)
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["sub"] = payload
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	token.Raw, err = token.SignedString(parsePrivateKey)
	if err != nil {
		return "", fmt.Errorf("createToken: signing string: %w", err)
	}

	_, err = jwt.ParseWithClaims(token.Raw, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return parsePublicKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("createToken: sign token: %w", err)
	}

	return token.Raw, nil
}

func ValidateToken(token string, publicKey string) (interface{}, error) {
	publicKeyData, err := b64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("validateToken: decode: public key: %w", err)
	}

	parsedPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, fmt.Errorf("validateToken: parse: public key: %w", err)
	}

	parsedToken, err := jwt.ParseWithClaims(token, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return parsedPublicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validateToken: parse token: %w", err)
	}

	claims, ok := parsedToken.Claims.(*jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validateToken: claims: %w", err)
	}

	return (*claims)["sub"], nil
}

func RefreshAccessToken(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	refreshTokenPublicKey := os.Getenv("REFRESH_TOKEN_PUBLIC_KEY")
	accessTokenPrivateKey := os.Getenv("ACCESS_TOKEN_PRIVATE_KEY")
	accessTokenPublicKey := os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")

	refreshToken := ""
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		fmt.Println("refreshAccessToken: get cookie: ", err)
		http.Error(w, "refreshAccessToken: get cookie: "+err.Error(), http.StatusBadRequest)
		return
	}

	accessTokenTTL := time.Hour * 24
	refreshToken = cookie.Value
	sub, err := ValidateToken(refreshToken, refreshTokenPublicKey)
	if err != nil {
		fmt.Println("refreshAccessToken: validate token: ", err)
		http.Error(w, "refreshAccessToken: validate token: "+err.Error(), http.StatusBadRequest)
		return
	}

	user := models.User{}
	db.Where("ID = ?", sub).First(&user)

	accessToken, err := CreateToken(accessTokenTTL, user.ID, accessTokenPrivateKey, accessTokenPublicKey)
	if err != nil {
		fmt.Println("refreshAccessToken: create token: ", err)
		http.Error(w, "refreshAccessToken: create token: "+err.Error(), http.StatusBadRequest)
		return
	}

	access_cookie := http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(accessTokenTTL),
		HttpOnly: true,
		Path:     "/",
	}

	http.SetCookie(w, &access_cookie)

	w.WriteHeader(http.StatusOK)
}
