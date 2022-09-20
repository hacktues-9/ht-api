package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hacktues-9/API/pkg/models"
	"gorm.io/gorm"
)

func GetDiscordInfo(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	query := r.URL.Query()
	code := query.Get("code")
	fmt.Printf("code: %s", code)

	client_id := "1009547623637712977"
	client_secret := "K-rktHBJ2haT1pqzZ8hs239M9n0PliFY"
	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := url.Values{
		"client_id":     {client_id},
		"client_secret": {client_secret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {"http://192.168.1.57:8080/api/discord"},
		"scope":         {"identify"},
	}

	postBody := bytes.NewBufferString(data.Encode())

	resp, err := http.Post("https://discord.com/api/v10/oauth2/token", "application/x-www-form-urlencoded", postBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	bearer := &models.DiscordBearer{}
	err = json.NewDecoder(resp.Body).Decode(&bearer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var bearerString = "Bearer " + bearer.AccessToken

	req, err := http.NewRequest("GET", "https://discord.com/api/v10/users/@me", nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req.Header.Set("Authorization", bearerString)

	client := &http.Client{}
	resps, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer resps.Body.Close()

	user := &models.Discord{}
	err = json.NewDecoder(resps.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db.Create(&user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

}
