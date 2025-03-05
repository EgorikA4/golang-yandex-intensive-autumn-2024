package utils

import (
	"encoding/json"
	"net/http"

	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/consts"
)

type Response interface{}

func SendResponse(w http.ResponseWriter, statusCode int, response Response) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, consts.InternalServerError, http.StatusInternalServerError)
	}
}
