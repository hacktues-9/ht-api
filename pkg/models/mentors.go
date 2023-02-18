package models

import (
	"gorm.io/gorm"
)

type Mentors struct {
	gorm.Model
	FirstName      string `gorm:"not null"`
	LastName       string `gorm:"not null"`
	Email          string `gorm:"unique, not null"`
	Mobile         string `gorm:"unique, not null"`
	ProfilePicture string `gorm:"not null"`
	Company        string `gorm:"not null"`
	Position       string `gorm:"not null"`
	Description    string `gorm:"not null"`
	Videos         string `gorm:"not null"`
	Online         bool   `gorm:"not null"`
	OnSite         bool   `gorm:"not null"`
	VerCode        string `gorm:"not null, default:generate_code()"`

	RoleID    uint `gorm:"not null"`
	Role      Role `gorm:"foreignKey:RoleID"`
	TeamID    uint
	Team      Team    `gorm:"foreignKey:TeamID"`
	DiscordID uint    `gorm:"unique"`
	Discord   Discord `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL; foreignKey:DiscordID"`
}

type MentorTechnologies struct {
	gorm.Model
	MentorID       uint         `gorm:"not null"`
	Mentor         Mentors      `gorm:"foreignKey:MentorID"`
	TechnologiesID uint         `gorm:"not null"`
	Technologies   Technologies `gorm:"foreignKey:TechnologiesID"`
}

type TimeFrames struct {
	gorm.Model
	Date      string `gorm:"not null"`
	StartTime string `gorm:"not null"`
	EndTime   string `gorm:"not null"`
}

type MentorsTimeFrames struct {
	gorm.Model
	Mentor       Mentors    `gorm:"foreignKey:MentorID"`
	MentorID     uint       `gorm:"not null"`
	TimeFrames   TimeFrames `gorm:"foreignKey:TimeFramesID"`
	TimeFramesID uint       `gorm:"not null"`
}

type ParseMentorView struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Pos            string `json:"pos"`
	ProfilePicture string `json:"profile_picture"`
	YtVideo        string `json:"yt_video"`
	TeamID         uint   `json:"team_id"`
}
