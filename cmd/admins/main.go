package admins

import (
	"encoding/json"
	"fmt"
	users2 "github.com/hacktues-9/API/cmd/users"
	"github.com/hacktues-9/API/pkg/models"
	"gorm.io/gorm"
	"net/http"
)

func GetTeams(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// get teams with 3 to 5 members
	var count int64

	db.Raw("SELECT * FROM team WHERE id IN (SELECT team_id FROM users GROUP BY team_id HAVING COUNT(*) >= 3 AND COUNT(*) <= 5)").Count(&count)

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
	if roleId != 2 {
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
