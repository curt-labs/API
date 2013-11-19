package videos_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/plate"
	. "github.com/curt-labs/GoAPI/models"
	"net/http"
)

func DistinctVideos(w http.ResponseWriter, r *http.Request) {

	videos, err := UniqueVideos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, videos)

}
