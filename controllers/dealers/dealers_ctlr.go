package dealers_ctlr

import (
	. "../../models"
	"../../plate"
	"net/http"
)

func Etailers(w http.ResponseWriter, r *http.Request) {

	dealers, err := GetEtailers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	plate.ServeFormatted(w, r, dealers)
}
