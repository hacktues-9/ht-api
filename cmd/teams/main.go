package teams

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hacktues-9/API/cmd/users"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"strconv"

	"github.com/hacktues-9/API/pkg/models"
)

func CreateTeam(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	team := models.Team{}
	parseTeam := models.ParseTeam{}
	user := models.Users{}

	sub, err := users.ReturnAuthID(r)
	if err != nil {
		fmt.Printf("[ ERROR ] [ CreateTeam ] %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, err.Error(), 0), err, http.StatusUnauthorized, "CreateTeam")
		return
	}

	db.Where("id = ?\n", sub).First(&user)

	err = json.NewDecoder(r.Body).Decode(&parseTeam)
	if err != nil {
		fmt.Printf("[ ERROR ] [ CreateTeam ] parse: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "parse: "+err.Error(), 0), err, http.StatusInternalServerError, "CreateTeam")
		return
	}

	if user.TeamID != 0 {
		fmt.Printf("[ ERROR ] [ CreateTeam ] user already has a team\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusForbidden, "user already has a team", 0), err, http.StatusForbidden, "CreateTeam")
		return
	}

	team = models.Team{
		Name:        parseTeam.Name,
		Description: parseTeam.Description,
		Logo:        "https://api.hacktues.bg/api/image/" + parseTeam.Name,
	}

	if result := db.Omit("ProjectID", "InvitesID").Create(&team); result.Error != nil {
		fmt.Printf("[ ERROR ] [ CreateTeam ] create: %v\n", result.Error)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "create: "+result.Error.Error(), 0), result.Error, http.StatusInternalServerError, "CreateTeam")
		return
	}

	var technologies []uint

	for _, tech := range parseTeam.Technologies {
		var tempTech models.Technologies
		db.Where("technology = ?\n", tech).First(&tempTech)
		technologies = append(technologies, tempTech.ID)
	}

	for _, tech := range technologies {
		db.Create(&models.TeamTechnologies{TeamID: team.ID, TechnologiesID: tech})
	}

	//modify user as captain of team

	user.TeamID = team.ID
	user.RoleID = 2

	if result := db.Save(&user); result.Error != nil {
		fmt.Printf("[ ERROR ] [ CreateTeam ] save: %v\n", result.Error)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "save: "+result.Error.Error(), 0), result.Error, http.StatusInternalServerError, "CreateTeam")
		return
	}

	//send invites to invitees
	for _, invitee := range parseTeam.Invitees {
		var tempUser models.Users
		db.Where("id = ?\n", invitee.ID).First(&tempUser)
		if tempUser.ID == 0 {
			fmt.Printf("[ ERROR ] [ CreateTeam ] user not found")
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusNotFound, "user not found", 0), err, http.StatusNotFound, "CreateTeam")
			return
		}

		if tempUser.TeamID != 0 {
			fmt.Printf("[ ERROR ] [ CreateTeam ] user already has a team\n")
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusForbidden, "user already has a team", 0), err, http.StatusForbidden, "CreateTeam")
			return
		}

		invite := models.Invite{
			UserID:      tempUser.ID,
			TeamID:      team.ID,
			Pending:     true,
			Application: false,
		}

		if result := db.Create(&invite); result.Error != nil {
			fmt.Printf("[ ERROR ] [ CreateTeam ] create: %v\n", result.Error)
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "create: "+result.Error.Error(), 0), result.Error, http.StatusInternalServerError, "CreateTeam")
			return
		}

	}

	models.RespHandler(w, r, models.DefaultPosResponse(team.ID), nil, http.StatusCreated, "CreateTeam")
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

// 	db.Where("id = ?\n", sub).First(&team)

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
// 		db.Where("technology = ?\n", tech).First(&tempTech)
// 		technologies = append(technologies, tempTech.ID)
// 	}

// 	for _, tech := range technologies {
// 		db.Create(&models.ProjectTechnologies{ProjectID: project.ID, TechnologiesID: tech})
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("Project created successfully"))
// }

