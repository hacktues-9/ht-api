package models

type MemberView struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	ProfilePicture string `json:"profile_picture"`
	Role           string `json:"role"`
	Class          string `json:"class"`
	Email          string `json:"email"`
	Discord        string `json:"discord"`
	Github         string `json:"github"`
}

type ProjectView struct {
	ID           uint     `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Logo         string   `json:"logo"`
	Github       string   `json:"github"`
	Website      string   `json:"website"`
	Technologies []string `json:"technologies"`
	Photos       []string `json:"photos"`
}

type TeamsView struct {
	ID           uint         `json:"id"`
	Name         string       `json:"name"`
	Logo         string       `json:"logo"`
	Members      []MemberView `json:"members"`
	Project      ProjectView  `json:"project"`
	Technologies []string     `json:"technologies"`
	IsVerified   bool         `json:"isVerified"`
}

type SearchView struct {
	ID                 uint   `json:"id"`
	Name               string `json:"name"`
	ProfilePicture     string `json:"profile_picture"`
	IsInvited          bool   `json:"isInvited"`
	ElsysEmail         string `json:"elsys_email"`
	ElsysEmailVerified bool   `json:"elsys_email_verified"`
}

type UserView struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`

	Email      string `json:"email"`
	ElsysEmail string `json:"elsysEmail"`
	Mobile     string `json:"mobile"`

	Class     string `json:"class"`
	ShirtSize string `json:"shirtSize"`

	EmailVerified      bool `json:"emailVerified"`
	ElsysEmailVerified bool `json:"elsysEmailVerified"`
	ProfilePicVerified bool `json:"profilePicVerified"`

	Discord        string `json:"discord"`
	Github         string `json:"github"`
	LookingForTeam bool   `json:"lookingForTeam"`

	ProfilePicture string   `json:"profilePicture"`
	Technologies   []string `json:"technologies"`
}

type MemberTeamView struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Role   string `json:"role"`
}

type MentorTeamView struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type ProjectTeamView struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Links       struct {
		Github  string `json:"github"`
		Website string `json:"website"`
	} `json:"links"`
}

type GetTeamView struct {
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Logo         string           `json:"logo"`
	Technologies []string         `json:"technologies"`
	Members      []MemberTeamView `json:"members"`
	Mentor       MentorTeamView   `json:"mentor"`
	Project      ProjectTeamView  `json:"project"`
}

type KurView struct {
	Class              string `json:"class"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	Email              string `json:"email"`
	ElsysEmail         string `json:"elsys_email"`
	Mobile             string `json:"mobile"`
	ShirtSize          string `json:"shirt_size"`
	EatingPreference   string `json:"eating_preference"`
	EmailVerified      bool   `json:"email_verified"`
	ElsysEmailVerified bool   `json:"elsys_email_verified"`
	ManualVerified     bool   `json:"manual_verified"`
	Discord            string `json:"discord"`
	Github             string `json:"github"`
	Team               string `json:"team"`
}
