package models

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name        string `gorm:"unique, not null"`
	Description string

	GithubLink string `gorm:"unique"`
	Votes      int    `gorm:"default:0"`

	Logo  string
	Color string
}

type Pictures struct {
	gorm.Model
	ProjectID uint    `gorm:"not null"`
	Project   Project `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL; foreignKey:ProjectID"`
	Picture   string
}

type Invites struct {
	gorm.Model
	LookingFor bool `gorm:"default:false"`
	Accepts    bool `gorm:"default:false"`
}

type Team struct {
	gorm.Model
	Name        string `gorm:"unique, not null"`
	Description string
	ProjectID   uint
	Project     Project `gorm:"foreignKey:ProjectID"`
	InvitesID   uint
	Invites     Invites `gorm:"foreignKey:InvitesID"`

	Logo  string
	Color string

	Approved bool `gorm:"default:false"`
	Points   int  `gorm:"default:0"`
}

type Invite struct {
	gorm.Model
	TeamID uint `gorm:"not null"`
	Team   Team `gorm:"foreignKey:TeamID"`
	UserID uint `gorm:"not null"`
	User   User `gorm:"foreignKey:UserID"`

	Pending     bool `gorm:"default:true"`
	Application bool `gorm:"default:false"`
}

type TeamTechnologies struct {
	gorm.Model
	TeamID         uint         `gorm:"not null"`
	Team           Team         `gorm:"foreignKey:TeamID"`
	TechnologiesID uint         `gorm:"not null"`
	Technologies   Technologies `gorm:"foreignKey:TechnologiesID"`
}

type ProjectTechnologies struct {
	gorm.Model
	ProjectID      uint         `gorm:"not null"`
	Project        Project      `gorm:"foreignKey:ProjectID"`
	TechnologiesID uint         `gorm:"not null"`
	Technologies   Technologies `gorm:"foreignKey:TechnologiesID"`
}

type InviteTechnologies struct {
	gorm.Model
	InviteID       uint         `gorm:"not null"`
	Invite         Invite       `gorm:"foreignKey:InviteID"`
	TechnologiesID uint         `gorm:"not null"`
	Technologies   Technologies `gorm:"foreignKey:TechnologiesID"`
}

type InvitesTechnologies struct {
	gorm.Model
	InvitesID      uint         `gorm:"not null"`
	Invites        Invites      `gorm:"foreignKey:InvitesID"`
	TechnologiesID uint         `gorm:"not null"`
	Technologies   Technologies `gorm:"foreignKey:TechnologiesID"`
}
