package teams

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hacktues-9/API/cmd/users"
	"github.com/hacktues-9/API/pkg/models"
	"gorm.io/gorm"
	"net/http"
)

// var (
// 	accessTokenTTL         = time.Hour * 24
// 	refreshTokenTTL        = time.Hour * 24 * 7
// 	accessTokenPrivateKey  = os.Getenv("ACCESS_TOKEN_PRIVATE_KEY")
// 	accessTokenPublicKey   = os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")
// 	refreshTokenPrivateKey = os.Getenv("REFRESH_TOKEN_PRIVATE_KEY")
// 	refreshTokenPublicKey  = os.Getenv("REFRESH_TOKEN_PUBLIC_KEY")
// )

func AddTeamMember(user *models.Users, team *models.Team, db *gorm.DB) error {
	//check if team is full max is 5 members
	return nil
}

func IsUserInTeam(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	//check if user is in team
	sub, err := users.ReturnAuthID(w, r, db)
	if err != nil {
		fmt.Printf("[ ERROR ] [ IsUserInTeam ] return auth id: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "return auth id: "+err.Error(), 0), err, http.StatusInternalServerError, "IsUserInTeam")
		return
	}

	id := mux.Vars(r)["id"]

	//check if user is in team
	var team models.Team
	if err := db.Where("id = ?", id).First(&team).Error; err != nil {
		fmt.Printf("[ ERROR ] [ IsUserInTeam ] get team: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "get team: "+err.Error(), 0), err, http.StatusInternalServerError, "IsUserInTeam")
		return
	}

	var user models.Users
	if err := db.Where("id = ?", sub).First(&user).Error; err != nil {
		fmt.Printf("[ ERROR ] [ IsUserInTeam ] get user: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "get user: "+err.Error(), 0), err, http.StatusInternalServerError, "IsUserInTeam")
		return
	}

	if user.TeamID != team.ID {
		fmt.Printf("[ ERROR ] [ IsUserInTeam ] user not in team: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "user not in team", 0), err, http.StatusUnauthorized, "IsUserInTeam")
		return
	}

	models.RespHandler(w, r, models.DefaultPosResponse("user in team"), nil, http.StatusOK, "IsUserInTeam")
}