// InviteUserToTeam func
func InviteUserToTeam(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	parseInvite := models.ParseInvite{}
	err := json.NewDecoder(r.Body).Decode(&parseInvite)
	if err != nil {
		fmt.Printf("[ ERROR ] [ InviteUserToTeam ] parse: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "parse: "+err.Error(), 0), err, http.StatusBadRequest, "InviteUserToTeam")
		return
	}

	sub, err := users.ReturnAuthID(r)

	if err != nil {
		fmt.Printf("[ ERROR ] [ InviteUserToTeam ] %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, err.Error(), 0), err, http.StatusUnauthorized, "InviteUserToTeam")
	}

	captain := models.Users{}
	db.Where("id = ?\n", sub).First(&captain)

	// Check if captain is team owner
	if captain.RoleID != 2 {
		fmt.Printf("[ ERROR ] [ InviteUserToTeam ] user not team owner: %v\n", captain.RoleID)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "user not team owner", 0), err, http.StatusUnauthorized, "InviteUserToTeam")
		return
	}

	team := models.Team{}
	db.Where("id = ?\n", captain.TeamID).First(&team)

	if team.ID == 0 {
		fmt.Printf("[ ERROR ] [ InviteUserToTeam ] team not match\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusNotFound, "team not match", 0), err, http.StatusNotFound, "InviteUserToTeam")
		return
	}

	if team.ID != parseInvite.TeamID {
		fmt.Printf("[ ERROR ] [ InviteUserToTeam ] team not found\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusNotFound, "team not found", 0), err, http.StatusNotFound, "InviteUserToTeam")
		return
	}

	// Check if user exists
	user := models.Users{}
	db.Where("id = ?\n", parseInvite.UserID).First(&user)

	if user.ID == 0 {
		fmt.Printf("[ ERROR ] [ InviteUserToTeam ] user not found\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusNotFound, "user not found", 0), err, http.StatusNotFound, "InviteUserToTeam")
		return
	}

	// Check if user is already in team
	if user.TeamID != 0 {
		fmt.Printf("[ ERROR ] [ InviteUserToTeam ] user already in team\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "user already in team", 0), err, http.StatusBadRequest, "InviteUserToTeam")
		return
	}

	// Check if user is already invited
	invite := models.Invite{}
	db.Where("user_id = ? AND team_id = ?\n", parseInvite.UserID, parseInvite.TeamID).First(&invite)

	if invite.ID != 0 {
		fmt.Printf("[ ERROR ] [ InviteUserToTeam ] user already invited\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "user already invited", 0), err, http.StatusBadRequest, "InviteUserToTeam")
		return
	}

	// Create invite
	invite = models.Invite{
		UserID:      parseInvite.UserID,
		TeamID:      parseInvite.TeamID,
		Pending:     true,
		Application: false,
	}

	if result := db.Create(&invite); result.Error != nil {
		fmt.Printf("[ ERROR ] [ InviteUserToTeam ] create: %v\n", result.Error)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "create: "+result.Error.Error(), 0), err, http.StatusInternalServerError, "InviteUserToTeam")
		return
	}

	models.RespHandler(w, r, models.DefaultPosResponse("success"), nil, http.StatusCreated, "InviteUserToTeam")
}

func ApplyToTeam(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	parseApply := models.ParseApply{}
	err := json.NewDecoder(r.Body).Decode(&parseApply)
	if err != nil {
		fmt.Printf("[ ERROR ] [ ApplyToTeam ] parse: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "parse: "+err.Error(), 0), err, http.StatusBadRequest, "ApplyToTeam")
		return
	}

	sub, err := users.ReturnAuthID(r)
	if err != nil {
		fmt.Printf("[ ERROR ] [ ApplyToTeam ] %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, err.Error(), 0), err, http.StatusBadRequest, "ApplyToTeam")
		return
	}

	user := models.Users{}
	db.Where("id = ?\n", sub).First(&user)

	if user.ID == 0 {
		fmt.Printf("[ ERROR ] [ ApplyToTeam ] user not found\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusNotFound, "user not found", 0), err, http.StatusNotFound, "ApplyToTeam")
		return
	}

	if user.ID != parseApply.UserID {
		fmt.Printf("[ ERROR ] [ ApplyToTeam ] user not match\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusNotFound, "user not match", 0), err, http.StatusNotFound, "ApplyToTeam")
		return
	}

	// Check if user is already in team
	if user.TeamID != 0 {
		fmt.Printf("[ ERROR ] [ ApplyToTeam ] user already in team\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "user already in team", 0), err, http.StatusBadRequest, "ApplyToTeam")
		return
	}

	// Check if user is already invited
	invite := models.Invite{}
	db.Where("user_id = ? AND team_id = ?\n", parseApply.UserID, parseApply.TeamID).First(&invite)

	if invite.ID != 0 {
		fmt.Printf("[ ERROR ] [ ApplyToTeam ] user already invited\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "user already invited", 0), err, http.StatusBadRequest, "ApplyToTeam")
		return
	}

	// Create invite
	invite = models.Invite{
		UserID:      parseApply.UserID,
		TeamID:      parseApply.TeamID,
		Pending:     true,
		Application: true,
	}

	if result := db.Create(&invite); result.Error != nil {
		fmt.Printf("[ ERROR ] [ ApplyToTeam ] create: %v\n", result.Error)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "create: "+result.Error.Error(), 0), err, http.StatusInternalServerError, "ApplyToTeam")
		return
	}

	models.RespHandler(w, r, models.DefaultPosResponse("success"), nil, http.StatusCreated, "ApplyToTeam")
}

func RecommendTeam(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	//get user
	user := models.Users{}

	sub, err := users.ReturnAuthID(r)
	if err != nil {
		fmt.Printf("[ ERROR ] [ RecommendTeam ] %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, err.Error(), 0), err, http.StatusBadRequest, "RecommendTeam")
		return
	}

	db.Where("id = ?\n", sub).First(&user)

	//get user technologies
	var userTechnologies []models.Technologies
	err = db.Model(&user).Association("Technologies").Find(&userTechnologies)
	if err != nil {
		fmt.Printf("[ ERROR ] [ RecommendTeam ] get user technologies: find: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "get user technologies: find: "+err.Error(), 0), err, http.StatusInternalServerError, "RecommendTeam")
		return
	}

	//get teams
	var teams []models.Team
	db.Find(&teams)

	//get teams technologies
	var teamsTechnologies []models.Technologies
	for _, team := range teams {
		err := db.Model(&team).Association("Technologies").Find(&teamsTechnologies)
		if err != nil {
			fmt.Printf("[ ERROR ] [ RecommendTeam ] get teams technologies: find: %v\n", err)
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "get teams technologies: find: "+err.Error(), 0), err, http.StatusInternalServerError, "RecommendTeam")
			return
		}
	}

	//get teams projects
	var teamsProjects []models.Project
	for _, team := range teams {
		err := db.Model(&team).Association("Projects").Find(&teamsProjects)
		if err != nil {
			fmt.Printf("[ ERROR ] [ RecommendTeam ] get teams projects: %v\n", err)
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "get teams projects: "+err.Error(), 0), err, http.StatusInternalServerError, "RecommendTeam")
			return
		}
	}

	//get teams projects technologies
	var teamsProjectsTechnologies []models.Technologies
	for _, project := range teamsProjects {
		err := db.Model(&project).Association("Technologies").Find(&teamsProjectsTechnologies)
		if err != nil {
			fmt.Printf("[ ERROR ] [ RecommendTeam ] get teams projects technologies: find: %v\n", err)
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "get teams projects technologies: find: "+err.Error(), 0), err, http.StatusInternalServerError, "RecommendTeam")
			return
		}
	}

	//compare user technologies with teams technologies
	teamsTechnologiesMap := make(map[uint]uint)
	for _, teamTech := range teamsTechnologies {
		for _, userTech := range userTechnologies {
			if teamTech.ID == userTech.ID {
				teamsTechnologiesMap[teamTech.ID]++
			}
		}
	}

	//compare user technologies with teams projects technologies
	teamsProjectsTechnologiesMap := make(map[uint]uint)
	for _, teamProjectTech := range teamsProjectsTechnologies {
		for _, userTech := range userTechnologies {
			if teamProjectTech.ID == userTech.ID {
				teamsProjectsTechnologiesMap[teamProjectTech.ID]++
			}
		}
	}

	var teamsWithMostCommonTechnologiesAndMostCommonTechnologiesInProjects []models.Team
	for _, team := range teams {
		if teamsTechnologiesMap[team.ID] > 0 && teamsProjectsTechnologiesMap[team.ID] > 0 {
			teamsWithMostCommonTechnologiesAndMostCommonTechnologiesInProjects = append(teamsWithMostCommonTechnologiesAndMostCommonTechnologiesInProjects, team)
		}
	}

	//show teams
	models.RespHandler(w, r, teamsWithMostCommonTechnologiesAndMostCommonTechnologiesInProjects, nil, http.StatusOK, "RecommendTeam")
}

func AcceptInvite(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	vars := mux.Vars(r)
	teamID, err := strconv.Atoi(vars["teamId"])
	if err != nil {
		fmt.Printf("[ ERROR ] [ AcceptInvite ] accept invite: parse teamID: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "accept invite: parse teamID: "+err.Error(), 0), err, http.StatusBadRequest, "AcceptInvite")
		return
	}
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		fmt.Printf("[ ERROR ] [ AcceptInvite ] accept invite: parse teamID: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "accept invite: parse teamID: "+err.Error(), 0), err, http.StatusBadRequest, "AcceptInvite")
		return
	}

	sub, err := users.ReturnAuthID(r)
	if err != nil {
		fmt.Printf("[ ERROR ] [ AcceptInvite ] %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, err.Error(), 0), err, http.StatusUnauthorized, "AcceptInvite")
		return
	}

	//application true => user applied to join a team
	//application false => user was invited to join a team

	// get invite
	invite := models.Invite{}
	db.Where("user_id = ? AND team_id = ?\n", userID, teamID).First(&invite)
	if invite.ID == 0 {
		fmt.Printf("[ ERROR ] [ AcceptInvite ] accept invite: invite not found\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "accept invite: invite not found", 0), err, http.StatusBadRequest, "AcceptInvite")
		return
	}

	if float64(sub) != float64(userID) {
		// check if sub is team leader of teamID
		user := models.Users{}
		db.Where("id = ?\n", sub).First(&user)
		if user.RoleID != 2 {
			fmt.Printf("[ ERROR ] [ AcceptInvite ] accept invite: user is not a team leader\n")
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "accept invite: user is not a team leader", 0), err, http.StatusUnauthorized, "AcceptInvite")
			return
		}

		if float64(user.TeamID) != float64(teamID) {
			fmt.Printf("[ ERROR ] [ AcceptInvite ] accept invite: user is not a team leader of team\n")
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "accept invite: user is not a team leader of team", 0), err, http.StatusUnauthorized, "AcceptInvite")
			return
		}

		if !invite.Application {
			fmt.Printf("[ ERROR ] [ AcceptInvite ] accept invite: can not accept user you have invited\n")
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "accept invite: can not accept user you have invited", 0), err, http.StatusUnauthorized, "AcceptInvite")
			return
		}

		//accept user to team
		db.Model(&models.Users{}).Where("id = ?", userID).Update("team_id", teamID)
		//delete invite
		db.Where("user_id = ? AND team_id = ?\n", userID, teamID).Delete(&models.Invite{})
		models.RespHandler(w, r, models.DefaultPosResponse("success"), nil, http.StatusOK, "AcceptInvite")
		return
	} else {
		if invite.Application {
			fmt.Printf("[ ERROR ] [ AcceptInvite ] accept invite: can not accept team you have applied to join\n")
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "accept invite: can not accept team you have applied to join", 0), err, http.StatusUnauthorized, "AcceptInvite")
			return
		}

		//check if team is full (max 5 members)
		var num int64
		db.Model(&models.Users{}).Where("team_id = ?", teamID).Count(&num)

		if num >= 5 {
			fmt.Printf("[ ERROR ] [ AcceptInvite ] accept invite: team is full\n")
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "accept invite: team is full", 0), err, http.StatusUnauthorized, "AcceptInvite")
			return
		}

		//accept user to team
		db.Model(&models.Users{}).Where("id = ?", userID).Update("team_id", teamID)
		//delete invite
		db.Where("user_id = ? AND team_id = ?\n", userID, teamID).Delete(&models.Invite{})
		models.RespHandler(w, r, models.DefaultPosResponse("success"), nil, http.StatusOK, "AcceptInvite")
		return
	}
}

