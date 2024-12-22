package server

import (
	"encoding/json"
	"net/http"

	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/pkg/calculation"
)

type Request struct {
	Expression string `json:"expression"`
}

type ResultResponse struct {
	Result float64 `json:"result"`
}

type ErrResponse struct {
	Error string `json:"error"`
}

type Response interface{}

func sendResponse(w http.ResponseWriter, statusCode int, response Response) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
	}
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		sendResponse(w, http.StatusMethodNotAllowed, ErrResponse{
			Error: MethodNotAllowed,
		})
		return
	}

	request := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		sendResponse(w, http.StatusBadRequest, ErrResponse{
			Error: err.Error(),
		})
		return
	}

	result, err := calculation.Calc(request.Expression)
	if err != nil {
		sendResponse(w, http.StatusBadRequest, ErrResponse{
			Error: err.Error(),
		})
	} else {
		sendResponse(w, http.StatusOK, ResultResponse{
			Result: result,
		})
	}
}
