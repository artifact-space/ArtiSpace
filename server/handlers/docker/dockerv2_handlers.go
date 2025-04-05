package docker

import (
	"fmt"
	"net/http"

	"github.com/artifact-space/ArtiSpace/config"
	"github.com/artifact-space/ArtiSpace/consts"
	"github.com/artifact-space/ArtiSpace/db"
	"github.com/artifact-space/ArtiSpace/handlers/common"
	"github.com/artifact-space/ArtiSpace/log"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Reference: https://docker-docs.uclv.cu/registry/spec/api/#api-version-check
func GetDockerV2APISupport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Docker-Distribution-API-Version", "registry/2.0")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"Docker-Distribution-API-Version": "registry/2.0"}`))
}

// Reference: https://docker-docs.uclv.cu/registry/spec/api/#pushing-an-image
func InitiateBlobUpload(w http.ResponseWriter, r *http.Request) {
	repoName := chi.URLParam(r, "docker_repository_name")
	nameSpace := chi.URLParam(r, "docker_namespace")

	if repoName == "" {
		common.HandleBadRequest(w, consts.ErrBadRequest, "Docker repository name is empty")
		return
	}

	if nameSpace == "" {
		nameSpace = "library"
	}

	sessionId := uuid.New().String()

	config, err := config.Config()
	if err != nil {
		log.Logger().Error().Err(err).Msg("unable to read server configuration")
		common.HandleInternalError(w, consts.ErrLoadingConfig, "")
		return
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	locationUrl := fmt.Sprintf("%s://%s:%d/v2/%s/%s/blob/uploads/%s", scheme, config.Server.HostnameOrNodeIP, config.Server.DockerV2Port, nameSpace, repoName, sessionId)

	err = db.Provider.CreateDockerNamespaceAndRepositoryIfMissing(r.Context(), nameSpace, repoName)
	if err != nil {
		common.HandleInternalError(w, consts.ErrDatabaseSaveFailed, "")
		return
	}

	w.Header().Set("Location", locationUrl)
	w.WriteHeader(http.StatusAccepted)

	w.Write([]byte(`{"Location":"` + locationUrl + `" }`))
}
