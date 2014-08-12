package videos_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/video"
	"net/http"
)

func DistinctVideos(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {

	videos, err := video.UniqueVideos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(videos))
}
