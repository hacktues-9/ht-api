package users

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/hacktues-9/API/pkg/models"
	"github.com/hacktues-9/API/pkg/pass"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Register(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var userInfo models.Info
	var user models.User
	var parseUser models.RegisterUser
	var userSocials models.Socials
	userSecurity := models.Security{}

	err := json.NewDecoder(r.Body).Decode(&parseUser)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	classID, eatingPreferenceID, shirtSizeID, roleID, allergies, technologies := returnDefaultIDs(db, &parseUser)

	password, err := pass.HashPassword(parseUser.Password)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userSocials = models.Socials{
		GithubLink:    "",
		LinkedInLink:  "",
		InstagramLink: "",
	}

	if result := db.Omit("DiscordID").Create(&userSocials); result.Error != nil {
		log.Println(result.Error)
		w.WriteHeader(http.StatusInternalServerError)
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
		log.Println(result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if result := db.Create(&userSecurity); result.Error != nil {
		log.Println(result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user = models.User{
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

	if result := db.Omit("TeamID").Create(&user); result.Error != nil {
		log.Println(result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, allergy := range allergies {
		db.Create(&models.InfoAllergies{InfoID: userInfo.ID, AllergiesID: allergy})
	}

	for _, tech := range technologies {
		db.Create(&models.UserTechnologies{UserID: user.ID, TechnologiesID: tech})
	}

	resp := ParseUser(user.ID, db)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
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
	db.Where("name = ?", "student").First(&role)

	for _, allergy := range user.Allergies {
		var tempAllergy models.Allergies
		db.Where("name = ?", allergy).First(&tempAllergy)
		allergies = append(allergies, tempAllergy.ID)
	}

	for _, tech := range user.Technologies {
		var tempTech models.Technologies
		db.Where("technology = ?", tech.Technology).First(&tempTech)
		technologies = append(technologies, tempTech.ID)
	}

	return class.ID, eatingPreference.ID, shirtSize.ID, role.ID, allergies, technologies
}

func FetchUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	id := r.URL.Query().Get("id")
	resp := ParseUser(id, db)

	json.NewEncoder(w).Encode(resp)
}

func ParseUser(id interface{}, db *gorm.DB) map[string]interface{} {
	user := models.User{}
	userTech := []models.UserTechnologies{}
	userAllergies := []models.InfoAllergies{}
	tech := []models.Technologies{}
	allergies := []models.Allergies{}

	db.Preload(clause.Associations).Preload("Info.Class").Preload("Info.EatingPreference").Preload("Info.Socials").Preload("Info.ShirtSize").Preload("Info.Socials.Discord").Preload("Team.Project").Preload("Team.Invites").Where("ID = ?", id).First(&user)
	db.Where("user_id = ?", id).Find(&userTech)
	db.Where("info_id = ?", user.InfoID).Find(&userAllergies)

	for _, techID := range userTech {
		var tempTech models.Technologies
		db.Where("ID = ?", techID.TechnologiesID).First(&tempTech)
		tech = append(tech, tempTech)
	}

	for _, allergyID := range userAllergies {
		var tempAllergy models.Allergies
		db.Where("ID = ?", allergyID.AllergiesID).First(&tempAllergy)
		allergies = append(allergies, tempAllergy)
	}

	resp := map[string]interface{}{
		"user":      user,
		"tech":      tech,
		"allergies": allergies,
	}

	return resp
}
