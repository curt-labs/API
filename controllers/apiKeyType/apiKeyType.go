package apiKeyType

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/apiKeyType"

	"net/http"
)

func GetApiKeyTypes(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	types, err := apiKeyType.GetAllApiKeyTypes()
	if err != nil {
		if err != nil {
			apierror.GenerateError("Trouble converting ID parameter", err, rw, req)
		}
	}
	return encoding.Must(enc.Encode(types))
}
