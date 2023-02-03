package oauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hacktues-9/API/pkg/jwt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/hacktues-9/API/pkg/models"
	"gorm.io/gorm"
)

var (
	discordClientSecret = os.Getenv("DISCORD_CLIENT_SECRET")
	githubClientSecret  = os.Getenv("GITHUB_CLIENT_SECRET")
)

func GetDiscordInfo(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	hostUrl := os.Getenv("HOST_URL")
	query := r.URL.Query()
	code := query.Get("code")
	id := query.Get("state")

	clientId := "1009547623637712977"
	clientSecret := discordClientSecret
	if code == "" {
		fmt.Printf("[ ERROR ][ GetDiscordInfo ] code is empty\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "code is empty", 0), errors.New("code is empty"), http.StatusBadRequest, "GetDiscordInfo")
		return
	}

	data := url.Values{
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {hostUrl + "api/user/discord"},
		"scope":         {"identify"},
	}

	postBody := bytes.NewBufferString(data.Encode())

	resp, err := http.Post("https://discord.com/api/v10/oauth2/token", "application/x-www-form-urlencoded", postBody)
	if err != nil {
		//fmt.Println("Discord: Error while getting token", err)
		//w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte("Discord: Error while getting token" + err.Error()))
		fmt.Printf("[ ERROR ][ GetDiscordInfo ] Error while getting token: %s\n", err.Error())
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "Discord: Error while getting token"+err.Error(), 0), err, http.StatusInternalServerError, "GetDiscordInfo")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			//fmt.Println("Discord: Error while closing body", err)
			//w.WriteHeader(http.StatusInternalServerError)
			//w.Write([]byte("Discord: Error while closing body" + err.Error()))
			fmt.Printf("[ ERROR ][ GetDiscordInfo ] Error while closing body: %s\n", err.Error())
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "Discord: Error while closing body"+err.Error(), 0), err, http.StatusInternalServerError, "GetDiscordInfo")
			return
		}
	}(resp.Body)

	bearer := &models.DiscordBearer{}
	err = json.NewDecoder(resp.Body).Decode(&bearer)
	if err != nil {
		//fmt.Println("Discord: bearer decode error", err)
		//w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte("Discord: bearer decode error" + err.Error()))
		fmt.Printf("[ ERROR ][ GetDiscordInfo ] bearer decode error: %s\n", err.Error())
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "Discord: bearer decode error"+err.Error(), 0), err, http.StatusInternalServerError, "GetDiscordInfo")
		return
	}

	var bearerString = "Bearer " + bearer.AccessToken

	req, err := http.NewRequest("GET", "https://discord.com/api/v10/users/@me", nil)
	if err != nil {
		//fmt.Println("Discord: Error on request:", err)
		//w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte("Discord: Error on request: " + err.Error()))
		fmt.Printf("[ ERROR ][ GetDiscordInfo ] Error on request: %s\n", err.Error())
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "Discord: Error on request: "+err.Error(), 0), err, http.StatusInternalServerError, "GetDiscordInfo")
		return
	}

	req.Header.Set("Authorization", bearerString)

	client := &http.Client{}
	resps, err := client.Do(req)
	if err != nil {
		//fmt.Println("Discord: Error on response.\n[ERRO] -", err)
		//w.WriteHeader(http.StatusBadRequest)
		//w.Write([]byte("Discord: Error on response.\n[ERRO] -" + err.Error()))
		fmt.Printf("[ ERROR ][ GetDiscordInfo ] Error on response: %s\n", err.Error())
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "Discord: Error on response.\n[ERRO] -"+err.Error(), 0), err, http.StatusBadRequest, "GetDiscordInfo")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			//fmt.Println("Discord: Error while closing body", err)
			//w.WriteHeader(http.StatusInternalServerError)
			//w.Write([]byte("Discord: Error while closing body" + err.Error()))
			fmt.Printf("[ ERROR ][ GetDiscordInfo ] Error while closing body: %s\n", err.Error())
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "Discord: Error while closing body"+err.Error(), 0), err, http.StatusInternalServerError, "GetDiscordInfo")
			return
		}
	}(resps.Body)

	discord := &models.Discord{}
	err = json.NewDecoder(resps.Body).Decode(&discord)
	if err != nil {
		//fmt.Println("discord: parse: ", err)
		//w.WriteHeader(http.StatusBadRequest)
		//w.Write([]byte("discord: parse: " + err.Error()))
		fmt.Printf("[ ERROR ][ GetDiscordInfo ] parse: %s\n", err.Error())
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "discord: parse: "+err.Error(), 0), err, http.StatusBadRequest, "GetDiscordInfo")
		return
	}

	db.Create(&discord)
	db.Model(&models.Users{}).Joins("info ON users.info_id = info.id").Joins("socials ON info.socials_id = socials.id").Where("users.id = ?", id).Update("socials.discord_id", discord.ID)

	http.Redirect(w, r, "https://discord.gg/q6GGxvjjGb", http.StatusMovedPermanently)
}

