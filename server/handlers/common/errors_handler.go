package common

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"error_message"`
	Code    int    `json:"error_code"`
}

func writeError(w http.ResponseWriter, statusCode int, errorCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Code: errorCode, Message:  errorMsg})
}

func HandleInternalError(w http.ResponseWriter, errorCode int, errorMessage string) {
	writeError(w, http.StatusInternalServerError, errorCode, errorMessage)
}

func HandleBadRequest(w http.ResponseWriter, errorCode int, errorMessage string) {
	writeError(w, http.StatusBadRequest, errorCode, errorMessage)
}
