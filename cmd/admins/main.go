package admins

import (
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

//func SearchWithFilters(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
//	//parse request body
//	var filters models.ParseFilterUsers
//	err := json.NewDecoder(r.Body).Decode(&filters)
//	if err != nil {
//		models.RespHandler(w, r, nil, err, http.StatusBadRequest, "AdminSearchWithFilters")
//		return
//	}
//
//	// get users with filters
//	var users []models.Users
//
//}
