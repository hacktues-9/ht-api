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
	FirstName        string   `json:"first_name"`
	LastName         string   `json:"last_name"`
	Email            string   `json:"email"`
	ElsysEmail       string   `json:"elsys_email"`
	Mobile           string   `json:"mobile"`
	Password         string   `json:"password"`
	Class            string   `json:"class"`
	Grade            string   `json:"grade"`
	ShirtSize        string   `json:"shirt_size"`
	EatingPreference string   `json:"eating_preference"`
	Allergies        []string `json:"allergies"`
	Technologies     []string `json:"technologies"`
}

type LoginUser struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}
