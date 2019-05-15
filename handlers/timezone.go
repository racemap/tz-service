package handlers

import (
	"fmt"
	"net/http"

	timezone "github.com/evanoberholster/timezoneLookup"
	"github.com/gorilla/mux"
)

func TimezoneHandler(tz timezone.TimezoneInterface) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Category: %v\n", vars["category"])
	}
}
