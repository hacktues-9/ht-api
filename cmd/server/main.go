package main

import (
	// db "github.com/hacktues-9/API/pkg/database"
	"net/http"

	"github.com/hacktues-9/API/pkg/discord"
)

func main() {
	// db.Init()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/discord", discord.GetDiscordInfo)

	http.ListenAndServe(":8080", mux)
}
