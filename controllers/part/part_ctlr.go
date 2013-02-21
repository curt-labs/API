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
	part := Part{
		PartId: id,
	}

	err := part.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusFound)
	}

	plate.ServeFormatted(w, r, part)
}
