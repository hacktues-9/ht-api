package jwt

import (
	"net/http"
	"os"
	"time"
)

var (
	accessTokenTTL         = time.Hour * 24
	refreshTokenTTL        = time.Hour * 24 * 7
	accessTokenPrivateKey  = os.Getenv("ACCESS_TOKEN_PRIVATE_KEY")
	accessTokenPublicKey   = os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")
	refreshTokenPrivateKey = os.Getenv("REFRESH_TOKEN_PRIVATE_KEY")
	refreshTokenPublicKey  = os.Getenv("REFRESH_TOKEN_PUBLIC_KEY")
)

func GenerateCookies(userID uint) (http.Cookie, http.Cookie, error) {
	accessToken, err := CreateToken(accessTokenTTL, userID, accessTokenPrivateKey, accessTokenPublicKey)
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}

	refreshToken, err := CreateToken(refreshTokenTTL, userID, refreshTokenPrivateKey, refreshTokenPublicKey)
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}

	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(refreshTokenTTL),
		HttpOnly: true,
		Domain:   os.Getenv("HOST_DOMAIN"),
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	accessCookie := http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(accessTokenTTL),
		HttpOnly: true,
		Domain:   os.Getenv("HOST_DOMAIN"),
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	return accessCookie, refreshCookie, nil
}
