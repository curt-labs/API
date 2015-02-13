package videos_ctlr

import (
	"encoding/json"
	"io/ioutil"
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

func GetPartVideos(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	// var p products.Part
	var prodId int

	var videos video.Videos
	if prodId, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting part ID", err, w, r)
		return ""
	}

	if videos, err = video.GetPartVideos(prodId); err != nil {
		apierror.GenerateError("Trouble getting part videos", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(videos))
}

func SaveVideo(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var v video.Video
	var err error

	if params["id"] != "" {
		if v.ID, err = strconv.Atoi(params["id"]); err != nil {
			apierror.GenerateError("Trouble getting video ID", err, w, r)
			return ""
		}
		if err = v.Get(); err != nil {
			apierror.GenerateError("Trouble getting video", err, w, r)
			return ""
		}
	}

	//json
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {

		apierror.GenerateError("Trouble reading request body while saving video", err, w, r)
		return encoding.Must(enc.Encode(false))
	}

	if err = json.Unmarshal(requestBody, &v); err != nil {
		apierror.GenerateError("Trouble unmarshalling json request body while saving video", err, w, r)
		return encoding.Must(enc.Encode(false))
	}

	//create or update
	if v.ID > 0 {
		err = v.Update(dtx)
	} else {
		err = v.Create(dtx)
	}

	if err != nil {
		msg := "Trouble creating video"
		if v.ID > 0 {
			msg = "Trouble updating video"
		}
		apierror.GenerateError(msg, err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(v))
}
func DeleteVideo(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var v video.Video
	var err error

	if v.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting video ID", err, w, r)
		return ""
	}

	if err = v.Delete(dtx); err != nil {
		apierror.GenerateError("Trouble deleting video", err, w, r)

		return ""
	}

	return encoding.Must(enc.Encode(v))
}

func SaveChannel(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c video.Channel
	var err error

	if params["id"] != "" {
		if c.ID, err = strconv.Atoi(params["id"]); err != nil {
			apierror.GenerateError("Trouble getting video channel ID", err, w, r)
			return ""
		}
		if err = c.Get(); err != nil {
			apierror.GenerateError("Trouble getting video channel", err, w, r)

			return ""
		}
	}

	//json
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {

		apierror.GenerateError("Trouble reading request body while saving video channel", err, w, r)
		return encoding.Must(enc.Encode(false))
	}

	if err = json.Unmarshal(requestBody, &c); err != nil {
		apierror.GenerateError("Trouble unmarshalling json request body while saving video channel", err, w, r)
		return encoding.Must(enc.Encode(false))

	}

	//create or update
	if c.ID > 0 {
		err = c.Update()
	} else {
		err = c.Create()
	}

	if err != nil {

		msg := "Trouble creating video channel"
		if c.ID > 0 {
			msg = "Trouble updating video channel"
		}
		apierror.GenerateError(msg, err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

func DeleteChannel(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c video.Channel
	var err error

	if c.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting video channel ID", err, w, r)
		return ""
	}

	if err = c.Delete(); err != nil {
		apierror.GenerateError("Trouble deleting video channel", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

func SaveCdn(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c video.CdnFile
	var err error

	if params["id"] != "" {
		if c.ID, err = strconv.Atoi(params["id"]); err != nil {
			apierror.GenerateError("Trouble getting video CDN ID", err, w, r)
			return ""
		}
		if err = c.Get(); err != nil {
			apierror.GenerateError("Trouble getting video CDN", err, w, r)
			return ""
		}
	}

	//json
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apierror.GenerateError("Trouble reading request body while saving video CDN", err, w, r)
		return encoding.Must(enc.Encode(false))
	}

	if err = json.Unmarshal(requestBody, &c); err != nil {
		apierror.GenerateError("Trouble unmarshalling json request body while saving video CDN", err, w, r)
		return encoding.Must(enc.Encode(false))
	}

	//create or update
	if c.ID > 0 {
		err = c.Update()
	} else {
		err = c.Create()
	}

	if err != nil {

		msg := "Trouble creating video CDN"
		if c.ID > 0 {
			msg = "Trouble updating video CDN"
		}
		apierror.GenerateError(msg, err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

func DeleteCdn(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c video.CdnFile
	var err error

	if c.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting video CDN ID", err, w, r)
		return ""
	}

	if err = c.Delete(); err != nil {
		apierror.GenerateError("Trouble deleting video CDN", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

func SaveVideoType(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c video.VideoType
	var err error

	if params["id"] != "" {
		if c.ID, err = strconv.Atoi(params["id"]); err != nil {
			apierror.GenerateError("Trouble getting video type ID", err, w, r)
			return ""
		}
		if err = c.Get(); err != nil {
			apierror.GenerateError("Trouble getting video type", err, w, r)
			return ""
		}
	}

	//json
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {

		apierror.GenerateError("Trouble reading request body while saving video type", err, w, r)
		return encoding.Must(enc.Encode(false))
	}

	if err = json.Unmarshal(requestBody, &c); err != nil {
		apierror.GenerateError("Trouble unmarshalling json request body while saving video type", err, w, r)
		return encoding.Must(enc.Encode(false))
	}

	//create or update
	if c.ID > 0 {
		err = c.Update()
	} else {
		err = c.Create()
	}

	if err != nil {

		msg := "Trouble creating video type"
		if c.ID > 0 {
			msg = "Trouble updating video type"
		}
		apierror.GenerateError(msg, err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}
func DeleteVideoType(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c video.VideoType
	var err error

	if c.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting video type ID", err, w, r)
		return ""
	}

	if err = c.Delete(); err != nil {
		apierror.GenerateError("Trouble deleting video type", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

func GetCdnType(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var v video.CdnFileType
	var err error

	if v.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting video CDN type ID", err, w, r)
		return ""
	}

	if err = v.Delete(); err != nil {
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

func SaveCdnType(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c video.CdnFileType
	var err error

	if params["id"] != "" {
		if c.ID, err = strconv.Atoi(params["id"]); err != nil {
			apierror.GenerateError("Trouble getting video CDN type ID", err, w, r)
			return ""
		}
		if err = c.Get(); err != nil {
			apierror.GenerateError("Trouble getting video CDN type", err, w, r)
			return ""
		}
	}

	//json
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {

		apierror.GenerateError("Trouble reading request body while saving video CDN type", err, w, r)
		return encoding.Must(enc.Encode(false))
	}

	if err = json.Unmarshal(requestBody, &c); err != nil {
		apierror.GenerateError("Trouble unmarshalling json request body while saving video CDN type", err, w, r)
		return encoding.Must(enc.Encode(false))
	}

	//create or update
	if c.ID > 0 {
		err = c.Update()
	} else {
		err = c.Create()
	}

	if err != nil {

		msg := "Trouble creating video CDN type"
		if c.ID > 0 {
			msg = "Trouble updating video CDN type"
		}
		apierror.GenerateError(msg, err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

func DeleteCdnType(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c video.CdnFileType
	var err error

	if c.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting video CDN type ID", err, w, r)
		return ""
	}

	if err = c.Delete(); err != nil {
		apierror.GenerateError("Trouble deleting video CDN type", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(c))
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

func SaveChannelType(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c video.ChannelType
	var err error

	if params["id"] != "" {
		if c.ID, err = strconv.Atoi(params["id"]); err != nil {
			apierror.GenerateError("Trouble getting video channel type ID", err, w, r)
			return ""
		}
		if err = c.Get(); err != nil {
			apierror.GenerateError("Trouble getting video channel type", err, w, r)
			return ""
		}
	}

	//json
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {

		apierror.GenerateError("Trouble reading request body while saving video channel type", err, w, r)
		return encoding.Must(enc.Encode(false))
	}

	if err = json.Unmarshal(requestBody, &c); err != nil {
		apierror.GenerateError("Trouble unmarshalling json request body while saving video channel type", err, w, r)
		return encoding.Must(enc.Encode(false))
	}

	//create or update
	if c.ID > 0 {
		err = c.Update()
	} else {
		err = c.Create()
	}

	if err != nil {

		msg := "Trouble creating video channel type"
		if c.ID > 0 {
			msg = "Trouble updating video channel type"
		}
		apierror.GenerateError(msg, err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

func DeleteChannelType(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c video.ChannelType
	var err error

	if c.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting video channel type ID", err, w, r)
		return ""
	}

	if err = c.Delete(); err != nil {
		apierror.GenerateError("Trouble deleting video channel type", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}
