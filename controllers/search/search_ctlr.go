package search_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/search"
	"github.com/go-martini/martini"
	"net/http"
	"strings"
)

func SearchPart(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	qs := r.URL.Query()
	terms := params["term"]
	key := qs.Get("key")

	qry := search.PartSearchResult{
		Request: search.SearchQuery{
			SearchTerms: strings.Replace(terms, ",", " ", -1),
			StartIndex:  0,
			Count:       0,
		},
	}

	qry.SearchParts(key)

	return encoding.Must(enc.Encode(qry))
}
