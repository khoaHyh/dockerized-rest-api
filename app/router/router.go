package router

import (
	"dockerized-rest-api/app/requestlog"
	"dockerized-rest-api/app/router/middleware"
	"dockerized-rest-api/app/server"

	"github.com/go-chi/chi"
)

func New(s *server.Server) *chi.Mux {
	l := s.Logger()

	r := chi.NewRouter()
	r.Method("GET", "/", requestlog.NewHandler(s.HandleIndex, l))

	// Routes for healthz
	r.Get("/healthz/liveness", server.HandleLive)
	r.Method("GET", "/healthz/readiness", requestlog.NewHandler(s.HandleReady, l))

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.ContentTypeJson)

		// Routes for books
		r.Method("GET", "/books", requestlog.NewHandler(s.HandleListBooks, l))
		r.Method("POST", "/books", requestlog.NewHandler(s.HandleCreateBook, l))
		r.Method("GET", "/books/{id}", requestlog.NewHandler(s.HandleReadBook, l))
		r.Method("PUT", "/books/{id}", requestlog.NewHandler(s.HandleUpdateBook, l))
		r.Method("DELETE", "/books/{id}", requestlog.NewHandler(s.HandleDeleteBook, l))
	})

	return r
}
