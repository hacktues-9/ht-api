package models

import (
	"gorm.io/gorm"
)

type Technologies struct {
	gorm.Model
	Technology  string
	Description string

	BgColor   string
	TextColor string
	Icon      string
}
