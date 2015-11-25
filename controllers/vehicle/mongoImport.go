package vehicle

import (
	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/products"

	// "log"
	"net/http"
	"strings"
)

type ErrorResp struct {
	ConversionErrs []error `json:"conversion_errors" xml:"conversion_errors"`
	InsertErrs     []error `json:"insert_errors" xml:"insert_errors"`
}

//requires the "Consolidated App Guides" that MJ produces in Excel
//intended to be a short term solution until Aries-Curt data merge is complete
//powers the Godzilla application

//Import a Csv
//Fields are expected to be: Part (oldpartnumber), Make, Model, Style, Year - 5 columns total
func ImportCsv(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	contentTypeHeader := r.Header.Get("Content-Type")
	contentTypeArr := strings.Split(contentTypeHeader, ";")
	if len(contentTypeArr) < 1 {
		apierror.GenerateError("Content-Type is not multipart/form-data", nil, w, r)
		return ""
	}
	contentType := contentTypeArr[0]
	if contentType != "multipart/form-data" {
		apierror.GenerateError("Content-Type is not multipart/form-data", nil, w, r)
		return ""
	}
	file, header, err := r.FormFile("file")

	if err != nil {
		apierror.GenerateError("Error getting file", err, w, r)
		return ""
	}
	defer file.Close()

	collectionName := header.Filename

	conversionErrs, insertErrs, err := products.Import(file, collectionName)
	if err != nil {
		apierror.GenerateError("Error importing", err, w, r)
		return ""
	}

	errResp := ErrorResp{
		ConversionErrs: conversionErrs,
		InsertErrs:     insertErrs,
	}

	return encoding.Must(enc.Encode(errResp))
}