func DeclineInvite(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	vars := mux.Vars(r)
	teamID, err := strconv.Atoi(vars["teamId"])
	if err != nil {
		fmt.Printf("[ ERROR ] [ DeclineInvite ] decline invite: parse teamID: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "decline invite: parse teamID: "+err.Error(), 0), err, http.StatusBadRequest, "DeclineInvite")
		return
	}
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		fmt.Printf("[ ERROR ] [ DeclineInvite ] decline invite: parse teamID: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "decline invite: parse teamID: "+err.Error(), 0), err, http.StatusBadRequest, "DeclineInvite")
		return
	}

	sub, err := users.ReturnAuthID(r)
	if err != nil {
		fmt.Printf("[ ERROR ] [ DeclineInvite ] %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, err.Error(), 0), err, http.StatusUnauthorized, "DeclineInvite")
		return
	}

	//application true => user applied to join a team
	//application false => user was invited to join a team

	// get invite
	invite := models.Invite{}
	db.Where("user_id = ? AND team_id = ?\n", userID, teamID).First(&invite)
	if invite.ID == 0 {
		fmt.Printf("[ ERROR ] [ DeclineInvite ] decline invite: invite not found\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "decline invite: invite not found", 0), err, http.StatusBadRequest, "DeclineInvite")
		return
	}

	if float64(sub) != float64(userID) {
		// check if sub is team leader of teamID
		user := models.Users{}
		db.Where("id = ?\n", sub).First(&user)
		if user.RoleID != 2 {
			fmt.Printf("[ ERROR ] [ DeclineInvite ] decline invite: user is not a team leader\n")
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "decline invite: user is not a team leader", 0), err, http.StatusUnauthorized, "DeclineInvite")
			return
		}

		if float64(user.TeamID) != float64(teamID) {
			fmt.Printf("[ ERROR ] [ DeclineInvite ] decline invite: user is not a team leader of team\n")
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "decline invite: user is not a team leader of team", 0), err, http.StatusUnauthorized, "DeclineInvite")
			return
		}

		if !invite.Application {
			fmt.Printf("[ ERROR ] [ DeclineInvite ] decline invite: can not decline user you have invited\n")
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "decline invite: can not decline user you have invited", 0), err, http.StatusUnauthorized, "DeclineInvite")
			return
		}

		//decline user from team
		db.Model(&models.Users{}).Where("id = ?", userID).Update("team_id", nil)
		//delete invite
		db.Where("user_id = ? AND team_id = ?\n", userID, teamID).Delete(&models.Invite{})
		models.RespHandler(w, r, models.DefaultPosResponse("success"), nil, http.StatusOK, "DeclineInvite")
	} else {
		if invite.Application {
			fmt.Printf("[ ERROR ] [ DeclineInvite ] decline invite: can not decline team you have applied to join\n")
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "decline invite: can not decline team you have applied to join", 0), err, http.StatusUnauthorized, "DeclineInvite")
			return
		}

		//decline user from team
		db.Model(&models.Users{}).Where("id = ?", userID).Update("team_id", nil)
		//delete invite
		db.Where("user_id = ? AND team_id = ?\n", userID, teamID).Delete(&models.Invite{})
		models.RespHandler(w, r, models.DefaultPosResponse("success"), nil, http.StatusOK, "DeclineInvite")
		return
	}

}

