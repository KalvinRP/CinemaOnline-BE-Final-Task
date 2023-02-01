package routes

import (
	"finaltask/handlers"
	"finaltask/pkg/middleware"
	"finaltask/pkg/mysql"
	"finaltask/repositories"

	"github.com/gorilla/mux"
)

func TransRoutes(r *mux.Router) {
	transRepository := repositories.RepoTrans(mysql.DB)
	h := handlers.HandlerTransaction(transRepository)

	r.HandleFunc("/transaction", h.GetAllTrans).Methods("GET")
	r.HandleFunc("/transaction/{id}", h.GetOneTrans).Methods("GET")
	r.HandleFunc("/account", middleware.Auth(h.UserHistory)).Methods("GET")
	r.HandleFunc("/own-films", middleware.Auth(h.UserFilms)).Methods("GET")
	r.HandleFunc("/transactions", middleware.Auth(h.MakeMoreTrans)).Methods("POST")
	r.HandleFunc("/notification", h.Notification).Methods("POST")
}
