package aces_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/aces"
	"net/http"
)

func ACES(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {

	resp, err := aces.GetACESPartData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(resp))
}
