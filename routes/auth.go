package routes

import (
	handlers "finaltask/handlers"
	"finaltask/pkg/middleware"
	"finaltask/pkg/mysql"
	repositories "finaltask/repositories"

	"github.com/gorilla/mux"
)

func AuthRoutes(r *mux.Router) {
	userRepository := repositories.RepoUser(mysql.DB)
	h := handlers.HandlerAuth(userRepository)

	r.HandleFunc("/register", h.Register).Methods("POST")
	r.HandleFunc("/verify", h.Verify).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")
	r.HandleFunc("/forget", h.ForgetPass).Methods("POST")
	r.HandleFunc("/reset", h.ResetPassword).Methods("PATCH")
	r.HandleFunc("/check-auth", middleware.Auth(h.CheckAuth)).Methods("GET")
}
