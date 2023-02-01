package routes

import (
	"finaltask/handlers"
	"finaltask/pkg/middleware"
	"finaltask/pkg/mysql"
	"finaltask/repositories"

	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router) {
	userRepository := repositories.RepoUser(mysql.DB)
	h := handlers.HandlerUser(userRepository)

	r.HandleFunc("/useracc", middleware.Auth(h.GetAcc)).Methods("GET")
	r.HandleFunc("/user", middleware.Auth(middleware.MayUploadFile(h.EditProfile))).Methods("PATCH")
	r.HandleFunc("/user/{id}", middleware.Auth(h.DeleteAcc)).Methods("DELETE")
}
