package users

import (
	"errors"
	"fmt"
	"github.com/disintegration/letteravatar"
	"github.com/gorilla/mux"
	"github.com/hacktues-9/API/pkg/jwt"
	"github.com/hacktues-9/API/pkg/models"
	"gorm.io/gorm"
	"image/png"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

func ReturnAuthID(w http.ResponseWriter, r *http.Request, db *gorm.DB) (uint, error) {
	cookie, err := r.Cookie("access_token")
	_, errR := r.Cookie("refresh_token")
	authorizationHeader := r.Header.Get("Authorization")
	fields := strings.Fields(authorizationHeader)
	accessToken := ""
	var sub uint

	if len(fields) != 0 && fields[0] == "Bearer" {
		accessToken = fields[1]
	} else if err == nil {
		accessToken = cookie.Value
	} else if errR == nil {
		jwt.RefreshAccessToken(w, r, db)
	} else {
		fmt.Println("get user: access token: get:", err)
		return 0, errors.New("no access token provided")
	}

	sub, err = jwt.ValidateToken(accessToken, accessTokenPublicKey)
	if err != nil {
		fmt.Println("get user: access token: validate:", err)
		return 0, err
	}

	return sub, nil
}

func returnDefaultIDs(db *gorm.DB, user *models.RegisterUser) (uint, uint, uint, uint, []uint, []uint) {
	class := models.Class{}
	eatingPreference := models.EatingPreference{}
	shirtSize := models.ShirtSize{}
	role := models.Role{}
	allergies := []uint{}
	technologies := []uint{}

	db.Where("name = ?", user.Class).First(&class)
	db.Where("name = ?", user.EatingPreference).First(&eatingPreference)
	db.Where("name = ?", user.ShirtSize).First(&shirtSize)
	db.Where("name = ?", "STUDENT").First(&role)

	for _, allergy := range user.Allergies {
		var tempAllergy models.Allergies
		db.Where("name = ?", allergy).First(&tempAllergy)
		allergies = append(allergies, tempAllergy.ID)
	}

	for _, tech := range user.Technologies {
		var tempTech models.Technologies
		db.Where("technology = ?", tech).First(&tempTech)
		technologies = append(technologies, tempTech.ID)
	}

	return class.ID, eatingPreference.ID, shirtSize.ID, role.ID, allergies, technologies
}

func GetNotifications(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	userId, err := ReturnAuthID(w, r, db)
	if err != nil {
		fmt.Printf("[ ERROR ] [ GetNotifications ] %v", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, err.Error(), 0), err, http.StatusUnauthorized, "GetNotifications")
		return
	}

	//var notifications []models.Notification
	//db.Table("invite").Where("user_id = ?", userId).Joins("JOIN team ON team.id = invite.team_id").Select("team.id, team.name, team.logo").Scan(&notifications)

	var notifications []models.Notification
	db.Raw("SELECT * FROM inviteview(?)", userId).Scan(&notifications)

	models.RespHandler(w, r, models.DefaultPosResponse(notifications), nil, http.StatusOK, "GetNotifications")
}

