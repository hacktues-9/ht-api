package router

import (
	"fmt"
	"github.com/hacktues-9/API/pkg/models"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hacktues-9/API/cmd/teams"
	"github.com/hacktues-9/API/cmd/users"
	db "github.com/hacktues-9/API/pkg/database"
	"github.com/hacktues-9/API/pkg/email"
	"github.com/hacktues-9/API/pkg/jwt"
	"github.com/hacktues-9/API/pkg/oauth"
	"github.com/rs/cors"
	"gorm.io/gorm"
)

func Init(DB *gorm.DB) {
	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.Use(mux.CORSMethodMiddleware(r))
	auth := r.PathPrefix("/auth").Subrouter()
	//admin := r.PathPrefix("/admin").Subrouter()
	// mentor := r.PathPrefix("/mentor").Subrouter()
	team := r.PathPrefix("/team").Subrouter()
	user := r.PathPrefix("/user").Subrouter()
	database := r.PathPrefix("/db").Subrouter()

	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { // route - /api/ping
		w.WriteHeader(http.StatusOK)
		models.RespHandler(w, r, models.DefaultPosResponse("pong"), nil, http.StatusOK, "ping")
	})

	user.HandleFunc("/discord", func(w http.ResponseWriter, r *http.Request) { // route - /api/user/discord
		oauth.GetDiscordInfo(w, r, DB)
	})

	user.HandleFunc("/github", func(w http.ResponseWriter, r *http.Request) { // route - /api/user/github
		oauth.GetGithubInfo(w, r, DB)
	})

	auth.HandleFunc("/github", func(w http.ResponseWriter, r *http.Request) { // route - /api/auth/github"
		oauth.LoginGithub(w, r, DB)
	})

	auth.HandleFunc("/discord", func(w http.ResponseWriter, r *http.Request) { // route - /api/auth/discord"
		oauth.LoginDiscord(w, r, DB)
	})

	database.HandleFunc("/migrate", func(w http.ResponseWriter, r *http.Request) { // route - /api/db/migrate
		db.Migrate(DB)
	})

	database.HandleFunc("/drop", func(w http.ResponseWriter, r *http.Request) { // route - /api/db/drop
		db.Drop(DB)
	})

	database.HandleFunc("/populate", func(w http.ResponseWriter, r *http.Request) { // route - /api/db/populate
		db.PopulateDefault(DB)
	})

	auth.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) { // route - /api/auth/register
		users.Register(w, r, DB)
	})

	user.HandleFunc("/verify/{elsys}/{token}", func(w http.ResponseWriter, r *http.Request) { // route - /api/user/verify/{elsys}/{token}
		email.ValidateEmail(w, r, DB)
	})

	//admin.HandleFunc("/search-user", func(w http.ResponseWriter, r *http.Request) { // route - /api/admin/search-user
	//	users.FetchUser(w, r, DB)
	//})

	auth.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) { // route - /api/auth/login
		users.Login(w, r, DB)
	})

	auth.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) { // route - /api/auth/refresh
		jwt.RefreshAccessToken(w, r, DB)
	})

	auth.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) { // route - /api/auth/logout
		users.Logout(w)
	})

	auth.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) { // route - /api/auth/me
		users.GetUserID(w, r)
	})

	user.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) { // route - /api/user/get
		users.GetUser(w, r, DB)
	})

	team.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/create
		teams.CreateTeam(w, r, DB)
	})

	team.HandleFunc("/invite", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/invite
		teams.InviteUserToTeam(w, r, DB)
	})

	team.HandleFunc("/apply", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/apply
		teams.ApplyToTeam(w, r, DB)
	})

	team.HandleFunc("/accept", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/accept
		teams.AcceptUserToTeam(w, r, DB)
	})

	team.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/get
		teams.GetTeams(w, r, DB)
	})

	team.HandleFunc("/users/search", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/users/search?search={search}
		teams.SearchInvitees(w, r, DB)
	})

	team.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/{id}
		teams.GetTeamID(w, r, DB)
	})

	team.HandleFunc("/{id}/captain", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/captain
		teams.GetCaptainID(w, r, DB)
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://hacktues.com", "http://localhost:8080", "https://hacktues.bg"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
	})

	handler := c.Handler(r)
	http.ListenAndServe(":8080", handler)
	fmt.Println("Server started on port 8080")
}
