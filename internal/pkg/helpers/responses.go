package helpers

import (
	"encoding/json"
	"net/http"
	"time"
)

func DatetimeToString(field time.Time) string {
	return field.Format("2006-01-02 15:04:05")
}

func WriteResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
