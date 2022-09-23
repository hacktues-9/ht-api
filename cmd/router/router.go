package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hacktues-9/API/cmd/users"
	db "github.com/hacktues-9/API/pkg/database"
	"github.com/hacktues-9/API/pkg/discord"
	"github.com/hacktues-9/API/pkg/email"
	"github.com/hacktues-9/API/pkg/jwt"
	"gorm.io/gorm"
)

func Init(DB *gorm.DB) {
	r := mux.NewRouter()

	r.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	r.HandleFunc("/api/discord", func(w http.ResponseWriter, r *http.Request) {
		discord.GetDiscordInfo(w, r, DB)
	})

	r.HandleFunc("/api/db/migrate", func(w http.ResponseWriter, r *http.Request) {
		db.Migrate(DB)
	})

	r.HandleFunc("/api/db/drop", func(w http.ResponseWriter, r *http.Request) {
		db.Drop(DB)
	})

	r.HandleFunc("/api/db/populate", func(w http.ResponseWriter, r *http.Request) {
		db.PopulateDefault(DB)
	})

	r.HandleFunc("/api/auth/register", func(w http.ResponseWriter, r *http.Request) {
		users.Register(w, r, DB)
	})

	r.HandleFunc("/api/auth/verify/{elsys}/{token}", func(w http.ResponseWriter, r *http.Request) {
		email.ValidateEmail(w, r, DB)
	})

	r.HandleFunc("/api/admin/auth/getUser", func(w http.ResponseWriter, r *http.Request) {
		users.FetchUser(w, r, DB)
	})

	r.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		users.Login(w, r, DB)
	})

	r.HandleFunc("/api/auth/refresh", func(w http.ResponseWriter, r *http.Request) {
		jwt.RefreshAccessToken(w, r, DB)
	})

	r.HandleFunc("/api/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		users.Logout(w, r, DB)
	})

	r.HandleFunc("/api/users/me", func(w http.ResponseWriter, r *http.Request) {
		users.GetUser(w, r, DB)
	})

	http.ListenAndServe(":8080", r)
	fmt.Println("Server started on port :8080")
}