func GetGithubInfo(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	query := r.URL.Query()
	code := query.Get("code")
	id := query.Get("state")

	client_id := "4f5f1918bf58eb0cccd4"
	client_secret := githubClientSecret
	if code == "" {
		//fmt.Println("Github: code is empty")
		//w.WriteHeader(http.StatusBadRequest)
		//w.Write([]byte("Github: code is empty"))
		fmt.Printf("[ ERROR ][ GetGithubInfo ] code is empty\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "Github: code is empty", 0), nil, http.StatusBadRequest, "GetGithubInfo")
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
		//fmt.Println("Github: Error while getting token", err)
		//w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte("Github: Error while getting token" + err.Error()))
		fmt.Printf("[ ERROR ][ GetGithubInfo ] Error while getting token: %s\n", err.Error())
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "Github: Error while getting token"+err.Error(), 0), err, http.StatusInternalServerError, "GetGithubInfo")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			//fmt.Println("Github: Error while closing body", err)
			//w.WriteHeader(http.StatusInternalServerError)
			//w.Write([]byte("Github: Error while closing body" + err.Error()))
			fmt.Printf("[ ERROR ][ GetGithubInfo ] Error while closing body: %s\n", err.Error())
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "Github: Error while closing body"+err.Error(), 0), err, http.StatusInternalServerError, "GetGithubInfo")
			return
		}
	}(resp.Body)

	respString, _ := io.ReadAll(resp.Body)

	access_token := strings.Split(strings.Split(string(respString), "=")[1], "&")[0]

	bearerString := "Bearer " + access_token

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		//fmt.Println("Github: Error on request:", err)
		//w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte("Github: Error on request: " + err.Error()))
		fmt.Printf("[ ERROR ][ GetGithubInfo ] Error on request: %s\n", err.Error())
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "Github: Error on request: "+err.Error(), 0), err, http.StatusInternalServerError, "GetGithubInfo")
		return
	}

	req.Header.Set("Authorization", bearerString)

	client := &http.Client{}
	resps, err := client.Do(req)
	if err != nil {
		//fmt.Println("Github: Error on response.\n[ERRO] -", err)
		//w.WriteHeader(http.StatusBadRequest)
		//w.Write([]byte("Github: Error on response.\n[ERRO] -" + err.Error()))
		fmt.Printf("[ ERROR ][ GetGithubInfo ] Error on response: %s\n", err.Error())
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "Github: Error on response.\n[ERRO] -"+err.Error(), 0), err, http.StatusBadRequest, "GetGithubInfo")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			//fmt.Println("Github: Error while closing body", err)
			//w.WriteHeader(http.StatusInternalServerError)
			//w.Write([]byte("Github: Error while closing body" + err.Error()))
			fmt.Printf("[ ERROR ][ GetGithubInfo ] Error while closing body: %s\n", err.Error())
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "Github: Error while closing body"+err.Error(), 0), err, http.StatusInternalServerError, "GetGithubInfo")
			return
		}
	}(resps.Body)

	github := &models.Github{}
	err = json.NewDecoder(resps.Body).Decode(&github)
	if err != nil {
		//fmt.Println("github: parse: ", err)
		//w.WriteHeader(http.StatusBadRequest)
		//w.Write([]byte("github: parse: " + err.Error()))
		fmt.Printf("[ ERROR ][ GetGithubInfo ] github: parse: %s\n", err.Error())
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "github: parse: "+err.Error(), 0), err, http.StatusBadRequest, "GetGithubInfo")
		return
	}

	db.Create(&github)
	//db.Model(&models.Socials{}).Where("ID = ?", id).Update("GithubID", github.ID)
	db.Model(&models.Users{}).Joins("info ON users.info_id = info.id").Joins("socials ON info.socials_id = socials.id").Where("users.id = ?", id).Update("socials.github_id", github.ID)

	http.Redirect(w, r, "https://fuckme.hacktues.bg/", http.StatusMovedPermanently)
}

