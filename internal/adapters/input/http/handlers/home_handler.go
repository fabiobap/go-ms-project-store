package handlers

import (
	"net/http"
	"runtime"

	"github.com/go-ms-project-store/internal/pkg/helpers"
)

func Home(w http.ResponseWriter, r *http.Request) {
	helpers.WriteResponse(w, http.StatusOK, "API powered by Go v"+runtime.Version())
}
