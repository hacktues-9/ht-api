package models

import (
	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	Url      string
	Endpoint string
	Method   string

	Headers string
	Body    string
	Ip      string
}