func LoginGithub(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	query := r.URL.Query()
	code := query.Get("code")

	client_id := "4f5f1918bf58eb0cccd4"
	client_secret := githubClientSecret

	if code == "" {
		//fmt.Println("Github: code is empty")
		//w.WriteHeader(http.StatusBadRequest)
		//w.Write([]byte("Github: code is empty"))
		fmt.Printf("[ ERROR ][ LoginGithub ] Github: code is empty\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "Github: code is empty", 0), nil, http.StatusBadRequest, "LoginGithub")
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
		//fmt.Println("Github: Error while getting token", err)
		//w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte("Github: Error while getting token" + err.Error()))
		fmt.Printf("[ ERROR ][ LoginGithub ] Github: Error while getting token: %s\n", err.Error())
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "Github: Error while getting token"+err.Error(), 0), err, http.StatusInternalServerError, "LoginGithub")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			//fmt.Println("Github: Error while closing body", err)
			//w.WriteHeader(http.StatusInternalServerError)
			//w.Write([]byte("Github: Error while closing body" + err.Error()))
			fmt.Printf("[ ERROR ][ LoginGithub ] Github: Error while closing body: %s\n", err.Error())
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "Github: Error while closing body"+err.Error(), 0), err, http.StatusInternalServerError, "LoginGithub")
			return
		}
	}(resp.Body)

	respString, _ := io.ReadAll(resp.Body)

	access_token := strings.Split(strings.Split(string(respString), "=")[1], "&")[0]

	bearerString := "Bearer " + access_token

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		//fmt.Println("Github: Error on request:", err)
		//w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte("Github: Error on request: " + err.Error()))
		fmt.Printf("[ ERROR ][ LoginGithub ] Github: Error on request: %s\n", err.Error())
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "Github: Error on request: "+err.Error(), 0), err, http.StatusInternalServerError, "LoginGithub")
		return
	}

	req.Header.Set("Authorization", bearerString)

	client := &http.Client{}
	resps, err := client.Do(req)
	if err != nil {
		//fmt.Println("Github: Error on response.\n[ERRO] -", err)
		//w.WriteHeader(http.StatusBadRequest)
		//w.Write([]byte("Github: Error on response.\n[ERRO] -" + err.Error()))
		fmt.Printf("[ ERROR ][ LoginGithub ] Github: Error on response: %s\n", err.Error())
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "Github: Error on response: "+err.Error(), 0), err, http.StatusBadRequest, "LoginGithub")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			//fmt.Println("Github: Error while closing body", err)
			//w.WriteHeader(http.StatusInternalServerError)
			//w.Write([]byte("Github: Error while closing body" + err.Error()))
			fmt.Printf("[ ERROR ][ LoginGithub ] Github: Error while closing body: %s\n", err.Error())
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "Github: Error while closing body"+err.Error(), 0), err, http.StatusInternalServerError, "LoginGithub")
			return
		}
	}(resps.Body)

	github := &models.Github{}
	user := &models.Users{}

	err = json.NewDecoder(resps.Body).Decode(&github)
	if err != nil {
		//fmt.Println("github: parse: ", err)
		//w.WriteHeader(http.StatusBadRequest)
		//w.Write([]byte("github: parse: " + err.Error()))
		fmt.Printf("[ ERROR ][ LoginGithub ] Github: parse: %s\n", err.Error())
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "Github: parse: "+err.Error(), 0), err, http.StatusBadRequest, "LoginGithub")
		return
	}

	db.Table("users").Joins("JOIN info ON users.info_id = info.id").Joins("JOIN socials ON info.socials_id = socials.id").Joins("JOIN github ON socials.github_id = github.id").Where("github.github_user_id = ?", github.GithubUserID).First(&user)

	if user.ID == 0 {
		//fmt.Println("Github: User not found")
		//w.WriteHeader(http.StatusBadRequest)
		//w.Write([]byte("Github: User not found"))
		fmt.Printf("[ ERROR ][ LoginGithub ] Github: User not found\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "Github: User not found", 0), err, http.StatusBadRequest, "LoginGithub")
		return
	}

	db.Table("users").Where("id = ?", user.ID).Update("last_login", time.Now())

	accessCookie, refreshCookie, err := jwt.GenerateCookies(user.ID)
	if err != nil {
		//fmt.Println("Github: Error while generating tokens", err)
		//w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte("Github: Error while generating tokens" + err.Error()))
		fmt.Printf("[ ERROR ][ LoginGithub ] Github: Error while generating tokens: %s\n", err.Error())
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "Github: Error while generating tokens"+err.Error(), 0), err, http.StatusInternalServerError, "LoginGithub")
		return
	}

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)
	http.Redirect(w, r, "https://fuckme.hacktues.bg/", http.StatusFound)
}

