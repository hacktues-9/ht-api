package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
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

	err := json.NewDecoder(r.Body).Decode(&parseUser)
	if err != nil {
		fmt.Printf("[ ERROR ] [ Register ] registerUser: parse: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "registerUser: parse: "+err.Error(), 0), err, http.StatusInternalServerError, "Register")
		return
	}

	//check if email is valid elsys email format [name].[middle name initial].[surname].[yearofEntry]@elsys-bg.org
	//split email by . and @
	splitEmail := strings.Split(parseUser.ElsysEmail, ".")
	splitEmail2 := strings.Split(splitEmail[3], "@")

	//check if splitEmail2[0] is a year after 2018
	year, err := strconv.Atoi(splitEmail2[0])
	if err != nil {
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "email is not valid elsys email", 0), nil, http.StatusUnauthorized, "Register")
		return
	}

	if year < 2018 {
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "email is not authorized to participate", 0), nil, http.StatusUnauthorized, "Register")
		return
	}

	//combine splitEmail2[1] and splitEmail[4]
	splitEmail2[1] = splitEmail2[1] + "." + splitEmail[4]

	if splitEmail2[1] != "elsys-bg.org" {
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "email is not valid elsys email", 0), nil, http.StatusUnauthorized, "Register")
		return
	}

	classID, eatingPreferenceID, shirtSizeID, roleID, allergies, technologies := returnDefaultIDs(db, &parseUser)

	password, err := pass.HashPassword(parseUser.Password)

	if err != nil {
		fmt.Printf("[ ERROR ] [ Register ] password: hash: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "password: hash: "+err.Error(), 0), err, http.StatusInternalServerError, "Register")
		return
	}

	userSocials = models.Socials{
		LinkedInLink:   "",
		InstagramLink:  "",
		ProfilePicture: "https://api.hacktues.bg/api/image/" + parseUser.FirstName + " " + parseUser.LastName,
	}

	if result := db.Omit("DiscordID", "GithubID").Create(&userSocials); result.Error != nil {
		fmt.Printf("[ ERROR ] [ Register ] userSocials: create: %v\n", result.Error)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "userSocials: create: "+result.Error.Error(), 0), result.Error, http.StatusInternalServerError, "Register")
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
		//delete prev Tables
		db.Delete(&userSocials)
		fmt.Printf("[ ERROR ] [ Register ] userInfo: create: %v\n", result.Error)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "userInfo: create: "+result.Error.Error(), 0), result.Error, http.StatusInternalServerError, "Register")
		return
	}

	if result := db.Create(&userSecurity); result.Error != nil {
		//delete prev Tables
		db.Delete(&userSocials)
		db.Delete(&userInfo)
		fmt.Printf("[ ERROR ] [ Register ] userSecurity: create: %v\n", result.Error)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "userSecurity: create: "+result.Error.Error(), 0), result.Error, http.StatusInternalServerError, "Register")
		return
	}

	user = models.Users{
		FirstName:      parseUser.FirstName,
		LastName:       parseUser.LastName,
		Email:          parseUser.Email,
		ElsysEmail:     parseUser.ElsysEmail,
		Mobile:         parseUser.Mobile,
		Password:       password,
		InfoID:         userInfo.ID,
		SecurityID:     userSecurity.ID,
		RoleID:         roleID,
		LookingForTeam: true,
	}
	verificationLinkTTL := time.Duration(24) * time.Hour
	deletionLinkTTL := time.Duration(256) * time.Hour

	if user.Email != "" {
		verificationLink := email.GenerateVerificationLink(parseUser.Email, accessTokenPrivateKey, accessTokenPublicKey, verificationLinkTTL)
		deletionLink := email.GenerateDeletionLink(parseUser.Email, accessTokenPrivateKey, accessTokenPublicKey, deletionLinkTTL)
		err = email.SendEmail(user.FirstName+" "+user.LastName, user.Email, verificationLink, deletionLink)
		if err != nil {
			fmt.Printf("[ ERROR ] [ Register ] email: send: %v\n", err)
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "email: send: "+err.Error(), 0), err, http.StatusInternalServerError, "Register")
			return
		}
	}

	verificationLink := email.GenerateVerificationLink(parseUser.ElsysEmail, accessTokenPrivateKey, accessTokenPublicKey, verificationLinkTTL)
	deletionLink := email.GenerateDeletionLink(parseUser.ElsysEmail, accessTokenPrivateKey, accessTokenPublicKey, deletionLinkTTL)
	err = email.SendEmail(user.FirstName+" "+user.LastName, user.ElsysEmail, verificationLink, deletionLink)
	if err != nil {
		fmt.Printf("[ ERROR ] [ Register ] email: send: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "email: send: "+err.Error(), 0), err, http.StatusInternalServerError, "Register")
		return
	}

	if result := db.Omit("TeamID").Create(&user); result.Error != nil {
		//delete prev Tables
		db.Delete(&userSocials)
		db.Delete(&userInfo)
		db.Delete(&userSecurity)
		fmt.Printf("[ ERROR ] [ Register ] user: create: %v\n", result.Error)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "user: create: "+result.Error.Error(), 0), result.Error, http.StatusInternalServerError, "Register")
		return
	}

	for _, allergy := range allergies {
		db.Create(&models.InfoAllergies{InfoID: userInfo.ID, AllergiesID: allergy})
	}

	for _, tech := range technologies {
		db.Create(&models.UserTechnologies{UserID: user.ID, TechnologiesID: tech})
	}

	accessToken, err := jwt.CreateToken(accessTokenTTL, user.ID, accessTokenPrivateKey, accessTokenPublicKey)
	if err != nil {
		fmt.Printf("[ ERROR ] [ Register ] access token: create: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "access token: create: "+err.Error(), 0), err, http.StatusInternalServerError, "Register")
		return
	}

	refreshToken, err := jwt.CreateToken(refreshTokenTTL, user.ID, refreshTokenPrivateKey, refreshTokenPublicKey)
	if err != nil {
		fmt.Printf("[ ERROR ] [ Register ] refresh token: create: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "refresh token: create: "+err.Error(), 0), err, http.StatusInternalServerError, "Register")
		return
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

	http.SetCookie(w, &refreshCookie)
	http.SetCookie(w, &accessCookie)

	models.RespHandler(w, r, models.DefaultPosResponse(strconv.Itoa(int(user.ID))), nil, http.StatusOK, "Register")
}

func Login(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	user := models.LoginUser{}
	userDB := models.Users{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Printf("[ ERROR ] [ Login ] loginUser: parse: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "loginUser: parse: "+err.Error(), 0), err, http.StatusInternalServerError, "Login")
		return
	}

	if result := db.Where("email = ?", user.Identifier).First(&userDB); result.Error != nil {
		if result := db.Where("elsys_email = ?", user.Identifier).First(&userDB); result.Error != nil {
			fmt.Printf("[ ERROR ] [ Login ] user: find: %v\n", result.Error)
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusNotFound, "user: find: "+result.Error.Error(), 0), result.Error, http.StatusNotFound, "Login")
			return
		}
	}

	if !pass.ComparePasswords(userDB.Password, user.Password) {
		fmt.Printf("[ ERROR ] [ Login ] password: compare: wrong password")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "password: compare: wrong password", 0), errors.New("password: compare: wrong password"), http.StatusUnauthorized, "Login")
		return
	}

	db.Model(&models.Users{}).Where("ID = ?", userDB.ID).Update("last_login", time.Now())

	accessToken, err := jwt.CreateToken(accessTokenTTL, userDB.ID, accessTokenPrivateKey, accessTokenPublicKey)
	if err != nil {
		fmt.Printf("[ ERROR ] [ Login ] access token: create: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "access token: create: "+err.Error(), 0), err, http.StatusInternalServerError, "Login")
		return
	}

	refreshToken, err := jwt.CreateToken(refreshTokenTTL, userDB.ID, refreshTokenPrivateKey, refreshTokenPublicKey)
	if err != nil {
		fmt.Printf("[ ERROR ] [ Login ] refresh token: create: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "refresh token: create: "+err.Error(), 0), err, http.StatusInternalServerError, "Login")
		return
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

	http.SetCookie(w, &refreshCookie)
	http.SetCookie(w, &accessCookie)

	models.RespHandler(w, r, models.DefaultPosResponse("success"), nil, http.StatusOK, "Login")
}

func GetUserID(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	sub, err := ReturnAuthID(w, r, db)
	if err != nil {
		fmt.Printf("[ ERROR ] [ GetUserID ] %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, err.Error(), 0), err, http.StatusUnauthorized, "GetUserID")
		return
	}

	//check if user exists
	user := models.Users{}
	if result := db.Where("ID = ?", sub).First(&user); result.Error != nil {
		//delete cookies
		Logout(w)
		fmt.Printf("[ ERROR ] [ GetUserID ] user: find: %v\n", result.Error)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusNotFound, "user: find: "+result.Error.Error(), 0), result.Error, http.StatusNotFound, "GetUserID")
		return
	}

	models.RespHandler(w, r, models.DefaultPosResponse(sub), nil, http.StatusOK, "GetUserID")
}

func GetUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// get id int from query
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Printf("[ ERROR ] [ GetUser ] id: parse: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "id: parse: "+err.Error(), 0), err, http.StatusBadRequest, "GetUser")
		return
	}
	sub, err := ReturnAuthID(w, r, db)
	if err != nil {
		fmt.Printf("[ ERROR ] [ GetUser ] %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, err.Error(), 0), err, http.StatusUnauthorized, "GetUser")
		return
	}

	if float64(sub) != float64(id) {
		fmt.Printf("[ ERROR ] [ GetUser ] access token: validate: wrong id %v %v\n", sub, id)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "access token: validate: wrong id", 0), err, http.StatusUnauthorized, "GetUser")
		return
	}

	var user models.UserView

	db.Raw("SELECT * FROM userview(?)", id).Scan(&user)

	if user.FirstName == "" {
		fmt.Printf("[ ERROR ] [ GetUser ] user: find: not found %v\n", user)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusNotFound, "user: find: not found", 0), err, http.StatusNotFound, "GetUser")
		return
	}

	//get Technologies for user
	var technologies []models.Technologies
	var technologiesView []string
	db.Table("technologies").Joins("JOIN user_technologies ON user_technologies.technologies_id = technologies.id").Where("user_technologies.user_id = ?", id).Scan(&technologies)
	for _, technology := range technologies {
		technologiesView = append(technologiesView, technology.Technology)
	}

	user.Technologies = technologiesView

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		fmt.Printf("[ ERROR ] [ GetUser ] json encode: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "json encode: "+err.Error(), 0), err, http.StatusInternalServerError, "GetUser")
		return
	}
}

