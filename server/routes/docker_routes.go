package routes

import (
	"github.com/artifact-space/ArtiSpace/handlers/docker"
	"github.com/artifact-space/ArtiSpace/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
)

func InitDockerV2Router() *chi.Mux {
	router := chi.NewRouter()

	//Cors configuration
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELTE", "OPTIONS", "HEAD"},
		AllowedHeaders: []string{"*"},
	})

	router.Use(cors.Handler)
	router.Use(httplog.RequestLogger(log.HttpLogger()))

	router.Route("/v2", func(r chi.Router) {
		r.Get("/", docker.GetDockerV2APISupport)

		//Blob
		r.Post("/{docker_namespace}/{docker_repository_name}/blobs/uploads/", docker.InitiateBlobUpload)
	})

	

	return router
}
