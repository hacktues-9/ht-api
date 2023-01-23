package teams

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"gorm.io/gorm"

	"github.com/hacktues-9/API/pkg/jwt"
	"github.com/hacktues-9/API/pkg/models"
)

var accessTokenPublicKey = os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")

func CreateTeam(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	team := models.Team{}
	parseTeam := models.ParseTeam{}
	user := models.Users{}

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

	err = json.NewDecoder(r.Body).Decode(&parseTeam)
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

	user.TeamID = team.ID
	user.RoleID = 2

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

// InviteUserToTeam func
func InviteUserToTeam(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	parseInvite := models.ParseInvite{}
	err := json.NewDecoder(r.Body).Decode(&parseInvite)
	if err != nil {
		fmt.Println("inviteUserToTeam: parse:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("inviteUserToTeam: parse: " + err.Error()))
		return
	}

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

	captain := models.Users{}
	db.Where("id = ?", sub).First(&captain)

	// Check if captain is team owner
	if captain.RoleID != 2 {
		fmt.Println("inviteUserToTeam: user not team owner")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("inviteUserToTeam: user not team owner"))
		return
	}

	team := models.Team{}
	db.Where("id = ?", captain.TeamID).First(&team)

	if team.ID == 0 {
		fmt.Println("inviteUserToTeam: team not match")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("inviteUserToTeam: team not match"))
		return
	}

	if team.ID != parseInvite.TeamID {
		fmt.Println("inviteUserToTeam: team not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("inviteUserToTeam: team not found"))
		return
	}

	// Check if user exists
	user := models.Users{}
	db.Where("id = ?", parseInvite.UserID).First(&user)

	if user.ID == 0 {
		fmt.Println("inviteUserToTeam: user not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("inviteUserToTeam: user not found"))
		return
	}

	// Check if user is already in team
	if user.TeamID != 0 {
		fmt.Println("inviteUserToTeam: user already in team")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("inviteUserToTeam: user already in team"))
		return
	}

	// Check if user is already invited
	invite := models.Invite{}
	db.Where("user_id = ? AND team_id = ?", parseInvite.UserID, parseInvite.TeamID).First(&invite)

	if invite.ID != 0 {
		fmt.Println("inviteUserToTeam: user already invited")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("inviteUserToTeam: user already invited"))
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
		fmt.Println("inviteUserToTeam: create:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("inviteUserToTeam: create: " + result.Error.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User invited successfully"))
}

func ApplyToTeam(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	parseApply := models.ParseApply{}
	err := json.NewDecoder(r.Body).Decode(&parseApply)
	if err != nil {
		fmt.Println("applyToTeam: parse:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("applyToTeam: parse: " + err.Error()))
		return
	}

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

	user := models.Users{}
	db.Where("id = ?", sub).First(&user)

	if user.ID == 0 {
		fmt.Println("applyToTeam: user not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("applyToTeam: user not found"))
		return
	}

	if user.ID != parseApply.UserID {
		fmt.Println("applyToTeam: user not match")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("applyToTeam: user not match"))
		return
	}

	// Check if user is already in team
	if user.TeamID != 0 {
		fmt.Println("applyToTeam: user already in team")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("applyToTeam: user already in team"))
		return
	}

	// Check if user is already invited
	invite := models.Invite{}
	db.Where("user_id = ? AND team_id = ?", parseApply.UserID, parseApply.TeamID).First(&invite)

	if invite.ID != 0 {
		fmt.Println("applyToTeam: user already invited")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("applyToTeam: user already invited"))
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
		fmt.Println("applyToTeam: create:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("applyToTeam: create: " + result.Error.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User applied successfully"))
}