func Logout(w http.ResponseWriter) {
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

func UpdateUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var parseChangeUser models.ParseChangeUser

	err := json.NewDecoder(r.Body).Decode(&parseChangeUser)
	if err != nil {
		fmt.Printf("[ ERROR ] [ UpdateUser ] json decode: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "json decode: "+err.Error(), 0), err, http.StatusInternalServerError, "UpdateUser")
		return
	}

	sub, err := ReturnAuthID(w, r, db)
	if err != nil {
		fmt.Printf("[ ERROR ] [ UpdateUser ] %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, err.Error(), 0), err, http.StatusUnauthorized, "UpdateUser")
		return
	}

	if float64(sub) != float64(parseChangeUser.ID) {
		fmt.Printf("[ ERROR ] [ UpdateUser ] access token: validate: wrong id %v %v\n", sub, parseChangeUser.ID)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "access token: validate: wrong id"+strconv.Itoa(int(sub))+" "+strconv.Itoa(int(parseChangeUser.ID)), 0), err, http.StatusUnauthorized, "UpdateUser")
	}

	//to change technologies we need to delete all technologies for user and add new ones

	db.Exec("DELETE FROM user_technologies WHERE user_id = ?", sub)
	for _, technology := range parseChangeUser.Technologies {
		var technologyID int
		db.Table("technologies").Where("technology = ?", technology).Select("id").Row().Scan(&technologyID)
		db.Create(&models.UserTechnologies{UserID: sub, TechnologiesID: uint(technologyID)})
	}

	//change LookingForTeam
	db.Model(&models.Users{}).Where("id = ?", sub).Update("looking_for_team", parseChangeUser.LookingForTeam)

	//get first name and last name
	var firstName, lastName string
	db.Table("users").Where("id = ?", sub).Select("first_name, last_name").Row().Scan(&firstName, &lastName)

	//update profile_picture
	db.Model(&models.Users{}).Joins("JOIN info ON users.info_id = info.id").Joins("JOIN socials ON info.socials_id = socials.id").Where("users.id = ?", sub).Update("socials.profile_picture", "https://api.hacktues.bg/api/image/"+firstName+"%20"+lastName)

	models.RespHandler(w, r, models.DefaultPosResponse("success"), nil, http.StatusOK, "UpdateUser")
}

