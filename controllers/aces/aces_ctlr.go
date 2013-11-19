package aces_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/plate"
	"github.com/curt-labs/GoAPI/models/aces"
	"net/http"
)

func ACES(w http.ResponseWriter, r *http.Request) {

	resp, err := aces.GetACESPartData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeStringAsXml(w, resp)
}
