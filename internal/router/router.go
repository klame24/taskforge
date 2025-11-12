package router

import (
	"net/http"
	"taskforge/internal/handlers"
)

type Router struct {
	mux         *http.ServeMux
	userHandler *handlers.UserHandler
}

func NewRouter(userHandler *handlers.UserHandler) *Router {
	return &Router{
		mux:         http.NewServeMux(),
		userHandler: userHandler,
	}
}

func (r *Router) SetupRoutes() *http.ServeMux {
	r.mux.HandleFunc("POST /api/v1/register", r.userHandler.Register)
	r.mux.HandleFunc("POST /api/v1/login", r.userHandler.Login)

	return r.mux
}
