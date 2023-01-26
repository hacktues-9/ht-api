package users

import (
	"fmt"
	"github.com/hacktues-9/API/pkg/jwt"
	"github.com/hacktues-9/API/pkg/models"
	"gorm.io/gorm"
	"net/http"
	"strings"
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
		return 0, err
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
