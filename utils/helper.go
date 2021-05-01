package utils

import (
	"encoding/json"
	"net/http"
)

// Function ErrorCheck - Utility function that can be expanded to handle more error cases
// !! Only a suggestion. Just scrap if its a stupid idea !!
func ErrorCheck(w http.ResponseWriter, err error) bool {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}

// Function JsonResponse - Utility function that can be expanded
// !! Only a suggestion. Just scrap if its a stupid idea !!
func JsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	payload, err := json.Marshal(data)
	if ErrorCheck(w, err) {
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}
