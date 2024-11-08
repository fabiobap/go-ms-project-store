package helpers

import (
	"encoding/json"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	if code == http.StatusNoContent {
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
