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
	id := query.Get("id")
	fmt.Printf("code: %s", code)

	client_id := "1009547623637712977"
	client_secret := "K-rktHBJ2haT1pqzZ8hs239M9n0PliFY"
	if code == "" {
		fmt.Println("Discord: code is empty")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Discord: code is empty"))
		return
	}

	data := url.Values{
		"client_id":     {client_id},
		"client_secret": {client_secret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {"http://localhost:8080/api/discord?id=" + id},
		"scope":         {"identify"},
	}

	postBody := bytes.NewBufferString(data.Encode())

	resp, err := http.Post("https://discord.com/api/v10/oauth2/token", "application/x-www-form-urlencoded", postBody)
	if err != nil {
		fmt.Println("Discord: Error while getting token", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Discord: Error while getting token" + err.Error()))
	}
	defer resp.Body.Close()

	bearer := &models.DiscordBearer{}
	err = json.NewDecoder(resp.Body).Decode(&bearer)
	if err != nil {
		fmt.Println("Discord: bearer decode error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Discord: bearer decode error" + err.Error()))
		return
	}

	var bearerString = "Bearer " + bearer.AccessToken

	req, err := http.NewRequest("GET", "https://discord.com/api/v10/users/@me", nil)
	if err != nil {
		fmt.Println("Discord: Error on request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Discord: Error on request: " + err.Error()))
		return
	}

	req.Header.Set("Authorization", bearerString)

	client := &http.Client{}
	resps, err := client.Do(req)
	if err != nil {
		fmt.Println("Discord: Error on response.\n[ERRO] -", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Discord: Error on response.\n[ERRO] -" + err.Error()))
		return
	}
	defer resps.Body.Close()

	discord := &models.Discord{}
	err = json.NewDecoder(resps.Body).Decode(&discord)
	if err != nil {
		fmt.Println("discord: parse: ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("discord: parse: " + err.Error()))
		return
	}

	db.Create(&discord)
	db.Model(&models.Socials{}).Where("ID = ?", id).Update("DiscordID", discord.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(discord)

}
