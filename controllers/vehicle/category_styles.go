package vehicle

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"strconv"
	"strings"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/products"
	"github.com/go-martini/martini"
)

var (
	// DefaultStatuses Normal statuses used to query products
	DefaultStatuses = []int{800, 900}
)

func QueryCategoryStyle(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	if err := database.Init(); err != nil {
		apierror.GenerateError("Trouble generating database connection", err, w, r)
		return ""
	}

	session := database.ProductMongoSession.Copy()
	defer session.Close()

	statuses := DefaultStatuses
	if r.URL.Query().Get("brands") != "" {
		segs := strings.Split(r.URL.Query().Get("brands"), ",")
		var ids []int
		for _, seg := range segs {
			id, err := strconv.Atoi(seg)
			if err == nil {
				ids = append(ids, id)
			}
		}
		statuses = ids
	}

	ctx := &products.LookupContext{
		Brands:   dtx.BrandArray,
		Statuses: statuses,
		Session:  session,
	}

	cats, err := products.Query(
		ctx,
		params["year"],
		params["make"],
		params["model"],
		params["category"],
	)
	if err != nil {
		apierror.GenerateError("Trouble getting part", err, w, r)
		return ""
	}

	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		return encoding.Must(enc.Encode(cats))
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	resp, err := enc.Encode(cats)
	if err != nil {
		apierror.GenerateError("Failed to encode", err, w, r)
		return ""
	}

	respBytes := []byte(resp)
	w.Header().Set("Content-Encoding", "gzip")
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", http.DetectContentType(respBytes))
	}

	gz.Write(respBytes)
	gz.Close()

	return buf.String()
}
