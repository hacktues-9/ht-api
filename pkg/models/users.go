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
	DiscordUserID    string `json:"id" gorm:"unique, not null"`
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

type Github struct {
	gorm.Model
	Login        string `json:"login"`
	GithubUserID int    `json:"id" gorm:"unique, not null"`
	NodeID       string `json:"node_id"`

	AvatarURL  string `json:"avatar_url"`
	GravatarID string `json:"gravatar_id"`

	URL     string `json:"url"`
	HTMLURL string `json:"html_url"`

	FollowersURL string `json:"followers_url"`
	FollowingURL string `json:"following_url"`

	GistsURL         string `json:"gists_url"`
	StarredURL       string `json:"starred_url"`
	SubscriptionsURL string `json:"subscriptions_url"`

	OrganizationsURL string `json:"organizations_url"`
	ReposURL         string `json:"repos_url"`

	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`

	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`

	Name    string `json:"name"`
	Company string `json:"company"`

	Blog     string `json:"blog"`
	Location string `json:"location"`

	Email string `json:"email"`

	Hireable        bool   `json:"hireable"`
	Bio             string `json:"bio"`
	TwitterUsername string `json:"twitter_username"`

	PublicRepos int `json:"public_repos"`
	PublicGists int `json:"public_gists"`

	Followers int `json:"followers"`
	Following int `json:"following"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Socials struct {
	gorm.Model
	LinkedInLink  string
	InstagramLink string

	ProfilePicture string  `gorm:"default:https://api.hacktues.bg/api/image/John%20Doe"`
	DiscordID      uint    `gorm:"unique, not null"`
	Discord        Discord `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL; foreignKey:DiscordID"`
	GithubID       uint    `gorm:"unique, not null"`
	Github         Github  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL; foreignKey:GithubID"`
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

type Users struct {
	gorm.Model
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`

	Email      string `gorm:"unique"`
	ElsysEmail string `gorm:"unique, not null"`
	Mobile     string `gorm:"unique, not null"`

	Password string `gorm:"not null"`

	LookingForTeam bool `gorm:"default:false"`

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
	User           Users        `gorm:"foreignKey:UserID"`
	TechnologiesID uint         `gorm:"not null"`
	Technologies   Technologies `gorm:"foreignKey:TechnologiesID"`
}
