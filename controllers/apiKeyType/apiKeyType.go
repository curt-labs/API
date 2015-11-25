package apiKeyType

import (
	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/apiKeyType"

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
