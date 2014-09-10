package search_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/search"
	"github.com/go-martini/martini"
	"net/http"
)

func Search(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	terms := params["term"]

	res, err := search.Dsl(terms, []string{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(res))
}
