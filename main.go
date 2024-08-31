package main

import (
	"chat-backend/database"
	"chat-backend/routes"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	database.ConnectDB()

	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	http.ListenAndServe(":8000", r)
}