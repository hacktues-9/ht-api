package oauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/hacktues-9/API/pkg/models"
	"gorm.io/gorm"
)

var (
	discord_client_secret = os.Getenv("DISCORD_CLIENT_SECRET")
	github_client_secret  = os.Getenv("GITHUB_CLIENT_SECRET")
)

func GetDiscordInfo(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	hostUrl := os.Getenv("HOST_URL")
	query := r.URL.Query()
	code := query.Get("code")
	id := query.Get("id")

	client_id := "1009547623637712977"
	client_secret := discord_client_secret
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
		"redirect_uri":  {hostUrl + "api/user/discord?id=" + id},
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

	http.Redirect(w, r, "https://discord.gg/q6GGxvjjGb", http.StatusMovedPermanently)
}

func GetGithubInfo(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	query := r.URL.Query()
	code := query.Get("code")
	id := query.Get("state")

	client_id := "4f5f1918bf58eb0cccd4"
	client_secret := github_client_secret
	if code == "" {
		fmt.Println("Github: code is empty")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Github: code is empty"))
		return
	}

	data := url.Values{
		"client_id":     {client_id},
		"client_secret": {client_secret},
		"code":          {code},
	}

	postBody := bytes.NewBufferString(data.Encode())

	resp, err := http.Post("https://github.com/login/oauth/access_token", "application/x-www-form-urlencoded", postBody)
	if err != nil {
		fmt.Println("Github: Error while getting token", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Github: Error while getting token" + err.Error()))
	}
	defer resp.Body.Close()

	respString, _ := io.ReadAll(resp.Body)

	access_token := strings.Split(strings.Split(string(respString), "=")[1], "&")[0]

	bearerString := "Bearer " + access_token

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		fmt.Println("Github: Error on request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Github: Error on request: " + err.Error()))
		return
	}

	req.Header.Set("Authorization", bearerString)

	client := &http.Client{}
	resps, err := client.Do(req)
	if err != nil {
		fmt.Println("Github: Error on response.\n[ERRO] -", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Github: Error on response.\n[ERRO] -" + err.Error()))
		return
	}
	defer resps.Body.Close()

	github := &models.Github{}
	err = json.NewDecoder(resps.Body).Decode(&github)
	if err != nil {
		fmt.Println("github: parse: ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("github: parse: " + err.Error()))
		return
	}

	db.Create(&github)
	db.Model(&models.Socials{}).Where("ID = ?", id).Update("GithubID", github.ID)

	http.Redirect(w, r, "https://hacktues.bg/", http.StatusMovedPermanently)

}

// func GetGithubRepoInfo(w http.ResponseWriter, r *http.Request, db *gorm.DB) { // url...&scope=user%20repo
// 	query := r.URL.Query()
// 	code := query.Get("code")
// 	id := query.Get("state")

// 	client_id := "4f5f1918bf58eb0cccd4"
// 	client_secret := github_client_secret
// 	if code == "" {
// 		fmt.Println("Github: code is empty")
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte("Github: code is empty"))
// 		return
// 	}

// 	data := url.Values{
// 		"client_id":     {client_id},
// 		"client_secret": {client_secret},
// 		"code":          {code},
// 	}

// 	postBody := bytes.NewBufferString(data.Encode())

// 	resp, err := http.Post("https://github.com/login/oauth/access_token", "application/x-www-form-urlencoded", postBody)
// 	if err != nil {
// 		fmt.Println("Github: Error while getting token", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("Github: Error while getting token" + err.Error()))
// 	}
// 	defer resp.Body.Close()

// 	respString, _ := io.ReadAll(resp.Body)

// 	access_token := strings.Split(strings.Split(string(respString), "=")[1], "&")[0]

// 	bearerString := "Bearer " + access_token

// 	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
// 	if err != nil {
// 		fmt.Println("Github: Error on request:", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("Github: Error on request: " + err.Error()))
// 		return
// 	}

// 	req.Header.Set("Authorization", bearerString)

// 	client := &http.Client{}
// 	resps, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Github: Error on response.\n[ERRO] -", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte("Github: Error on response.\n[ERRO] -" + err.Error()))
// 		return
// 	}
// 	defer resps.Body.Close()

// 	github := &models.Github{}
// 	err = json.NewDecoder(resps.Body).Decode(&github)
// 	if err != nil {
// 		fmt.Println("github: parse: ", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte("github: parse: " + err.Error()))
// 		return
// 	}

// 	db.Create(&github)
// 	db.Model(&models.Socials{}).Where("ID = ?", id).Update("GithubID", github.ID)

// 	http.Redirect(w, r, "https://hacktues.bg/", http.StatusMovedPermanently)
// }
