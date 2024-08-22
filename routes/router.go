package routes

import (
	"chat-backend/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/login", handlers.LoginUser).Methods("POST")
	r.HandleFunc("/send", handlers.SendMessage).Methods("POST")
	r.HandleFunc("/messages", handlers.GetMessages).Methods("GET")
}
