package main

import (
	db "github.com/hacktues-9/API/pkg/database"
	_ "github.com/joho/godotenv/autoload"

	"github.com/hacktues-9/API/cmd/router"
)

func main() {
	DB := db.Init()
	router.Init(DB)
}
