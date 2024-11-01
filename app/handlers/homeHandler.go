package handlers

import (
	"encoding/json"
	"net/http"
	"runtime"
)

func Home(w http.ResponseWriter, r *http.Request) {
	WriteResponse(w, http.StatusOK, "API powered by Go v"+runtime.Version())
}

func WriteResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
