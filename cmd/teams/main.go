package teams

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/hacktues-9/API/pkg/jwt"
	"github.com/hacktues-9/API/pkg/models"
)

var (
	accessTokenTTL         = time.Hour * 24
	refreshTokenTTL        = time.Hour * 24 * 7
	accessTokenPrivateKey  = os.Getenv("ACCESS_TOKEN_PRIVATE_KEY")
	accessTokenPublicKey   = os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")
	refreshTokenPrivateKey = os.Getenv("REFRESH_TOKEN_PRIVATE_KEY")
	refreshTokenPublicKey  = os.Getenv("REFRESH_TOKEN_PUBLIC_KEY")
)

func CreateTeam(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	team := models.Team{}
	parseTeam := models.ParseTeam{}

	err := json.NewDecoder(r.Body).Decode(&parseTeam)
	if err != nil {
		fmt.Println("createTeam: parse:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("createTeam: parse: " + err.Error()))
		return
	}

	team = models.Team{
		Name:        parseTeam.Name,
		Description: parseTeam.Description,
		Logo:        parseTeam.Logo,
		Color:       parseTeam.Color,
	}

	if result := db.Omit("ProjectID", "InvitesID").Create(&team); result.Error != nil {
		fmt.Println("createTeam: create:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("createTeam: create: " + result.Error.Error()))
		return
	}

	technologies := []uint{}

	for _, tech := range parseTeam.Technologies {
		var tempTech models.Technologies
		db.Where("technology = ?", tech).First(&tempTech)
		technologies = append(technologies, tempTech.ID)
	}

	for _, tech := range technologies {
		db.Create(&models.TeamTechnologies{TeamID: team.ID, TechnologiesID: tech})
	}

	//modify user as captain of team
	user := models.User{}

	cookie, err := r.Cookie("access_token")
	authorizationHeader := r.Header.Get("Authorization")
	fields := strings.Fields(authorizationHeader)
	accessToken := ""

	if len(fields) != 0 && fields[0] == "Bearer" {
		accessToken = fields[1]
	} else if err == nil {
		accessToken = cookie.Value
	} else {
		fmt.Println("get user: access token: get:", err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("get user: access token: get: " + err.Error()))
		return
	}

	sub, err := jwt.ValidateToken(accessToken, accessTokenPublicKey)
	if err != nil {
		fmt.Println("get user: access token: validate:", err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("get user: access token: validate: " + err.Error()))
		return
	}

	db.Where("id = ?", sub).First(&user)

	user.TeamID = team.ID
	user.RoleID = 3

	if result := db.Save(&user); result.Error != nil {
		fmt.Println("createTeam: save:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("createTeam: save: " + result.Error.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Team created successfully"))
}

// func CreateProject(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
// 	project := models.Project{}
// 	parseProject := models.ParseProject{}

// 	err := json.NewDecoder(r.Body).Decode(&parseProject)
// 	if err != nil {
// 		fmt.Println("createProject: parse:", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte("createProject: parse: " + err.Error()))
// 		return
// 	}

// 	team := models.Team{}

// 	cookie, err := r.Cookie("access_token")
// 	authorizationHeader := r.Header.Get("Authorization")
// 	fields := strings.Fields(authorizationHeader)
// 	accessToken := ""

// 	if len(fields) != 0 && fields[0] == "Bearer" {
// 		accessToken = fields[1]
// 	} else if err == nil {
// 		accessToken = cookie.Value
// 	} else {
// 		fmt.Println("get user: access token: get:", err)
// 		w.WriteHeader(http.StatusUnauthorized)
// 		w.Write([]byte("get user: access token: get: " + err.Error()))
// 		return
// 	}

// 	sub, err := jwt.ValidateToken(accessToken, accessTokenPublicKey)
// 	if err != nil {
// 		fmt.Println("get user: access token: validate:", err)
// 		w.WriteHeader(http.StatusUnauthorized)
// 		w.Write([]byte("get user: access token: validate: " + err.Error()))
// 		return
// 	}

// 	db.Where("id = ?", sub).First(&team)

// 	project = models.Project{
// 		Name:        parseProject.Name,
// 		Description: parseProject.Description,
// 		TeamID:      team.ID,
// 	}

// 	if result := db.Omit("TeamID").Create(&project); result.Error != nil {
// 		fmt.Println("createProject: create:", result.Error)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("createProject: create: " + result.Error.Error()))
// 		return
// 	}

// 	technologies := []uint{}

// 	for _, tech := range parseProject.Technologies {
// 		var tempTech models.Technologies
// 		db.Where("technology = ?", tech).First(&tempTech)
// 		technologies = append(technologies, tempTech.ID)
// 	}

// 	for _, tech := range technologies {
// 		db.Create(&models.ProjectTechnologies{ProjectID: project.ID, TechnologiesID: tech})
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("Project created successfully"))
// }
