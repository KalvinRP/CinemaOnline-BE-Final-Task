package routes

import (
	handlers "finaltask/handlers"
	"finaltask/pkg/middleware"
	"finaltask/pkg/mysql"
	repositories "finaltask/repositories"

	"github.com/gorilla/mux"
)

func GenreRoutes(r *mux.Router) {
	genreRepository := repositories.RepoGenre(mysql.DB)
	h := handlers.HandlerGenre(genreRepository)

	r.HandleFunc("/genre", h.WholeGenre).Methods("GET")
	r.HandleFunc("/genre", middleware.AuthAdmin(h.MakeGenre)).Methods("POST")
	r.HandleFunc("/genre/{id}", middleware.AuthAdmin(h.EditGenre)).Methods("PATCH")
	r.HandleFunc("/genre/{id}", middleware.AuthAdmin(h.DeleteGenre)).Methods("DELETE")
}