func GetTeams(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var parseTeams []models.ParseTeamView
	db.Raw("SELECT * FROM teams").Scan(&parseTeams)

	// parse teams to teams
	var teams []models.TeamsView

	for _, parseTeam := range parseTeams {

		teams = append(teams, models.TeamsView{
			ID:           parseTeam.ID,
			Name:         parseTeam.Name,
			Logo:         parseTeam.Logo,
			Members:      []models.MemberView{},
			Project:      models.ProjectView{},
			Technologies: []string{},
			IsVerified:   parseTeam.Approved,
		})

		var members []models.Users //get team members with info table, socials, class,discord, github, role
		db.Preload("Info.Socials.Discord").Preload("Info.Class").Preload("Info.Socials.Github").Table("users").Preload(clause.Associations).Where("team_id = ?\n", parseTeam.ID).Find(&members)

		//add the member to the team
		for _, member := range members {
			// discord = member.Info.Socials.Discord.Username + "#" + member.Info.Socials.Discord.Discriminator, if member.Info.Socials.Discord.Username == "" { discord = "" }
			discord := string(member.Info.Socials.Discord.Username + "#" + member.Info.Socials.Discord.Discriminator)
			if discord == "#" {
				discord = ""
			}
			teams[len(teams)-1].Members = append(teams[len(teams)-1].Members, models.MemberView{
				ID:             member.ID,
				Name:           member.FirstName + " " + member.LastName,
				ProfilePicture: member.Info.Socials.ProfilePicture,
				Role:           member.Role.Name,
				Class:          strconv.Itoa(member.Info.Grade) + " " + member.Info.Class.Name,
				Email:          member.Email,
				Github:         member.Info.Socials.Github.Login,
				Discord:        discord,
			})
		}

		//get team project
		if parseTeam.ProjectID != 0 {
			var teamProject models.Project
			db.Table("projects").Where("id = ?\n", parseTeam.ProjectID).First(&teamProject)

			//get team project technologies
			var teamProjectTechnologies []models.Technologies
			db.Table("technologies").Joins("JOIN project_technologies ON project_technologies.project_id = ?\n", parseTeam.ProjectID).Where("project_technologies.technology_id = technologies.id").Find(&teamProjectTechnologies)

			//parse team project technologies
			var teamProjectTechnologiesParsed []string
			for _, teamProjectTechnology := range teamProjectTechnologies {
				teamProjectTechnologiesParsed = append(teamProjectTechnologiesParsed, teamProjectTechnology.Technology)
			}

			//add team project to team
			teams[len(teams)-1].Project = models.ProjectView{
				ID:           teamProject.ID,
				Name:         teamProject.Name,
				Description:  teamProject.Description,
				Technologies: teamProjectTechnologiesParsed,
			}

		}

		//get team technologies
		var teamTechnologies []string
		db.Table("technologies").Joins("JOIN team_technologies ON team_technologies.team_id = ?\n", parseTeam.ID).Where("team_technologies.technologies_id = technologies.id").Pluck("technologies.technology", &teamTechnologies)

		//add team technologies to team
		teams[len(teams)-1].Technologies = teamTechnologies
	}

	models.RespHandler(w, r, models.DefaultPosResponse(teams), nil, http.StatusOK, "GetTeams")
}

