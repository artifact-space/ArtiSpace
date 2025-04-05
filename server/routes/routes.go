package routes

import (
	"github.com/artifact-space/ArtiSpace/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
)

func InitRouter() *chi.Mux {
	router := chi.NewRouter()

	//Cors configuration
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "HEAD"},
		AllowedHeaders: []string{"*"},
	})

	router.Use(cors.Handler)
	router.Use(httplog.RequestLogger(log.HttpLogger()))

	//Routes

	return router
}
