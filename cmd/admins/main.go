package admins

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	users2 "github.com/hacktues-9/API/cmd/users"
	"github.com/hacktues-9/API/pkg/email"
	"github.com/hacktues-9/API/pkg/models"
	"gorm.io/gorm"
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

func GetTeams(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// get teams with 3 to 5 members
	var count int64

	db.Raw("SELECT * FROM team WHERE id IN (SELECT team_id FROM users WHERE deleted_at IS NOT NULL GROUP BY team_id HAVING COUNT(*) >= 3 AND COUNT(*) <= 5)").Count(&count)

	models.RespHandler(w, r, models.DefaultPosResponse(count), nil, http.StatusOK, "AdminGetTeams")
}

func SearchWithFilters(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	sub, err := users2.ReturnAuthID(w, r, db)
	if err != nil {
		fmt.Printf("[ ERROR ] [ SearchWithFilters ] auth: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "auth: "+err.Error(), 0), err, http.StatusInternalServerError, "AdminSearchWithFilters")
		return
	}

	//get roleId
	var roleId int
	db.Table("users").Select("role_id").Where("id = ?", sub).Find(&roleId)
	if roleId != 5 {
		fmt.Printf("[ ERROR ] [ SearchWithFilters ] unauthorized: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "unauthorized", 0), err, http.StatusUnauthorized, "AdminSearchWithFilters")
		return
	}
	//parse request body
	var filters models.ParseFilterUsers
	err = json.NewDecoder(r.Body).Decode(&filters)
	if err != nil {
		fmt.Printf("[ ERROR ] [ SearchWithFilters ] parse: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "parse: "+err.Error(), 0), err, http.StatusInternalServerError, "AdminSearchWithFilters")
		return
	}

	// get users with filters
	var users []models.KurView

	db.Raw("select * from get_user_filtered(?,?,?,?,?,?,?,?,?)", filters.ShirtSize, filters.Grade, filters.Class, filters.Name, filters.Email, filters.Mobile, filters.ElsysEmail, filters.Team, filters.EatingPreference).Scan(&users)

	models.RespHandler(w, r, models.DefaultPosResponse(users), nil, http.StatusOK, "AdminSearchWithFilters")
}

func ResendVerification(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	sub, err := users2.ReturnAuthID(w, r, db) // admin id
	if err != nil {
		fmt.Printf("[ ERROR ] [ ResendVerification ] auth: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "auth: "+err.Error(), 0), err, http.StatusInternalServerError, "AdminResendVerification")
		return
	}
	id := mux.Vars(r)["id"] // user id

	//get roleId
	var roleId int
	db.Table("users").Select("role_id").Where("id = ?", sub).Find(&roleId)
	if roleId != 5 {
		fmt.Printf("[ ERROR ] [ ResendVerification ] unauthorized: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "unauthorized", 0), err, http.StatusUnauthorized, "AdminResendVerification")
		return
	}

	// check if user exists
	var user models.Users
	db.Table("users").Where("id = ?", id).Find(&user)
	if user.ID == 0 {
		fmt.Printf("[ ERROR ] [ ResendVerification ] user not found: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusNotFound, "user not found", 0), err, http.StatusNotFound, "AdminResendVerification")
		return
	}

	var security models.Security
	db.Table("security").Where("id = ?", user.SecurityID).Find(&security)
	if security.ID == 0 {
		fmt.Printf("[ ERROR ] [ ResendVerification ] security not found: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusNotFound, "security not found", 0), err, http.StatusNotFound, "AdminResendVerification")
		return
	}

	if security.EmailVerified {
		fmt.Printf("[ ERROR ] [ ResendVerification ] email already verified: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "email already verified", 0), err, http.StatusBadRequest, "AdminResendVerification")
		return
	}

	verificationLinkTTL := time.Duration(24) * time.Hour

	// send verification email
	verificationLink := email.GenerateVerificationLink(user.Email, accessTokenPrivateKey, accessTokenPublicKey, verificationLinkTTL)
	err = email.SendEmail(user.FirstName+" "+user.LastName, user.Email, verificationLink, "")
	if err != nil {
		fmt.Printf("[ ERROR ] [ Register ] email: send: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "email: send: "+err.Error(), 0), err, http.StatusInternalServerError, "Register")
		return
	}

	models.RespHandler(w, r, models.DefaultPosResponse("verification email sent"), nil, http.StatusOK, "AdminResendVerification")
}

func ResendVerificationElsys(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	sub, err := users2.ReturnAuthID(w, r, db) // admin id
	if err != nil {
		fmt.Printf("[ ERROR ] [ ResendVerificationElsys ] auth: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "auth: "+err.Error(), 0), err, http.StatusInternalServerError, "AdminResendVerificationElsys")
		return
	}

	id := mux.Vars(r)["id"] // user id

	//get roleId
	var roleId int
	db.Table("users").Select("role_id").Where("id = ?", sub).Find(&roleId)
	if roleId != 5 {
		fmt.Printf("[ ERROR ] [ ResendVerificationElsys ] unauthorized: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "unauthorized", 0), err, http.StatusUnauthorized, "AdminResendVerificationElsys")
		return
	}

	// check if user exists
	var user models.Users
	db.Table("users").Where("id = ?", id).Find(&user)
	if user.ID == 0 {
		fmt.Printf("[ ERROR ] [ ResendVerificationElsys ] user not found: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusNotFound, "user not found", 0), err, http.StatusNotFound, "AdminResendVerificationElsys")
		return
	}

	var security models.Security
	db.Table("security").Where("id = ?", user.SecurityID).Find(&security)
	if security.ID == 0 {
		fmt.Printf("[ ERROR ] [ ResendVerificationElsys ] security not found: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusNotFound, "security not found", 0), err, http.StatusNotFound, "AdminResendVerificationElsys")
		return
	}

	if security.ElsysEmailVerified {
		fmt.Printf("[ ERROR ] [ ResendVerificationElsys ] email already verified: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "email already verified", 0), err, http.StatusBadRequest, "AdminResendVerificationElsys")
		return
	}

	verificationLinkTTL := time.Duration(24) * time.Hour

	// send verification email
	verificationLink := email.GenerateVerificationLink(user.ElsysEmail, accessTokenPrivateKey, accessTokenPublicKey, verificationLinkTTL)
	err = email.SendEmail(user.FirstName+" "+user.LastName, user.ElsysEmail, verificationLink, "")
	if err != nil {
		fmt.Printf("[ ERROR ] [ Register ] email: send: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "email: send: "+err.Error(), 0), err, http.StatusInternalServerError, "Register")
		return
	}

	models.RespHandler(w, r, models.DefaultPosResponse("verification email sent"), nil, http.StatusOK, "AdminResendVerificationElsys")
}
