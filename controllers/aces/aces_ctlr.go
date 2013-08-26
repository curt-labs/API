package aces_ctlr

import (
	"../../helpers/plate"
	"../../models/aces"
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
