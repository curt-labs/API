package vehicle

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/products"
	"github.com/go-martini/martini"

	"net/http"
)

func GetAllCollectionApplications(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext, params martini.Params) string {
	collection := params["collection"]
	if collection == "" {
		apierror.GenerateError("No Collection in URL", nil, w, r)
		return ""
	}
	apps, err := products.GetAllCollectionApplications(collection)
	if err != nil {
		apierror.GenerateError(err.Error(), err, w, r)
		return ""
	}
	return encoding.Must(enc.Encode(apps))
}
