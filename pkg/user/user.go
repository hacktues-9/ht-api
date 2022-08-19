package modelUser

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	Id        int       `json:"id"`
	name 	  string    `json:"name"`
	email     string    `json:"email"`
	password  string    `json:"password"`
	telNumber string    `json:"tel_number"`
	createdAt time.Time `json:"created_at"`
	verified  bool      `json:"verified"`
	grade     int       `json:"grade"`
	shirtSize string    `json:"shirt_size"`
	preferences string   `json:"preferences"`
	discordId string    `json:"discord_id"`
	profilePic string    `json:"profile_pic"`
	allergies []string    `json:"allergies"`
	technologies []string `json:"technologies"`
	notifications []string `json:"notifications"`
}
