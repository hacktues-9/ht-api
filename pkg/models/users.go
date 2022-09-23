package models

import (
	"time"

	"gorm.io/gorm"
)

type Class struct {
	gorm.Model
	Name string `gorm:"unique, not null"`
}

type Role struct {
	gorm.Model
	Name string `gorm:"unique, not null"`
}

type EatingPreference struct {
	gorm.Model
	Name string `gorm:"unique, not null"`
}

type Allergies struct {
	gorm.Model
	Name string `gorm:"unique, not null"`
}

type ShirtSize struct {
	gorm.Model
	Name string `gorm:"unique, not null"`
}

type Discord struct {
	gorm.Model
	DiscordUserID    string `json:"id"`
	Username         string `json:"username"`
	Avatar           string `json:"avatar"`
	AvatarDecoration string `json:"avatar_decoration"`
	Discriminator    string `json:"discriminator"`
	PublicFlags      int    `json:"public_flags"`
	Flags            int    `json:"flags"`
	Banner           string `json:"banner"`
	BannerColor      string `json:"banner_color"`
	AccentColor      int    `json:"accent_color"`
	Locale           string `json:"locale"`
	MfaEnabled       bool   `json:"mfa_enabled"`
}

type Socials struct {
	gorm.Model
	GithubLink    string
	LinkedInLink  string
	InstagramLink string

	ProfilePicture string `gorm:"default:https://i.stack.imgur.com/l60Hf.png"`
	DiscordID      uint
	Discord        Discord `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL; foreignKey:DiscordID"`
}

type Info struct {
	gorm.Model
	Grade   int   `gorm:"not null, check:grade > 7 and grade < 13"`
	ClassID uint  `gorm:"not null"`
	Class   Class `gorm:"foreignKey:ClassID"`

	EatingPreferenceID uint             `gorm:"not null"`
	EatingPreference   EatingPreference `gorm:"foreignKey:EatingPreferenceID"`
	SocialsID          uint             `gorm:"unique, not null"`
	Socials            Socials          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ShirtSizeID        uint             `gorm:"not null"`
	ShirtSize          ShirtSize        `gorm:"foreignKey:ShirtSizeID"`
}

type Security struct {
	gorm.Model
	EmailVerified      bool `gorm:"default:false"`
	ElsysEmailVerified bool `gorm:"default:false"`
	ManualVerified     bool `gorm:"default:false"`
}

type User struct {
	gorm.Model
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`

	Email      string `gorm:"unique"`
	ElsysEmail string `gorm:"unique, not null"`
	Mobile     string `gorm:"unique, not null"`

	Password string `gorm:"not null"`

	InfoID     uint     `gorm:"unique"`
	Info       Info     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL; foreignKey:InfoID"`
	SecurityID uint     `gorm:"unique"`
	Security   Security `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL; foreignKey:SecurityID"`
	RoleID     uint     `gorm:"not null"`
	Role       Role     `gorm:"foreignKey:RoleID"`
	TeamID     uint
	Team       Team `gorm:"foreignKey:TeamID"`

	LastLogin time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type InfoAllergies struct {
	gorm.Model
	InfoID      uint      `gorm:"not null"`
	Info        Info      `gorm:"foreignKey:InfoID"`
	AllergiesID uint      `gorm:"not null"`
	Allergies   Allergies `gorm:"foreignKey:AllergiesID"`
}

type UserTechnologies struct {
	gorm.Model
	UserID         uint         `gorm:"not null"`
	User           User         `gorm:"foreignKey:UserID"`
	TechnologiesID uint         `gorm:"not null"`
	Technologies   Technologies `gorm:"foreignKey:TechnologiesID"`
}
