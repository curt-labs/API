package search_ctlr

import (
	"net/http"

	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/search"
	"github.com/go-martini/martini"
)

func Search(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	terms := params["term"]

	res, err := search.Dsl(terms, []string{})
	if err != nil {
		apierror.GenerateError("Trouble searching", err, rw, r)
	}

	return encoding.Must(enc.Encode(res))
}
