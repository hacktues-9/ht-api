package models

import (
	"gorm.io/gorm"
)

type Technologies struct {
	gorm.Model
	Technology  string `json:"technology" gorm:"unique, not null"`
	Description string `json:"description"`

	BgColor   string `json:"bgColor"`
	TextColor string `json:"textColor"`
	Icon      string `json:"icon"`
}

type Log struct {
	gorm.Model
	Url      string
	Endpoint string
	Method   string

	Headers string
	Body    string
	Ip      string
}

type DiscordBearer struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type RegisterUser struct {
	FirstName        string   `json:"firstName"`
	LastName         string   `json:"lastName"`
	Email            string   `json:"email"`
	ElsysEmail       string   `json:"elsysEmail"`
	Mobile           string   `json:"phone"`
	Password         string   `json:"password"`
	ConfirmPassword  string   `json:"confirmPassword"`
	Class            string   `json:"classLetter"`
	Grade            string   `json:"classNumber"`
	ShirtSize        string   `json:"shirtSize"`
	EatingPreference string   `json:"eatingPreferences"`
	Allergies        []string `json:"allergies"`
	Technologies     []string `json:"technologies"`
}

type LoginUser struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type ParseTeam struct {
	Name         string       `json:"teamName"`
	Description  string       `json:"teamDescription"`
	Technologies []string     `json:"teamTechnologies"`
	Invitees     []SearchView `json:"teamInvitees"`
}

type ParseInvite struct {
	UserID uint `json:"userId"`
	TeamID uint `json:"teamId"`
}

type ParseApply struct {
	UserID uint `json:"userId"`
	TeamID uint `json:"teamId"`
}

type ParseAccept struct {
	UserID uint `json:"userId"`
	TeamID uint `json:"teamId"`
}

type ParseTeamView struct {
	TID      uint   `json:"tid"`
	Name     string `json:"name"`
	Logo     string `json:"logo"`
	PID      uint   `json:"pid"`
	Approved bool   `json:"approved"`
}

type ParseChangeUser struct {
	ID             uint     `json:"id"`
	Technologies   []string `json:"technologies"`
	LookingForTeam bool     `json:"lookingForTeam"`
}

type Notification struct {
	TeamID   uint   `json:"teamId"`
	TeamName string `json:"teamName"`
	TeamLogo string `json:"teamLogo"`
}
