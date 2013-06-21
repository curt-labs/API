package search_ctlr

import (
	"../../helpers/plate"
	. "../../models"
	"net/http"
	"strings"
)

func SearchPart(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	terms := params.Get(":term")
	key := params.Get("key")

	qry := PartSearchResult{
		Request: SearchQuery{
			SearchTerms: strings.Replace(terms, ",", " ", -1),
			StartIndex:  0,
			Count:       0,
		},
	}

	qry.SearchParts(key)

	plate.ServeFormatted(w, r, qry)
}
