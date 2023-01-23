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
}