func SearchInvitees(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// return first three searchView that match the search query

	// get query search=...
	query := r.URL.Query().Get("search")

	sub, err := users.ReturnAuthID(r)
	if err != nil {
		fmt.Printf("[ ERROR ] [ SearchInvitees ] %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, err.Error(), 0), err, http.StatusUnauthorized, "SearchInvitees")
		return
	}

	// get user teamId from db
	var teamId uint
	db.Table("users").Where("id = ?", sub).Select("team_id").Scan(&teamId)

	// get searchView from db
	var searchView []models.SearchView
	// use searchuser function to get searchView from db
	db.Raw("SELECT * FROM searchuser(?, ?, ?)", query, teamId, sub).Scan(&searchView)

	// return searchView
	models.RespHandler(w, r, models.DefaultPosResponse(searchView), nil, http.StatusOK, "SearchInvitees")
}

func GetTeamID(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// get team id from user id

	// get user id from url
	userID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Printf("[ ERROR ] [ GetTeamID ] parse: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "parse: "+err.Error(), 0), err, http.StatusInternalServerError, "GetTeamID")
		return
	}

	// get team id from db
	var teamID int
	err = db.Table("users").Where("id = ?\n", userID).Select("team_id").Row().Scan(&teamID)
	if err != nil {
		fmt.Printf("[ ERROR ] [ GetTeamID ] select: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "select: "+err.Error(), 0), err, http.StatusInternalServerError, "GetTeamID")
		return
	}

	// return team id
	models.RespHandler(w, r, models.DefaultPosResponse(teamID), nil, http.StatusOK, "GetTeamID")
}

func GetCaptainID(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// get captain id from team id

	// get team id from url
	teamID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Printf("[ ERROR ] [ GetCaptainID ] parse: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "parse: "+err.Error(), 0), err, http.StatusInternalServerError, "GetCaptainID")
		return
	}

	// get captain id from db
	var captainID int
	err = db.Table("users").Where("team_id = ? AND role_id = 2", teamID).Select("id").Row().Scan(&captainID)
	if err != nil {
		fmt.Printf("[ ERROR ] [ GetCaptainID ] select: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "select: "+err.Error(), 0), err, http.StatusInternalServerError, "GetCaptainID")
		return
	}

	// return captain id
	models.RespHandler(w, r, models.DefaultPosResponse(captainID), nil, http.StatusOK, "GetCaptainID")
}

func GetTeam(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	teamID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Printf("[ ERROR ] [ GetTeam ] parse: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "parse: "+err.Error(), 0), err, http.StatusInternalServerError, "GetTeam")
		return
	}

	// get team from db
	var team models.GetTeamView
	db.Table("team").Select("team.name, team.description, team.logo").Where("team.id = ?\n", teamID).Scan(&team)
	db.Table("users").Select("users.id, concat(users.first_name, ' ', users.last_name) AS name, socials.profile_picture AS avatar, role.name AS role").Joins("JOIN info ON users.info_id = info.id").Joins("JOIN socials ON info.socials_id = socials.id").Joins("JOIN role ON role.id = users.role_id").Where("users.team_id = ?\n", teamID).Scan(&team.Members)

	// get team technologies from db
	var teamTechnologies []string
	db.Table("team_technologies").Select("technologies.technology").Joins("JOIN technologies ON team_technologies.technologies_id = technologies.id").Where("team_technologies.team_id = ?\n", teamID).Scan(&teamTechnologies)
	team.Technologies = teamTechnologies

	// get team projects from db
	var teamProject models.ProjectTeamView
	team.Project = teamProject

	// return team
	models.RespHandler(w, r, models.DefaultPosResponse(team), nil, http.StatusOK, "GetTeam")
}

func GetInvitees(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// get invitees from team id

	// get team id from url
	teamID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Printf("[ ERROR ] [ GetInvitees ] parse: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "parse: "+err.Error(), 0), err, http.StatusInternalServerError, "GetInvitees")
		return
	}

	// get invitees from db
	var invitees []models.SearchView //(id bigint, name text, profile_picture text, isinvited boolean)
	db.Table("users").Select("users.id, concat(users.first_name, ' ', users.last_name) AS name, socials.profile_picture AS profile_picture, true AS isinvited").Joins("JOIN info ON users.info_id = info.id").Joins("JOIN socials ON info.socials_id = socials.id").Joins("LEFT JOIN invite ON invite.user_id = u.id AND invite.team_id = teamid").Where("users.team_id = ? AND invite.id IS NULL", teamID).Scan(&invitees)

	// return invitees
	models.RespHandler(w, r, models.DefaultPosResponse(invitees), nil, http.StatusOK, "GetInvitees")
}

func KickUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	sub, err := users.ReturnAuthID(r)
	if err != nil {
		fmt.Printf("[ ERROR ] [ KickUser ] return auth id: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "return auth id: "+err.Error(), 0), err, http.StatusInternalServerError, "KickUser")
		return
	}

	// get user id from url
	userID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Printf("[ ERROR ] [ KickUser ] parse: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "parse: "+err.Error(), 0), err, http.StatusInternalServerError, "KickUser")
		return
	}

	// get team id from user
	var teamID int
	err = db.Table("users").Select("team_id").Where("id = ?", userID).Row().Scan(&teamID)

	// check if user is captain
	var captainID int
	err = db.Table("users").Where("team_id = ? AND role_id = 2", teamID).Select("id").Row().Scan(&captainID)
	if err != nil {
		fmt.Printf("[ ERROR ] [ KickUser ] select: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "select: "+err.Error(), 0), err, http.StatusInternalServerError, "KickUser")
		return
	}

	if float64(sub) != float64(captainID) {
		fmt.Printf("[ ERROR ] [ KickUser ] user is not captain: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "user is not captain: "+err.Error(), 0), err, http.StatusInternalServerError, "KickUser")
		return
	}

	// kick user
	err = db.Table("users").Where("id = ?", userID).Update("team_id", nil).Error
	if err != nil {
		fmt.Printf("[ ERROR ] [ KickUser ] update: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "update: "+err.Error(), 0), err, http.StatusInternalServerError, "KickUser")
		return
	}

	// return success
	models.RespHandler(w, r, models.DefaultPosResponse("success"), nil, http.StatusOK, "KickUser")
}

func UpdateTeam(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	sub, err := users.ReturnAuthID(r)
	if err != nil {
		fmt.Printf("[ ERROR ] [ UpdateTeam ] return auth id: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "return auth id: "+err.Error(), 0), err, http.StatusInternalServerError, "UpdateTeam")
		return
	}

	// get team id from url
	teamID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Printf("[ ERROR ] [ UpdateTeam ] parse: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "parse: "+err.Error(), 0), err, http.StatusInternalServerError, "UpdateTeam")
		return
	}

	// check if user is captain
	var captainID int
	err = db.Table("users").Where("team_id = ? AND role_id = 2", teamID).Select("id").Row().Scan(&captainID)
	if err != nil {
		fmt.Printf("[ ERROR ] [ UpdateTeam ] select: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "select: "+err.Error(), 0), err, http.StatusInternalServerError, "UpdateTeam")
		return
	}

	if float64(sub) != float64(captainID) {
		fmt.Printf("[ ERROR ] [ UpdateTeam ] user is not captain: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "user is not captain: "+err.Error(), 0), err, http.StatusInternalServerError, "UpdateTeam")
		return
	}

	var team models.GetTeamView
	err = json.NewDecoder(r.Body).Decode(&team)
	if err != nil {
		fmt.Printf("[ ERROR ] [ UpdateTeam ] decode: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "decode: "+err.Error(), 0), err, http.StatusInternalServerError, "UpdateTeam")
		return
	}

	// update team
	err = db.Table("team").Where("id = ?", teamID).Updates(map[string]interface{}{"name": team.Name, "description": team.Description, "logo": "https://api.hacktues.bg/api/image/" + team.Name}).Error
	if err != nil {
		fmt.Printf("[ ERROR ] [ UpdateTeam ] update: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "update: "+err.Error(), 0), err, http.StatusInternalServerError, "UpdateTeam")
		return
	}

	// update technologies
	db.Where("team_id = ?", teamID).Delete(models.TeamTechnologies{})

	for _, tech := range team.Technologies {
		//get technology id
		var techID int
		db.Table("technologies").Select("id").Where("technology = ?", tech).Row().Scan(&techID)
		teamTech := models.TeamTechnologies{
			TeamID:         uint(teamID),
			TechnologiesID: uint(techID),
		}
		db.Table("team_technologies").Create(&teamTech)
	}
	// return success
	models.RespHandler(w, r, models.DefaultPosResponse("success"), nil, http.StatusOK, "UpdateTeam")
}

func LeaveTeam(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	sub, err := users.ReturnAuthID(r)
	if err != nil {
		fmt.Printf("[ ERROR ] [ LeaveTeam ] return auth id: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "return auth id: "+err.Error(), 0), err, http.StatusInternalServerError, "LeaveTeam")
		return
	}

	// get team id from url
	teamID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Printf("[ ERROR ] [ LeaveTeam ] parse: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "parse: "+err.Error(), 0), err, http.StatusInternalServerError, "LeaveTeam")
		return
	}

	// check if user is captain
	var captainID int
	err = db.Table("users").Where("team_id = ? AND role_id = 2", teamID).Select("id").Row().Scan(&captainID)
	if err != nil {
		fmt.Printf("[ ERROR ] [ LeaveTeam ] select: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "select: "+err.Error(), 0), err, http.StatusInternalServerError, "LeaveTeam")
		return
	}

	if float64(sub) == float64(captainID) {
		// delete team
		//clear all members from team
		err = db.Table("users").Where("team_id = ?", teamID).Update("team_id", nil).Error
		if err != nil {
			fmt.Printf("[ ERROR ] [ LeaveTeam ] update: %v\n", err)
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "update: "+err.Error(), 0), err, http.StatusInternalServerError, "LeaveTeam")
			return
		}

		//delete team
		err = db.Table("team").Where("id = ?", teamID).Delete(models.Team{}).Error
		if err != nil {
			fmt.Printf("[ ERROR ] [ LeaveTeam ] delete: %v\n", err)
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "delete: "+err.Error(), 0), err, http.StatusInternalServerError, "LeaveTeam")
			return
		}

		// return success
		models.RespHandler(w, r, models.DefaultPosResponse("success"), nil, http.StatusOK, "LeaveTeam")
		return
	}

	// leave team
	err = db.Table("users").Where("id = ?", sub).Updates(map[string]interface{}{"team_id": nil, "role_id": 1}).Error
	if err != nil {
		fmt.Printf("[ ERROR ] [ LeaveTeam ] update: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "update: "+err.Error(), 0), err, http.StatusInternalServerError, "LeaveTeam")
		return
	}

	// delete team if no members
	var count int64
	db.Table("users").Where("team_id = ?", teamID).Count(&count)

	if count == 0 {
		err = db.Table("team").Where("id = ?", teamID).Delete(models.Team{}).Error
		if err != nil {
			fmt.Printf("[ ERROR ] [ LeaveTeam ] delete: %v\n", err)
			models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "delete: "+err.Error(), 0), err, http.StatusInternalServerError, "LeaveTeam")
			return
		}
	}

	// return success
	models.RespHandler(w, r, models.DefaultPosResponse("success"), nil, http.StatusOK, "LeaveTeam")
}

func DeleteTeam(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	sub, err := users.ReturnAuthID(r)
	if err != nil {
		fmt.Printf("[ ERROR ] [ DeleteTeam ] return auth id: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "return auth id: "+err.Error(), 0), err, http.StatusInternalServerError, "DeleteTeam")
		return
	}

	// get team id from url
	teamID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Printf("[ ERROR ] [ DeleteTeam ] parse: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "parse: "+err.Error(), 0), err, http.StatusInternalServerError, "DeleteTeam")
		return
	}

	// check if user is captain
	var captainID int
	err = db.Table("users").Where("team_id = ? AND role_id = 2", teamID).Select("id").Row().Scan(&captainID)
	if err != nil {
		fmt.Printf("[ ERROR ] [ DeleteTeam ] select: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "select: "+err.Error(), 0), err, http.StatusInternalServerError, "DeleteTeam")
		return
	}

	if float64(sub) != float64(captainID) {
		fmt.Printf("[ ERROR ] [ DeleteTeam ] user is not captain: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "user is not captain: "+err.Error(), 0), err, http.StatusInternalServerError, "DeleteTeam")
		return
	}

	//clear all users from team
	err = db.Table("users").Where("team_id = ?", teamID).Update("team_id", nil).Error
	if err != nil {
		fmt.Printf("[ ERROR ] [ DeleteTeam ] update: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "update: "+err.Error(), 0), err, http.StatusInternalServerError, "DeleteTeam")
		return
	}

	// delete team
	err = db.Table("team").Where("id = ?", teamID).Delete(models.Team{}).Error
	if err != nil {
		fmt.Printf("[ ERROR ] [ DeleteTeam ] delete: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "delete: "+err.Error(), 0), err, http.StatusInternalServerError, "DeleteTeam")
		return
	}

	// return success
	models.RespHandler(w, r, models.DefaultPosResponse("success"), nil, http.StatusOK, "DeleteTeam")
}

func UpdateCaptain(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	sub, err := users.ReturnAuthID(r)
	if err != nil {
		fmt.Printf("[ ERROR ] [ UpdateCaptain ] return auth id: %v\n ", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "return auth id: "+err.Error(), 0), err, http.StatusInternalServerError, "UpdateCaptain")
		return
	}

	// get team id from user
	var teamID int
	db.Table("users").Where("id = ?", sub).Select("team_id").Row().Scan(&teamID)

	// check if user is captain
	var captainID int
	db.Table("users").Where("team_id = ? AND role_id = 2", teamID).Select("id").Row().Scan(&captainID)

	if float64(sub) != float64(captainID) {
		fmt.Printf("[ ERROR ] [ UpdateCaptain ] user is not captain: %v\n ", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "user is not captain: "+err.Error(), 0), err, http.StatusInternalServerError, "UpdateCaptain")
		return
	}

	// get new captain id from url
	newCaptainID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Printf("[ ERROR ] [ UpdateCaptain ] parse: %v\n ", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "parse: "+err.Error(), 0), err, http.StatusInternalServerError, "UpdateCaptain")
		return
	}

	// check if new captain is in team
	var newCaptainTeamID int
	db.Table("users").Where("id = ?", newCaptainID).Select("team_id").Row().Scan(&newCaptainTeamID)

	if float64(newCaptainTeamID) != float64(teamID) {
		fmt.Printf("[ ERROR ] [ UpdateCaptain ] user is not in team: %v\n ", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusInternalServerError, "user is not in team: "+err.Error(), 0), err, http.StatusInternalServerError, "UpdateCaptain")
		return
	}

	// update new captain
	db.Table("users").Where("id = ?", newCaptainID).Update("role_id", 2)

	// update old captain
	db.Table("users").Where("id = ?", captainID).Update("role_id", 1)

	// return success
	models.RespHandler(w, r, models.DefaultPosResponse("success"), nil, http.StatusOK, "UpdateCaptain")
}
