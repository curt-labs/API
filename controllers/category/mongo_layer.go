package category_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/products"
	"github.com/go-martini/martini"

	"net/http"
	"strconv"
)

// GetCategory
func GetCategoryFromMongo(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var c products.Category
	var err error
	c.ID, err = strconv.Atoi(params["id"])
	if err != nil || c.ID == 0 {
		apierror.GenerateError("Trouble getting category identifier", err, rw, r)
		return ""
	}

	err = c.FromMongo()
	if err != nil || c.ID == 0 {
		apierror.GenerateError("Trouble getting category", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

func GetCategoryTree(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	cats, err := products.GetCategoryTree(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting categories", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(cats))
}

func GetCategoryParts(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	catIdStr := params["id"]
	catId, err := strconv.Atoi(catIdStr)
	if err != nil {
		apierror.GenerateError("Trouble getting category Id", err, rw, r)
		return ""
	}

	parts, err := products.GetCategoryParts(catId)
	if err != nil {
		apierror.GenerateError("Trouble getting parts", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(parts))
}
