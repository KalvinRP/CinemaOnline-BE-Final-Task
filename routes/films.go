package routes

import (
	handlers "finaltask/handlers"
	"finaltask/pkg/middleware"
	"finaltask/pkg/mysql"
	repositories "finaltask/repositories"

	"github.com/gorilla/mux"
)

func FilmsRoutes(r *mux.Router) {
	filmsRepository := repositories.RepoFilms(mysql.DB)
	h := handlers.HandlerFilms(filmsRepository)

	r.HandleFunc("/films", h.AllFilms).Methods("GET")
	r.HandleFunc("/topfilms", h.ForBanner).Methods("GET")
	r.HandleFunc("/films/{id}", middleware.Auth(h.GetFilms)).Methods("GET")
	r.HandleFunc("/public-films/{id}", h.GetPublicFilms).Methods("GET")
	r.HandleFunc("/films", middleware.AuthAdmin(middleware.UploadFile(h.MakeFilms))).Methods("POST")
	r.HandleFunc("/films/{id}", middleware.AuthAdmin(middleware.MayUploadFile(h.EditFilms))).Methods("PATCH")
	r.HandleFunc("/films/{id}", middleware.AuthAdmin(h.DeleteFilms)).Methods("DELETE")
}
