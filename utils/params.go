package utils

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetParamID(r *http.Request) uint {
	vars := mux.Vars(r)
	idStr := vars["id"]

	idUint64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0 // handle safely if needed
	}
	return uint(idUint64)
}
