package videos_ctlr

import (
	"../../helpers/plate"
	. "../../models"
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
