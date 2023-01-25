package users

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hacktues-9/API/pkg/email"
	"github.com/hacktues-9/API/pkg/jwt"
	"github.com/hacktues-9/API/pkg/models"
	"github.com/hacktues-9/API/pkg/pass"
	"gorm.io/gorm"
)

var (
	accessTokenTTL         = time.Hour * 24
	refreshTokenTTL        = time.Hour * 24 * 7
	accessTokenPrivateKey  = os.Getenv("ACCESS_TOKEN_PRIVATE_KEY")
	accessTokenPublicKey   = os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")
	refreshTokenPrivateKey = os.Getenv("REFRESH_TOKEN_PRIVATE_KEY")
	refreshTokenPublicKey  = os.Getenv("REFRESH_TOKEN_PUBLIC_KEY")
)

func Register(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	userInfo := models.Info{}
	user := models.Users{}
	parseUser := models.RegisterUser{}
	userSocials := models.Socials{}
	userSecurity := models.Security{}
	//TODO : Anti Radoslav Filipov func
	//TODO : Anti Krum Stefanov func
	//TODO : Anti David ot G class func
	//TODO : Anti Ivan Ivanov func
	//TODO : Anti Ivan Georgiev func
	//TODO : Anti Vasil Kolev func
	//TODO : Anti C-- func
	//TODO : Anti Petyo Miladinov func

	err := json.NewDecoder(r.Body).Decode(&parseUser)
	if err != nil {
		fmt.Println("register: registerUser: parse:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("register: registerUser: parse: " + err.Error()))
		return
	}

	classID, eatingPreferenceID, shirtSizeID, roleID, allergies, technologies := returnDefaultIDs(db, &parseUser)

	password, err := pass.HashPassword(parseUser.Password)

	if err != nil {
		fmt.Println("register: password: hash:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("register: password: hash: " + err.Error()))
		return
	}

	userSocials = models.Socials{
		LinkedInLink:  "",
		InstagramLink: "",
	}

	if result := db.Omit("DiscordID", "GithubID").Create(&userSocials); result.Error != nil {
		fmt.Println("register: userSocials: create:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("register: userSocials: create: " + result.Error.Error()))
		return
	}

	grade, _ := strconv.Atoi(parseUser.Grade)

	userInfo = models.Info{
		Grade:              grade,
		ClassID:            classID,
		EatingPreferenceID: eatingPreferenceID,
		SocialsID:          userSocials.ID,
		ShirtSizeID:        shirtSizeID,
	}

	if result := db.Create(&userInfo); result.Error != nil {
		fmt.Println("register: userInfo: create:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("register: userInfo: create: " + result.Error.Error()))
		return
	}

	if result := db.Create(&userSecurity); result.Error != nil {
		fmt.Println("register: userSecurity: create:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("register: userSecurity: create: " + result.Error.Error()))
		return
	}

	user = models.Users{
		FirstName:  parseUser.FirstName,
		LastName:   parseUser.LastName,
		Email:      parseUser.Email,
		ElsysEmail: parseUser.ElsysEmail,
		Mobile:     parseUser.Mobile,
		Password:   password,
		InfoID:     userInfo.ID,
		SecurityID: userSecurity.ID,
		RoleID:     roleID,
	}
	verificationLinkTTL := time.Duration(24) * time.Hour

	if user.Email != "" {
		verificationLink := email.GenerateVerificationLink(parseUser.Email, accessTokenPrivateKey, accessTokenPublicKey, verificationLinkTTL)
		email.SendEmail(user.FirstName+" "+user.LastName, user.Email, verificationLink)
	}

	verificationLink := email.GenerateVerificationLink(parseUser.ElsysEmail, accessTokenPrivateKey, accessTokenPublicKey, verificationLinkTTL)
	email.SendEmail(user.FirstName+" "+user.LastName, user.ElsysEmail, verificationLink)

	if result := db.Omit("TeamID").Create(&user); result.Error != nil {
		fmt.Println("register: user: create:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("register: user: create: " + result.Error.Error()))
		return
	}

	for _, allergy := range allergies {
		db.Create(&models.InfoAllergies{InfoID: userInfo.ID, AllergiesID: allergy})
	}

	for _, tech := range technologies {
		db.Create(&models.UserTechnologies{UserID: user.ID, TechnologiesID: tech})
	}

	resp := ParseUser(user.ID, db)

	accessToken, err := jwt.CreateToken(accessTokenTTL, user.ID, accessTokenPrivateKey, accessTokenPublicKey)
	if err != nil {
		fmt.Println("register: access token: create:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("register: access token: create: " + err.Error()))
		return
	}

	refreshToken, err := jwt.CreateToken(refreshTokenTTL, user.ID, refreshTokenPrivateKey, refreshTokenPublicKey)
	if err != nil {
		fmt.Println("register: refresh token: create:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("register: refresh token: create: " + err.Error()))
		return
	}

	refresh_cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(refreshTokenTTL),
		HttpOnly: true,
		Domain:   os.Getenv("HOST_DOMAIN"),
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	access_cookie := http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(accessTokenTTL),
		HttpOnly: true,
		Domain:   os.Getenv("HOST_DOMAIN"),
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	http.SetCookie(w, &refresh_cookie)
	http.SetCookie(w, &access_cookie)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func Login(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	user := models.LoginUser{}
	userDB := models.Users{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("login: loginUser: parse:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("login: loginUser: parse: " + err.Error()))
		return
	}

	if result := db.Where("email = ?", user.Identifier).First(&userDB); result.Error != nil {
		if result := db.Where("elsys_email = ?", user.Identifier).First(&userDB); result.Error != nil {
			fmt.Println("login: user: find:", result.Error)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("login: user: find: " + result.Error.Error()))
			return
		}
	}

	if !pass.ComparePasswords(userDB.Password, user.Password) {
		fmt.Println("login: password: compare: wrong password")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("login: password: compare: wrong password"))
		return
	}

	db.Model(&models.Users{}).Where("ID = ?", userDB.ID).Update("last_login", time.Now())

	accessToken, err := jwt.CreateToken(accessTokenTTL, userDB.ID, accessTokenPrivateKey, accessTokenPublicKey)
	if err != nil {
		fmt.Println("login: access token: create:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("login: access token: create: " + err.Error()))
		return
	}

	refreshToken, err := jwt.CreateToken(refreshTokenTTL, userDB.ID, refreshTokenPrivateKey, refreshTokenPublicKey)
	if err != nil {
		fmt.Println("login: refresh token: create:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("login: refresh token: create: " + err.Error()))
		return
	}

	refresh_cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(refreshTokenTTL),
		HttpOnly: true,
		Domain:   os.Getenv("HOST_DOMAIN"),
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	access_cookie := http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(accessTokenTTL),
		HttpOnly: true,
		Domain:   os.Getenv("HOST_DOMAIN"),
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	http.SetCookie(w, &refresh_cookie)
	http.SetCookie(w, &access_cookie)
	http.Redirect(w, r, "http://localhost:3000/", http.StatusFound)
}

func GetUserID(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	cookie, err := r.Cookie("access_token")
	authorizationHeader := r.Header.Get("Authorization")
	fields := strings.Fields(authorizationHeader)
	accessToken := ""

	if len(fields) != 0 && fields[0] == "Bearer" {
		accessToken = fields[1]
	} else if err == nil {
		accessToken = cookie.Value
	} else {
		fmt.Println("get user: access token: get:", err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("get user: access token: get: " + err.Error()))
		return
	}

	sub, err := jwt.ValidateToken(accessToken, accessTokenPublicKey)
	if err != nil {
		fmt.Println("get user: access token: validate:", err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("get user: access token: validate: " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"id\": \"" + fmt.Sprintf("%v", sub) + "\"}"))
}

func GetUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// get id int from query
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println("get user: id: parse:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("get user: id: parse: " + err.Error()))
		return
	}
	cookie, err := r.Cookie("access_token")
	authorizationHeader := r.Header.Get("Authorization")
	fields := strings.Fields(authorizationHeader)
	accessToken := ""

	if len(fields) != 0 && fields[0] == "Bearer" {
		accessToken = fields[1]
	} else if err == nil {
		accessToken = cookie.Value
	} else {
		fmt.Println("get user: access token: get:", err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("get user: access token: get: " + err.Error()))
		return
	}

	sub, err := jwt.ValidateToken(accessToken, accessTokenPublicKey)
	if err != nil {
		fmt.Println("get user: access token: validate:", err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("get user: access token: validate: " + err.Error()))
		return
	}

	if sub.(float64) != float64(id) {
		fmt.Println("get user: access token: validate: wrong id", sub, id)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("get user: access token: validate: wrong id"))
		return
	}

	var user models.UserView

	db.Raw("SELECT * FROM userview(?)", id).Scan(&user)

	if user.FirstName == "" {
		fmt.Println("get user: user: find: not found", user)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("get user: user: find: not found"))
		return
	}

	//get Technologies for user
	var technologies []models.Technologies
	var technologiesView []string
	db.Table("technologies").Joins("JOIN user_technologies ON user_technologies.technology_id = technologies.id").Where("user_technologies.user_id = ?", id).Scan(&technologies)
	for _, technology := range technologies {
		technologiesView = append(technologiesView, technology.Technology)
	}

	user.Technologies = technologiesView

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func Logout(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Domain:   os.Getenv("HOST_DOMAIN"),
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	accessCookie := http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Domain:   os.Getenv("HOST_DOMAIN"),
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	http.SetCookie(w, &refreshCookie)
	http.SetCookie(w, &accessCookie)
	w.WriteHeader(http.StatusOK)
}
