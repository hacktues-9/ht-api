package models

import "gorm.io/gorm"

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
