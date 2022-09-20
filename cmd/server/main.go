package main

import (
	"net/http"

	db "github.com/hacktues-9/API/pkg/database"

	"github.com/hacktues-9/API/cmd/users"
	"github.com/hacktues-9/API/pkg/discord"
)

func main() {
	DB := db.Init()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/discord", func(w http.ResponseWriter, r *http.Request) {
		discord.GetDiscordInfo(w, r, DB)
	})

	mux.HandleFunc("/api/db/migrate", func(w http.ResponseWriter, r *http.Request) {
		db.Migrate(DB)
	})

	mux.HandleFunc("/api/db/drop", func(w http.ResponseWriter, r *http.Request) {
		db.Drop(DB)
	})

	mux.HandleFunc("/api/db/populate", func(w http.ResponseWriter, r *http.Request) {
		db.PopulateDefault(DB)
	})

	mux.HandleFunc("/api/auth/register", func(w http.ResponseWriter, r *http.Request) {
		users.Register(w, r, DB)
	})

	mux.HandleFunc("/api/auth/getUser", func(w http.ResponseWriter, r *http.Request) {
		users.FetchUser(w, r, DB)
	})

	http.ListenAndServe(":8080", mux)
}
