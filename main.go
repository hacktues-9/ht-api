package main

import (
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Set the return Content-Type as JSON
	w.Header().Set("Content-Type", "application/json")

	// Customize the response depending on the method
	switch r.Method {
	case "GET":
		w.Write([]byte("{\"message\": \"GET\"}"))
	case "POST":
		w.Write([]byte("{\"message\": \"POST\"}"))
	case "PUT":
		w.Write([]byte("{\"message\": \"PUT\"}"))
	case "DELETE":
		w.Write([]byte("{\"message\": \"DELETE\"}"))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"message\": \"Not Found\"}"))
	}
}

func main() {
	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(":7000", nil))
}