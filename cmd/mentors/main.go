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
