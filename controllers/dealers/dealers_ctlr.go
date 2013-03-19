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

// Sample Data
// 
// latlng: 43.853282,-95.571675,45.800981,-90.468526
// center: 44.83536,-93.0201
// 
// Old Path: http://curtmfg.com/WhereToBuy/getLocalDealersJSON?latlng=43.853282,-95.571675,45.800981,-90.468526&center=44.83536,-93.0201
func LocalDealers(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	latlng := params.Get("latlng")
	center := params.Get("center")

	dealers, err := GetLocalDealers(center, latlng)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, dealers)
}
