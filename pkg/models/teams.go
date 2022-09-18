package models

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name        string
	Description string

	GithubLink string
	Votes      int

	Logo  string
	Color string
}

type Pictures struct {
	gorm.Model
	ProjectID uint
	Project   Project
	Picture   string
}

type Invites struct {
	gorm.Model
	LookingFor bool
	Accepts    bool
}

type Team struct {
	gorm.Model
	Name        string
	Description string
	ProjectID   uint
	Project     Project
	InvitesID   uint
	Invites     Invites

	Logo  string
	Color string

	Approved bool
	Points   int
}

type Invite struct {
	gorm.Model
	TeamID uint
	Team   Team
	UserID uint
	User   User

	Pending     bool
	Application bool
}

type TeamTechnologies struct {
	gorm.Model
	TeamID         uint
	Team           Team
	TechnologiesID uint
	Technologies   Technologies
}

type ProjectTechnologies struct {
	gorm.Model
	ProjectID      uint
	Project        Project
	TechnologiesID uint
	Technologies   Technologies
}

type InviteTechnologies struct {
	gorm.Model
	InviteID       uint
	Invite         Invite
	TechnologiesID uint
	Technologies   Technologies
}

type InvitesTechnologies struct {
	gorm.Model
	InvitesID      uint
	Invites        Invites
	TechnologiesID uint
	Technologies   Technologies
}
