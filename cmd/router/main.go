package router

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	// Login handler
	mux.HandleFunc("/login", login)
	// Logout handler
	mux.HandleFunc("/logout", logout)
	// Register handler
	mux.HandleFunc("/register", register)
	// Get user handler
	mux.HandleFunc("/refresh", getUser)

	log.Println("Listening on port 8080")
	http.ListenAndServe(":8080", mux)
}