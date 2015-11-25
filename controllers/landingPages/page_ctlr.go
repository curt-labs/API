package landingPage

import (
	//"encoding/json"
	"net/http"
	"strconv"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/landingPages"
	"github.com/go-martini/martini"
)

func Get(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var lp landingPage.LandingPage
	var err error
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		apierror.GenerateError("Must provide a Landing Page ID", err, rw, req)
	}

	lp.Id = id
	err = lp.Get(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting Landing Page by Id.", err, rw, req)
	}
	return encoding.Must(enc.Encode(lp))
}