func LoginDiscord(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	query := r.URL.Query()
	code := query.Get("code")
	hostUrl := os.Getenv("HOST_URL")

	clientId := "1009547623637712977"
	clientSecret := discordClientSecret
	if code == "" {
		//fmt.Println("Discord: code is empty")
		//w.WriteHeader(http.StatusBadRequest)
		//w.Write([]byte("Discord: code is empty"))
		fmt.Printf("[ ERROR ][ LoginDiscord ] Discord: code is empty\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "Discord: code is empty", 0), nil, http.StatusBadRequest, "LoginDiscord")
		return
	}

	data := url.Values{
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {hostUrl + "api/auth/discord"},
		"scope":         {"identify"},
	}

	postBody := bytes.NewBufferString(data.Encode())

	resp, err := http.Post("https://discord.com/api/v10/oauth2/token", "application/x-www-form-urlencoded", postBody)
	if err != nil {
		//fmt.Println("Discord: Error while getting token", err)
		//w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte("Discord: Error while getting token" + err.Error()))
		fmt.Printf("[ ERROR ][ LoginDiscord ] Discord: Error while getting token: %s\n", err.Error())
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "Discord: Error while getting token"+err.Error(), 0), err, http.StatusInternalServerError, "LoginDiscord")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			//fmt.Println("Discord: Error while closing body", err)
			//w.WriteHeader(http.StatusInternalServerError)
			//w.Write([]byte("Discord: Error while closing body" + err.Error()))
			fmt.Printf("[ ERROR ][ LoginDiscord ] Discord: Error while closing body: %s\n", err.Error())
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "Discord: Error while closing body"+err.Error(), 0), err, http.StatusInternalServerError, "LoginDiscord")
			return
		}
	}(resp.Body)

	bearer := &models.DiscordBearer{}
	err = json.NewDecoder(resp.Body).Decode(&bearer)
	if err != nil {
		//fmt.Println("Discord: bearer decode error", err)
		//w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte("Discord: bearer decode error" + err.Error()))
		fmt.Printf("[ ERROR ][ LoginDiscord ] Discord: bearer decode error: %s\n", err.Error())
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "Discord: bearer decode error"+err.Error(), 0), err, http.StatusInternalServerError, "LoginDiscord")
		return
	}

	var bearerString = "Bearer " + bearer.AccessToken

	req, err := http.NewRequest("GET", "https://discord.com/api/v10/users/@me", nil)
	if err != nil {
		//fmt.Println("Discord: Error on request:", err)
		//w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte("Discord: Error on request: " + err.Error()))
		fmt.Printf("[ ERROR ][ LoginDiscord ] Discord: Error on request: %s\n", err.Error())
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "Discord: Error on request: "+err.Error(), 0), err, http.StatusInternalServerError, "LoginDiscord")

		return
	}

	req.Header.Set("Authorization", bearerString)

	client := &http.Client{}
	resps, err := client.Do(req)
	if err != nil {
		//fmt.Println("Discord: Error on response.\n[ERRO] -", err)
		//w.WriteHeader(http.StatusBadRequest)
		//w.Write([]byte("Discord: Error on response.\n[ERRO] -" + err.Error()))
		fmt.Printf("[ ERROR ][ LoginDiscord ] Discord: Error on response: %s\n", err.Error())
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "Discord: Error on response: "+err.Error(), 0), err, http.StatusBadRequest, "LoginDiscord")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			//fmt.Println("Discord: Error while closing body", err)
			//w.WriteHeader(http.StatusInternalServerError)
			//w.Write([]byte("Discord: Error while closing body" + err.Error()))
			fmt.Printf("[ ERROR ][ LoginDiscord ] Discord: Error while closing body: %s\n", err.Error())
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "Discord: Error while closing body"+err.Error(), 0), err, http.StatusInternalServerError, "LoginDiscord")
			return
		}
	}(resps.Body)

	discord := &models.Discord{}
	user := &models.Users{}

	err = json.NewDecoder(resps.Body).Decode(&discord)
	if err != nil {
		//fmt.Println("discord: parse: ", err)
		//w.WriteHeader(http.StatusBadRequest)
		//w.Write([]byte("discord: parse: " + err.Error()))
		fmt.Printf("[ ERROR ][ LoginDiscord ] discord: parse: %s\n", err.Error())
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "discord: parse: "+err.Error(), 0), err, http.StatusBadRequest, "LoginDiscord")
		return
	}

	db.Table("users").Joins("JOIN info ON users.info_id = info.id").Joins("JOIN socials ON info.socials_id = socials.id").Joins("JOIN discord ON socials.discord_id = discord.id").Where("discord.discord_user_id = ?", discord.DiscordUserID).First(&user)

	if user.ID == 0 {
		//fmt.Println("Discord: User not found")
		//w.WriteHeader(http.StatusBadRequest)
		//w.Write([]byte("Discord: User not found"))
		fmt.Printf("[ ERROR ][ LoginDiscord ] Discord: User not found\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "Discord: User not found", 0), err, http.StatusBadRequest, "LoginDiscord")
		return
	}

	db.Table("users").Where("id = ?", user.ID).Update("last_login", time.Now())

	accessCookie, refreshCookie, err := jwt.GenerateCookies(user.ID)
	if err != nil {
		//fmt.Println("Discord: Error while generating tokens", err)
		//w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte("Discord: Error while generating tokens" + err.Error()))
		fmt.Printf("[ ERROR ][ LoginDiscord ] Discord: Error while generating tokens: %s\n", err.Error())
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "Discord: Error while generating tokens"+err.Error(), 0), err, http.StatusInternalServerError, "LoginDiscord")
		return
	}

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)

	http.Redirect(w, r, "https://fuckme.hacktues.bg/", http.StatusFound)

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
