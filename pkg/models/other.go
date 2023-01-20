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
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Technologies []string `json:"technologies"`
	Logo         string   `json:"logo"`
	Color        string   `json:"color"`
}
