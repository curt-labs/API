package search_ctlr

import (
	"github.com/curt-labs/API/helpers/apicontext"
	"net/http"
	"strconv"

	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/search"
	"github.com/go-martini/martini"
)

func Search(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	terms := params["term"]
	qs := r.URL.Query()
	page, _ := strconv.Atoi(qs.Get("page"))
	count, _ := strconv.Atoi(qs.Get("count"))
	brand, _ := strconv.Atoi(qs.Get("brand"))
	rawPartNumber := qs.Get("raw")

	res, err := search.Dsl(terms, page, count, brand, dtx, rawPartNumber)
	if err != nil {
		apierror.GenerateError("Trouble searching", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(res))
}

func SearchExactAndClose(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	terms := params["term"]
	qs := r.URL.Query()
	page, _ := strconv.Atoi(qs.Get("page"))
	count, _ := strconv.Atoi(qs.Get("count"))
	brand, _ := strconv.Atoi(qs.Get("brand"))

	res, err := search.ExactAndCloseDsl(terms, page, count, brand, dtx)
	if err != nil {
		apierror.GenerateError("Trouble searching", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(res))
}
