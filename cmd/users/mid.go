package users

import (
	"encoding/json"
	"net/http"

	"github.com/hacktues-9/API/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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
		db.Where("technology = ?", tech).First(&tempTech)
		technologies = append(technologies, tempTech.ID)
	}

	return class.ID, eatingPreference.ID, shirtSize.ID, role.ID, allergies, technologies
}
