package videos_ctlr

import (
	"net/http"
	"strconv"

	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/video"
	"github.com/go-martini/martini"
)

//gets old videos
func DistinctVideos(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	videos, err := video.UniqueVideos(dtx)
	if err != nil {
		apierror.GenerateError("Touble getting distinct videos", err, w, r)
		return ""
	}
	return encoding.Must(enc.Encode(videos))
}

// New videos, literally from the "video_new" table
func Get(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var v video.Video
	var err error

	if v.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting video ID", err, w, r)
		return ""
	}

	if err = v.Get(); err != nil {
		apierror.GenerateError("Trouble getting video", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(v))
}

func GetVideoDetails(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var v video.Video
	var err error

	if v.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting video ID", err, w, r)
		return ""
	}

	if err = v.GetVideoDetails(); err != nil {
		apierror.GenerateError("Trouble getting video details", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(v))
}

func GetAllVideos(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	vs, err := video.GetAllVideos(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting all videos", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(vs))
}

func GetChannel(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var vchan video.Channel
	var err error

	if vchan.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting video channel ID", err, w, r)
		return ""
	}

	if err = vchan.Get(); err != nil {
		apierror.GenerateError("Trouble getting video channel", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(vchan))
}

func GetAllChannels(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {

	vchans, err := video.GetAllChannels()

	if err != nil {
		apierror.GenerateError("Trouble getting all video channels", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(vchans))
}

func GetCdn(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var cdn video.CdnFile
	var err error

	if cdn.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting video CDN ID", err, w, r)
		return ""
	}

	if err = cdn.Get(); err != nil {
		apierror.GenerateError("Trouble getting video cdn", err, w, r)

		return ""
	}

	return encoding.Must(enc.Encode(cdn))
}

func GetAllCdns(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	cdns, err := video.GetAllCdnFiles()

	if err != nil {
		apierror.GenerateError("Trouble getting all video CDNs", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(cdns))
}

func GetVideoType(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var vt video.VideoType
	var err error

	if vt.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting video type ID", err, w, r)
		return ""
	}

	if err = vt.Get(); err != nil {
		apierror.GenerateError("Trouble getting video type", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(vt))
}

func GetAllVideoTypes(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	vts, err := video.GetAllVideoTypes()

	if err != nil {
		apierror.GenerateError("Trouble getting all video types", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(vts))
}

func GetCdnType(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var v video.CdnFileType
	var err error

	if v.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting video CDN type ID", err, w, r)
		return ""
	}

	if err = v.Get(); err != nil {
		apierror.GenerateError("Trouble deleting video CDN type", err, w, r)

		return ""
	}

	return encoding.Must(enc.Encode(v))
}

func GetAllCdnTypes(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	ct, err := video.GetAllCdnFileTypes()

	if err != nil {

		apierror.GenerateError("Trouble getting all video CDN types", err, w, r)

		return ""
	}

	return encoding.Must(enc.Encode(ct))
}

func GetChannelType(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var v video.ChannelType
	var err error

	if v.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting video channel type ID", err, w, r)
		return ""
	}

	if err = v.Get(); err != nil {
		apierror.GenerateError("Trouble getting video channel type", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(v))
}

func GetAllChannelTypes(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	cs, err := video.GetAllChannelTypes()

	if err != nil {

		apierror.GenerateError("Trouble getting all video channel types", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(cs))
}
