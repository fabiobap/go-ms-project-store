package app

import (
	"encoding/json"
	"net/http"
	"runtime"
)

func Home(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, http.StatusOK, "API powered by Go v"+runtime.Version())
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