func GenerateImage(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	name := mux.Vars(r)["name"]

	//if images/ folder does not exist, create it
	if _, err := os.Stat("images"); os.IsNotExist(err) {
		err := os.Mkdir("images", 0755)
		if err != nil {
			fmt.Printf("[ ERROR ] [ GenerateImage ] %v", err)
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, err.Error(), 0), err, http.StatusInternalServerError, "GenerateImage")
			return
		}
	}

	if _, err := os.Stat("images/" + name + ".png"); err == nil {
		// file exists
		file, err := os.Open("images/" + name + ".png")
		if err != nil {
			fmt.Printf("[ ERROR ] [ GenerateImage ] %v", err)
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, err.Error(), 0), err, http.StatusInternalServerError, "GenerateImage")
			return
		}

		defer file.Close()

		//return image
		w.Header().Set("Content-Type", "image/png")
		_, err = io.Copy(w, file)
		if err != nil {
			fmt.Printf("[ ERROR ] [ GenerateImage ] %v", err)
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, err.Error(), 0), err, http.StatusInternalServerError, "GenerateImage")
			return
		}

		w.WriteHeader(http.StatusOK)

	} else {
		// file does not exist
		firstLetter, _ := utf8.DecodeRuneInString(name)

		img, err := letteravatar.Draw(75, firstLetter, nil)
		if err != nil {
			fmt.Printf("[ ERROR ] [ GenerateImage ] %v", err)
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, err.Error(), 0), err, http.StatusInternalServerError, "GenerateImage")
			return
		}

		file, err := os.Create("images/" + name + ".png")
		if err != nil {
			fmt.Printf("[ ERROR ] [ GenerateImage ] %v", err)
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, err.Error(), 0), err, http.StatusInternalServerError, "GenerateImage")
			return
		}

		defer file.Close()

		err = png.Encode(file, img)
		if err != nil {
			fmt.Printf("[ ERROR ] [ GenerateImage ] %v", err)
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, err.Error(), 0), err, http.StatusInternalServerError, "GenerateImage")
			return
		}

		//return image
		w.Header().Set("Content-Type", "image/png")
		err = png.Encode(w, img)
		if err != nil {
			fmt.Printf("[ ERROR ] [ GenerateImage ] %v", err)
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, err.Error(), 0), err, http.StatusInternalServerError, "GenerateImage")
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func CheckEmail(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	email := mux.Vars(r)["email"]

	var user models.Users
	db.Where("email = ?", email).First(&user)

	if user.ID == 0 {
		models.RespHandler(w, r, models.DefaultPosResponse("email is available"), nil, http.StatusOK, "CheckEmail")
	} else {
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "email is already in use", 0), nil, http.StatusUnauthorized, "CheckEmail")
	}
}

func CheckElsysEmail(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	email := mux.Vars(r)["email"]

	//check if email is valid elsys email format [name].[middle name initial].[surname].[yearofEntry]@elsys-bg.org
	//split email by . and @
	splitEmail := strings.Split(email, ".")
	splitEmail2 := strings.Split(splitEmail[3], "@")

	//check if splitEmail2[0] is a year after 2018
	year, err := strconv.Atoi(splitEmail2[0])
	if err != nil {
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "email is not valid elsys email", 0), nil, http.StatusUnauthorized, "CheckElsysEmail")
		return
	}

	if year < 2018 {
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "email is not authorized to participate", 0), nil, http.StatusUnauthorized, "CheckElsysEmail")
		return
	}

	if splitEmail2[1] != "elsys-bg.org" {
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "email is not valid elsys email", 0), nil, http.StatusUnauthorized, "CheckElsysEmail")
		return
	}

	var user models.Users
	db.Where("elsys_email = ?", email).First(&user)

	if user.ID == 0 {
		models.RespHandler(w, r, models.DefaultPosResponse("email is available"), nil, http.StatusOK, "CheckElsysEmail")
	} else {
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "email is already in use", 0), nil, http.StatusUnauthorized, "CheckElsysEmail")
	}
}

func IsVerified(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	sub, err := ReturnAuthID(w, r, db)

	if err != nil {
		fmt.Printf("[ ERROR ] [ IsVerified ] %v", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, err.Error(), 0), err, http.StatusUnauthorized, "IsVerified")
		return
	}

	var isEmailVerified bool
	//check if email is verified in db security table
	db.Table("security").Joins("JOIN users ON users.security_id = security.id").Where("users.id = ?", sub).Pluck("email_verified", &isEmailVerified)

	if isEmailVerified {
		models.RespHandler(w, r, models.DefaultPosResponse("true"), nil, http.StatusOK, "IsVerified")
	} else {
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "false", 0), nil, http.StatusUnauthorized, "IsVerified")
	}
}
