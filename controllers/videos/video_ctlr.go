package videos_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/video"
	"github.com/go-martini/martini"
	// "log"
	"net/http"
	"strconv"
)

//gets old videos
func DistinctVideos(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {

	videos, err := video.UniqueVideos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(videos))
}

// New videos, literally from the "video_new" table
func Get(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var v video.Video
	var err error

	id, err := strconv.Atoi(params["id"])
	v.ID = id
	err = v.Get()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(v))
}

func GetAllVideos(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error

	vs, err := video.GetAllVideos()
	// log.Print(vs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(vs))
}

func GetChannel(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var ch video.Channel
	var err error

	id, err := strconv.Atoi(params["id"])
	ch.ID = id
	err = ch.Get()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(ch))
}

func GetAllChannels(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error

	cs, err := video.GetAllChannels()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(cs))
}
func GetCdn(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var cdn video.CdnFile
	var err error

	id, err := strconv.Atoi(params["id"])
	cdn.ID = id
	err = cdn.Get()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(cdn))
}

func GetAllCdns(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error

	cdns, err := video.GetAllCdnFiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(cdns))
}
func GetType(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var vt video.VideoType
	var err error

	id, err := strconv.Atoi(params["id"])
	vt.ID = id
	err = vt.Get()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(vt))
}

func GetAllTypes(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error

	vts, err := video.GetAllVideoTypes()
	// log.Print(vs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(vts))
}