func ForgotPassword(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	mail := mux.Vars(r)["elsys_email"]

	if mail == "" {
		fmt.Printf("[ ERROR ] [ ForgotPassword ] mail: empty\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "mail: empty", 0), errors.New(" "), http.StatusBadRequest, "ForgotPassword")
		return
	}

	var user models.Users
	db.Table("users").Where("elsys_email = ?", mail).Scan(&user)

	if user.ID == 0 {
		fmt.Printf("[ ERROR ] [ ForgotPassword ] user: find: not found %v\n", user)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusNotFound, "user: find: not found", 0), errors.New(" "), http.StatusNotFound, "ForgotPassword")
		return
	}
	resetLinkTTL := time.Duration(24) * time.Hour

	resetLink := email.GenerateResetLink(mail, accessTokenPrivateKey, accessTokenPublicKey, resetLinkTTL)

	err := email.SendResetLink(user.FirstName+" "+user.LastName, mail, resetLink)
	if err != nil {
		fmt.Printf("[ ERROR ] [ ForgotPassword ] send email: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "send email: "+err.Error(), 0), err, http.StatusInternalServerError, "ForgotPassword")
		return
	}

	models.RespHandler(w, r, models.DefaultPosResponse("success"), nil, http.StatusOK, "ForgotPassword")
}