func RecommendTeam(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	//get user
	user := models.Users{}

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

	//get user technologies
	userTechnologies := []models.Technologies{}
	db.Model(&user).Association("Technologies").Find(&userTechnologies)

	//get teams
	teams := []models.Team{}
	db.Find(&teams)

	//get teams technologies
	teamsTechnologies := []models.Technologies{}
	for _, team := range teams {
		db.Model(&team).Association("Technologies").Find(&teamsTechnologies)
	}

	//get teams projects
	teamsProjects := []models.Project{}
	for _, team := range teams {
		db.Model(&team).Association("Projects").Find(&teamsProjects)
	}

	//get teams projects technologies
	var teamsProjectsTechnologies []models.Technologies
	for _, project := range teamsProjects {
		err := db.Model(&project).Association("Technologies").Find(&teamsProjectsTechnologies)
		if err != nil {
			fmt.Println("get teams projects technologies: find:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("get teams projects technologies: find: " + err.Error()))
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(teamsWithMostCommonTechnologiesAndMostCommonTechnologiesInProjects)
	if err != nil {
		fmt.Println("recommend team: encode:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("recommend team: encode: " + err.Error()))
		return
	}
}

func AcceptUserToTeam(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	//get user
	user := models.Users{}
	var parseAccept models.ParseAccept

	err := json.NewDecoder(r.Body).Decode(&parseAccept)
	if err != nil {
		fmt.Println("accept user to team: decode:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("accept user to team: decode: " + err.Error()))
		return
	}
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

	var parseUser models.Users

	db.Where("id = ?", parseAccept.UserID).First(&parseUser)
	if parseUser.ID == 0 {
		fmt.Println("accept user to team: user not found")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("accept user to team: user not found"))
		return
	}

	if parseUser.TeamID != 0 {
		fmt.Println("accept user to team: user already in team")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("accept user to team: user already in team"))
		return
	}

	//get team
	team := models.Team{}
	db.Where("id = ?", parseAccept.TeamID).First(&team)

	//if user.ID == parseUser.ID => user is accepting an invitation to join a team
	//if user.ID == parseUser.ID => user is accepting parseUser to join user's team (user is a team leader)
	if user.ID == parseUser.ID {
		//accept user to team
		db.Model(&parseUser).Update("team_id", team.ID)
	} else {
		//accept parseUser to team
		if user.RoleID == 2 {
			db.Model(&parseUser).Update("team_id", team.ID)
		} else {
			fmt.Println("accept user to team: user is not a team leader")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("accept user to team: user is not a team leader"))
			return
		}
	}

	//delete invitation
	db.Where("user_id = ? AND team_id = ?", parseUser.ID, team.ID).Delete(&models.Invite{})
}

func GetTeams(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// json : { "teams" : [{ "id" : 1, "name" : "team1", "logo" : "https://cdn.thebrandingjournal.com/wp-content/uploads/2019/05/chanel_logo_the_branding_journal.jpg", "members": [{ "id" : 1, "firstName" : "John", "lastName" : "Doe", "profilePicture" : "https://cdn.thebrandingjournal.com/wp-content/uploads/2019/05/chanel_logo_the_branding_journal.jpg", "role" : 3, "grade" : 11, "class" : "А", "email" : "martin@bozhilov.me", "discordUsername" : "TechXTT", "discordDiscriminator" : "0196", "github" : "TechXTT"}, ...]}, ...]}

	//get teams from db
	// with Query = "SELECT * FROM teams" we get all teams from db
	var parseTeams []models.ParseTeamView
	db.Raw("SELECT * FROM teams").Scan(&parseTeams)

	// parse teams to teams
	var teams []models.TeamsView

	//every row in parseTeams is a member of a team
	for _, parseTeam := range parseTeams {
		//check if the team is already in teams
		teamAlreadyInTeams := false
		for i, team := range teams {
			if team.ID == parseTeam.ID {
				//if the team is already in teams, add the member to the team
				teamAlreadyInTeams = true

				var member models.MemberView
				member.ID = parseTeam.UID
				member.Name = parseTeam.FirstName + " " + parseTeam.LastName
				member.ProfilePicture = parseTeam.ProfilePicture
				member.Role = parseTeam.Role
				member.Class = parseTeam.Grade + parseTeam.Class
				member.Email = parseTeam.Email
				member.Discord = parseTeam.UserName + "#" + parseTeam.Discriminator
				member.Github = parseTeam.Login
				teams[i].Members = append(teams[i].Members, member)
			}
		}

		//if the team is not in teams, add it
		if !teamAlreadyInTeams {
			teams = append(teams, models.TeamsView{
				ID:           parseTeam.ID,
				Name:         parseTeam.Name,
				Logo:         parseTeam.Logo,
				Members:      []models.MemberView{},
				Project:      models.ProjectView{},
				Technologies: []string{},
			})

			//add the member to the team
			var member models.MemberView
			member.ID = parseTeam.UID
			member.Name = parseTeam.FirstName + " " + parseTeam.LastName
			member.ProfilePicture = parseTeam.ProfilePicture
			member.Role = parseTeam.Role
			member.Class = parseTeam.Grade + parseTeam.Class
			member.Email = parseTeam.Email
			member.Discord = parseTeam.UserName + "#" + parseTeam.Discriminator
			member.Github = parseTeam.Login

			teams[len(teams)-1].Members = append(teams[len(teams)-1].Members, member)

			//get team project
			if parseTeam.PID != 0 {
				var teamProject models.Project
				db.Table("projects").Where("id = ?", parseTeam.PID).First(&teamProject)

				//get team project technologies
				var teamProjectTechnologies []models.Technologies
				db.Table("technologies").Joins("JOIN project_technologies ON project_technologies.project_id = ?", parseTeam.PID).Where("project_technologies.technology_id = technologies.id").Find(&teamProjectTechnologies)

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
			var teamTechnologies []models.Technologies
			db.Table("technologies").Joins("JOIN team_technologies ON team_technologies.team_id = ?", parseTeam.ID).Where("team_technologies.technology_id = technologies.id").Find(&teamTechnologies)

			//parse team technologies
			var teamTechnologiesParsed []string
			for _, teamTechnology := range teamTechnologies {
				teamTechnologiesParsed = append(teamTechnologiesParsed, teamTechnology.Technology)
			}

			//add team technologies to team
			teams[len(teams)-1].Technologies = teamTechnologiesParsed
		}
	}

	//return teams
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(teams)
	if err != nil {
		fmt.Println("get teams: encode:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("get teams: encode: " + err.Error()))
		return
	}
}
