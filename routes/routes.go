package routes

import (
	"github.com/gorilla/mux"
)

func RouteInit(r *mux.Router) {
	UserRoutes(r)
	FilmsRoutes(r)
	GenreRoutes(r)
	AuthRoutes(r)
	TransRoutes(r)
}
