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
	"net/http"
	"os"
	"strings"
	"unicode/utf8"
)

func ReturnAuthID(r *http.Request) (uint, error) {
	cookie, err := r.Cookie("access_token")
	authorizationHeader := r.Header.Get("Authorization")
	fields := strings.Fields(authorizationHeader)
	accessToken := ""
	var sub uint

	if len(fields) != 0 && fields[0] == "Bearer" {
		accessToken = fields[1]
	} else if err == nil {
		accessToken = cookie.Value
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
	userId, err := ReturnAuthID(r)
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