func ResetPassword(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var parseReset models.ParseReset

	token := mux.Vars(r)["token"]
	err := json.NewDecoder(r.Body).Decode(&parseReset)
	if err != nil {
		fmt.Printf("[ ERROR ] [ ResetPassword ] json decode: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "json decode: "+err.Error(), 0), err, http.StatusInternalServerError, "ResetPassword")
		return
	}

	//validate reset link
	elsysEmail, err := email.ValidateResetLink(token)
	if err != nil {
		fmt.Printf("[ ERROR ] [ ResetPassword ] validate reset link: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "validate reset link: "+err.Error(), 0), err, http.StatusBadRequest, "ResetPassword")
		return
	}

	//get user id
	var userID int
	db.Table("users").Where("elsys_email = ?", elsysEmail).Select("id").Row().Scan(&userID)

	if userID == 0 {
		fmt.Printf("[ ERROR ] [ ResetPassword ] user: find: not found %v\n", userID)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusNotFound, "user: find: not found", 0), errors.New(" "), http.StatusNotFound, "ResetPassword")
		return
	}

	//hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(parseReset.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("[ ERROR ] [ ResetPassword ] hash password: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "hash password: "+err.Error(), 0), err, http.StatusInternalServerError, "ResetPassword")
		return
	}

	//update password
	db.Model(&models.Users{}).Where("id = ?", userID).Update("password", string(hashedPassword))
}

func DeleteUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	token := mux.Vars(r)["token"]

	mail, err := email.ValidateEmailToken(token)
	if err != nil {
		fmt.Printf("[ ERROR ] [ ValidateEmail ] validate: token: %s", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "validate: token", 0), err, http.StatusUnauthorized, "ValidateEmail")
		return
	}

	var user models.Users
	db.Table("users").Where("elsys_email = ?", mail).Scan(&user)
	if user.ID == 0 {
		fmt.Printf("[ ERROR ] [ DeleteUser ] user: find: not found %v\n", user)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusNotFound, "user: find: not found", 0), errors.New(" "), http.StatusNotFound, "DeleteUser")
		return
	}

	db.Delete(&user)

	models.RespHandler(w, r, models.DefaultPosResponse("success"), nil, http.StatusOK, "DeleteUser")
}

func GetUserRole(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// get id int from query
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Printf("[ ERROR ] [ GetUser ] id: parse: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "id: parse: "+err.Error(), 0), err, http.StatusBadRequest, "GetUser")
		return
	}
	sub, err := ReturnAuthID(w, r, db)
	if err != nil {
		fmt.Printf("[ ERROR ] [ GetUser ] %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, err.Error(), 0), err, http.StatusUnauthorized, "GetUser")
		return
	}

	if float64(sub) != float64(id) {
		fmt.Printf("[ ERROR ] [ GetUser ] access token: validate: wrong id %v %v\n", sub, id)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "access token: validate: wrong id", 0), err, http.StatusUnauthorized, "GetUser")
		return
	}

	var roleId int
	db.Table("users").Where("id = ?", id).Select("role_id").Row().Scan(&roleId)

	models.RespHandler(w, r, models.DefaultPosResponse(roleId), nil, http.StatusOK, "GetUserRole")
}
