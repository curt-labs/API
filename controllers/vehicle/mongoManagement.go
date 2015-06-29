package vehicle

// import (
// 	"github.com/curt-labs/GoAPI/helpers/apicontext"
// 	"github.com/curt-labs/GoAPI/helpers/encoding"
// 	"github.com/curt-labs/GoAPI/helpers/error"
// 	"github.com/curt-labs/GoAPI/models/products"
// 	"strconv"

// 	"net/http"
// )

// func Collections(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
// 	cols, err := products.GetAriesVehicleCollections()
// 	if err != nil {
// 		apierror.GenerateError(err.Error(), err, w, r)
// 		return ""
// 	}

// 	return encoding.Must(enc.Encode(cols))
// }
