package mentors

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hacktues-9/API/cmd/users"
	"github.com/hacktues-9/API/pkg/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func SaveMentor(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	vars := mux.Vars(r)
	teamID, err := strconv.Atoi(vars["team_id"])
	if err != nil {
		fmt.Printf("[ ERROR ] [ SaveMentor ] accept invite: parse teamID: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "accept invite: parse teamID: "+err.Error(), 0), err, http.StatusBadRequest, "SaveMentor")
		return
	}

	mentorID, err := strconv.Atoi(vars["mentor_id"])
	if err != nil {
		fmt.Printf("[ ERROR ] [ SaveMentor ] accept invite: parse mentorID: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "accept invite: parse mentorID: "+err.Error(), 0), err, http.StatusBadRequest, "SaveMentor")
		return
	}

	sub, err := users.ReturnAuthID(w, r, db)
	if err != nil {
		fmt.Printf("[ ERROR ] [ SaveMentor ] %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, err.Error(), 0), err, http.StatusBadRequest, "SaveMentor")
		return
	}

	mentor := models.Mentors{}
	// check if mentor is available
	db.Where("id = ?", mentorID).First(&mentor)

	if mentor.ID == 0 {
		fmt.Printf("[ ERROR ] [ SaveMentor ] mentor not found\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "mentor not found", 0), err, http.StatusBadRequest, "SaveMentor")
		return
	}

	if mentor.TeamID != 0 {
		fmt.Printf("[ ERROR ] [ SaveMentor ] mentor is not available\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "mentor is not available", 0), err, http.StatusBadRequest, "SaveMentor")
		return
	}

	//check if sub is team leader
	var user models.Users
	db.Where("id = ?", sub).First(&user)

	if float64(user.TeamID) != float64(teamID) {
		fmt.Printf("[ ERROR ] [ SaveMentor ] user is not team leader\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "user is not team leader", 0), err, http.StatusBadRequest, "SaveMentor")
		return
	}

	if user.RoleID != 2 {
		fmt.Printf("[ ERROR ] [ SaveMentor ] user is not team leader\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "user is not team leader", 0), err, http.StatusBadRequest, "SaveMentor")
		return
	}

	// set mentor team id
	db.Model(&models.Mentors{}).Where("id = ?", mentor.ID).Update("team_id", teamID)

	models.RespHandler(w, r, models.DefaultPosResponse("success"), nil, http.StatusOK, "SaveMentor")

}

func IsAvailable(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	vars := mux.Vars(r)
	mentorID, err := strconv.Atoi(vars["mentor_id"])
	if err != nil {
		fmt.Printf("[ ERROR ] [ IsAvailable ] accept invite: parse mentorID: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "accept invite: parse mentorID: "+err.Error(), 0), err, http.StatusBadRequest, "IsAvailable")
		return
	}

	mentor := models.Mentors{}
	// check if mentor is available
	db.Where("id = ?", mentorID).First(&mentor)

	if mentor.ID == 0 {
		fmt.Printf("[ ERROR ] [ IsAvailable ] mentor not found\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "mentor not found", 0), err, http.StatusBadRequest, "IsAvailable")
		return
	}

	if mentor.TeamID != 0 {
		fmt.Printf("[ ERROR ] [ IsAvailable ] mentor is not available\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "mentor is not available", 0), err, http.StatusBadRequest, "IsAvailable")
		return
	}

	models.RespHandler(w, r, models.DefaultPosResponse("success"), nil, http.StatusOK, "IsAvailable")
}

func HasMentor(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	vars := mux.Vars(r)
	teamID, err := strconv.Atoi(vars["team_id"])
	if err != nil {
		fmt.Printf("[ ERROR ] [ HasMentor ] accept invite: parse teamID: %v\n", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "accept invite: parse teamID: "+err.Error(), 0), err, http.StatusBadRequest, "HasMentor")
		return
	}

	mentor := models.Mentors{}
	// check if mentor is available
	db.Where("team_id = ?", teamID).First(&mentor)

	if mentor.ID == 0 {
		fmt.Printf("[ ERROR ] [ HasMentor ] mentor not found\n")
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusBadRequest, "mentor not found", 0), err, http.StatusBadRequest, "HasMentor")
		return
	}

	models.RespHandler(w, r, models.DefaultPosResponse(mentor.ID), nil, http.StatusOK, "HasMentor")
}

func GetMentors(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	query := r.URL.Query()
	name := query.Get("sname")
	tech := query.Get("stech")
	var parseMentors []models.Mentors
	// get mentors name like
	if name != "" {
		db.Raw("SELECT * FROM mentors WHERE concat(first_name, ' ', last_name) LIKE '%?%'", name).Scan(&parseMentors)
	} else {
		db.Raw("SELECT * FROM mentors").Scan(&parseMentors)
	}
	var mentors []models.MentorView
	for _, parseMentor := range parseMentors {
		mentor := models.MentorView{
			ID:             parseMentor.ID,
			Name:           parseMentor.FirstName + " " + parseMentor.LastName,
			Description:    parseMentor.Description,
			Position:       parseMentor.Position,
			Technologies:   []string{},
			ProfilePicture: parseMentor.ProfilePicture,
			Video:          parseMentor.Videos,
			TeamID:         parseMentor.TeamID,
			TimeFrames:     []uint{},
			OnSite:         parseMentor.OnSite,
			Online:         parseMentor.Online,
		}

		// get mentor technologies
		var mentorTechs []models.Technologies
		db.Joins("JOIN mentor_technologies ON mentor_technologies.technologies_id = technologies_id").Where("mentor_technologies.mentor_id = ?", parseMentor.ID).Find(&mentorTechs)
		// check if tech is in mentor techs
		for _, mentorTech := range mentorTechs {
			mentor.Technologies = append(mentor.Technologies, mentorTech.Technology)
		}
		// check if tech is in mentor techs
		if tech != "" {
			var isTech bool
			for _, mentorTech := range mentorTechs {
				if mentorTech.Technology == tech {
					isTech = true
				}
			}
			if !isTech {
				continue
			}
		}

		// get mentor time frames
		var mentorTimeFrames []models.TimeFrames
		db.Joins("JOIN mentor_time_frames ON mentor_time_frames.time_frames_id = time_frames_id").Where("mentor_time_frames.mentor_id = ?", parseMentor.ID).Find(&mentorTimeFrames)
		for _, mentorTimeFrame := range mentorTimeFrames {
			mentor.TimeFrames = append(mentor.TimeFrames, mentorTimeFrame.ID)
		}

		mentors = append(mentors, mentor)
	}
	models.RespHandler(w, r, models.DefaultPosResponse(mentors), nil, http.StatusOK, "GetMentors")
}
