package part_ctlr

import (
	. "../../models"
	"../../plate"
	"net/http"
	"strconv"
)

func Get(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, _ := strconv.Atoi(params.Get(":part"))
	key := params.Get("key")
	part := Part{
		PartId: id,
	}

	err := part.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, part)
}

func Vehicles(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, err := strconv.Atoi(params.Get(":part"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusFound)
	}

	vehicles, err := ReverseLookup(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, vehicles)
}
