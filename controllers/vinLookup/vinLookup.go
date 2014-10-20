package vinLookup

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/vinLookup"
	"github.com/go-martini/martini"
	"log"
	"net/http"
)

func GetParts(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	vin := params["vin"]

	parts, err := vinLookup.VinPartLookup(vin)
	if err != nil {
		log.Print(err)
		return ""
	}

	return encoding.Must(enc.Encode(parts))

}
