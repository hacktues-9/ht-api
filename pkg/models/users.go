package models

import (
	"time"

	"gorm.io/gorm"
)

type Class struct {
	gorm.Model
	Name string
}

type Role struct {
	gorm.Model
	Name string
}

type EatingPreference struct {
	gorm.Model
	Name string
}

type Allergies struct {
	gorm.Model
	Name string
}

type ShirtSize struct {
	gorm.Model
	Name string
}

type Discord struct {
	gorm.Model
	DiscordID string
	Email     string
	Token     string

	ExpiresIn int
	CreatedAt time.Time
}

type Socials struct {
	gorm.Model
	GithubLink    string
	LinkedInLink  string
	InstagramLink string

	ProfilePicture string
	DiscordID      uint
	Discord        Discord
}

type Info struct {
	gorm.Model
	Grade   int
	ClassID uint
	Class   Class

	EatingPreferenceID uint
	EatingPreference   EatingPreference
	SocialsID          uint
	Socials            Socials
	ShirtSizeID        uint
	ShirtSize          ShirtSize
}

type Security struct {
	gorm.Model
	EmailVerified      bool
	ElsysEmailVerified bool
	ManualVerified     bool
}

type User struct {
	gorm.Model
	FirstName string
	LastName  string

	Email      string
	ElsysEmail string
	Mobile     string

	Password string

	InfoID     uint
	Info       Info
	SecurityID uint
	Security   Security
	RoleID     uint
	Role       Role
	TeamID     uint
	Team       Team

	LastLogin time.Time
}

type InfoAllergies struct {
	gorm.Model
	InfoID      uint
	Info        Info
	Allergies   Allergies
	AllergiesID uint
}

type UserTechnologies struct {
	gorm.Model
	UserID         uint
	User           User
	TechnologiesID uint
	Technologies   Technologies
}
