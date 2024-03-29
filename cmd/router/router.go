package router

import (
	"fmt"
	"github.com/hacktues-9/API/cmd/admins"
	"github.com/hacktues-9/API/cmd/mentors"
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
	admin := r.PathPrefix("/admin").Subrouter()
	mentor := r.PathPrefix("/mentor").Subrouter()
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

	mentor.HandleFunc("/discord", func(w http.ResponseWriter, r *http.Request) { // route - /api/mentor/discord
		oauth.GetMentorDiscordInfo(w, r, DB)
	})

	mentor.HandleFunc("/save/{mentor_id}", func(w http.ResponseWriter, r *http.Request) { // route - /api/mentor/save/{team_id}/{mentor_id}
		mentors.SaveMentor(w, r, DB)
	})

	mentor.HandleFunc("/isAvailable/{mentor_id}", func(w http.ResponseWriter, r *http.Request) { // route - /api/mentor/isAvailable/{mentor_id}
		mentors.IsAvailable(w, r, DB)
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

	user.HandleFunc("/reset/{token}", func(w http.ResponseWriter, r *http.Request) { // route - /api/user/reset/{token}
		users.ResetPassword(w, r, DB)
	})

	admin.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) { // route - /api/admin/search
		admins.SearchWithFilters(w, r, DB)
	})
	auth.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) { // route - /api/auth/login
		users.Login(w, r, DB)
	})

	auth.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) { // route - /api/auth/refresh
		jwt.RefreshAccessToken(w, r, DB)
	})

	auth.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) { // route - /api/auth/logout
		users.Logout(w)
	})

	auth.HandleFunc("/forgot/{elsys_email}", func(w http.ResponseWriter, r *http.Request) { // route - /api/auth/forgot/{elsys_email}
		users.ForgotPassword(w, r, DB)
	})

	//auth.HandleFunc("/delete/{token}", func(w http.ResponseWriter, r *http.Request) { // route - /api/auth/delete/{token}
	//	users.DeleteUser(w, r, DB)
	//})

	auth.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) { // route - /api/auth/me
		users.GetUserID(w, r, DB)
	})

	auth.HandleFunc("/check/email/{email}", func(w http.ResponseWriter, r *http.Request) { // route - /api/auth/check/email/{email}
		users.CheckEmail(w, r, DB)
	})

	auth.HandleFunc("/check/elsys_email/{email}", func(w http.ResponseWriter, r *http.Request) { // route - /api/auth/check/elsys_email/{email}
		users.CheckElsysEmail(w, r, DB)
	})

	auth.HandleFunc("/isVerified", func(w http.ResponseWriter, r *http.Request) { // route - /api/auth/isVerified/{elsys_email}
		users.IsVerified(w, r, DB)
	})

	user.HandleFunc("/get/role/{id}", func(w http.ResponseWriter, r *http.Request) { // route - /api/user/get/role/{id}
		users.GetUserRole(w, r, DB)
	})

	user.HandleFunc("/get/{id}", func(w http.ResponseWriter, r *http.Request) { // route - /api/user/get
		users.GetUser(w, r, DB)
	})

	user.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) { // route - /api/user/update
		users.UpdateUser(w, r, DB)
	})

	team.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/create
		w.WriteHeader(http.StatusInternalServerError)
	})

	//team.HandleFunc("/invite", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/invite
	//	teams.InviteUserToTeam(w, r, DB)
	//})

	//team.HandleFunc("/apply", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/apply
	//	teams.ApplyToTeam(w, r, DB)
	//})

	//team.HandleFunc("/accept/{teamId}/{userId}", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/accept/{teamId}/{userId}
	//	teams.AcceptInvite(w, r, DB)
	//})

	//team.HandleFunc("/decline/{teamId}/{userId}", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/decline/{teamId}/{userId}
	//	teams.DeclineInvite(w, r, DB)
	//})

	team.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/get
		teams.GetTeams(w, r, DB)
	})

	team.HandleFunc("/users/search", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/users/search?search={search}
		teams.SearchInvitees(w, r, DB)
	})

	team.HandleFunc("/users/in-team/{id}", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/users/in-team/{id}
		teams.IsUserInTeam(w, r, DB)
	})

	team.HandleFunc("/get/invitees/{id}", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/get/invitees/{id}
		teams.GetInvitees(w, r, DB)
	})

	team.HandleFunc("/get/mentor/{team_id}", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/get/mentor/{team_id}
		mentors.HasMentor(w, r, DB)
	})

	user.HandleFunc("/get/mentor/{user_id}", func(w http.ResponseWriter, r *http.Request) { // route - /api/user/get/mentor/{user_id}
		mentors.HasMentorUID(w, r, DB)
	})

	team.HandleFunc("/get/{id}", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/get/{id}
		teams.GetTeam(w, r, DB)
	})

	team.HandleFunc("/captain/{id}", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/captain
		teams.GetCaptainID(w, r, DB)
	})

	team.HandleFunc("/kick/{id}", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/kick/{id}
		teams.KickUser(w, r, DB)
	})

	team.HandleFunc("/leave/{id}", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/leave/{id}
		teams.LeaveTeam(w, r, DB)
	})

	team.HandleFunc("/delete/{id}", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/delete/{id}
		teams.DeleteTeam(w, r, DB)
	})

	team.HandleFunc("/update/captain/{id}", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/update/captain/{id}
		teams.UpdateCaptain(w, r, DB)
	})

	team.HandleFunc("/update/{id}", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/update/{id}
		teams.UpdateTeam(w, r, DB)
	})

	team.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) { // route - /api/team/{id}
		teams.GetTeamID(w, r, DB)
	})

	user.HandleFunc("/notifications", func(w http.ResponseWriter, r *http.Request) { // route - /api/user/notifications
		users.GetNotifications(w, r, DB)
	})

	r.HandleFunc("/image/{name}", func(w http.ResponseWriter, r *http.Request) { // route - /api/image/{name}
		users.GenerateImage(w, r, DB)
	})

	admin.HandleFunc("/get/teams", func(w http.ResponseWriter, r *http.Request) { // route - /api/admins/get/teams
		admins.GetTeams(w, r, DB)
	})

	mentor.HandleFunc("/get/mentors", func(w http.ResponseWriter, r *http.Request) { // route - /api/mentors/get/mentors
		mentors.GetMentors(w, r, DB)
	})

	admin.HandleFunc("/resend/verification/elsys/{id}", func(w http.ResponseWriter, r *http.Request) { // route - /api/admins/resend/verification/elsys/{id}
		admins.ResendVerificationElsys(w, r, DB)
	})

	admin.HandleFunc("/resend/verification/{id}", func(w http.ResponseWriter, r *http.Request) { // route - /api/admins/resend/verification/{id}
		admins.ResendVerification(w, r, DB)
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://hacktues.com", "http://localhost:8080", "https://*.vercel.app", "https://hacktues.bg", "https://*.hacktues.bg", "http://localhost:3000/", "http://localhost:8080/", "https://*.vercel.app/", "https://*.hacktues.bg/", "https://hacktues.bg/"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
	})

	handler := c.Handler(r)
	http.ListenAndServe(":8080", handler)
	fmt.Println("Server started on port 8080")
}
