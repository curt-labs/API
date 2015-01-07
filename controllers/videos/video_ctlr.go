package videos_ctlr

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/products"
	"github.com/curt-labs/GoAPI/models/video"
	"github.com/go-martini/martini"
	"io/ioutil"
	"net/http"
	"strconv"
)

//gets old videos
func DistinctVideos(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {

	videos, err := video.UniqueVideos(dtx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	if len(videos) == 0 {
		http.Error(w, "No videos.", http.StatusInternalServerError)
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

func GetVideoDetails(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var v video.Video
	var err error

	id, err := strconv.Atoi(params["id"])
	v.ID = id

	err = v.GetVideoDetails()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(v))
}
func GetAllVideos(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var err error

	vs, err := video.GetAllVideos(dtx)

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
func GetVideoType(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
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

func GetAllVideoTypes(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error

	vts, err := video.GetAllVideoTypes()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(vts))
}

func GetPartVideos(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var p products.Part
	id, err := strconv.Atoi(params["id"])
	p.ID = id
	videos, err := video.GetPartVideos(p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(videos))
}

func SaveVideo(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var v video.Video
	var err error
	idStr := params["id"]
	if idStr != "" {
		v.ID, err = strconv.Atoi(idStr)
		err = v.Get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return ""
		}
	}
	//json
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}
	err = json.Unmarshal(requestBody, &v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}
	//create or update
	if v.ID > 0 {
		err = v.Update()
	} else {
		err = v.Create()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(v))
}
func DeleteVideo(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var v video.Video
	var err error
	idStr := params["id"]

	v.ID, err = strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	err = v.Delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(v))

}

func SaveChannel(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c video.Channel
	var err error
	idStr := params["id"]
	if idStr != "" {
		c.ID, err = strconv.Atoi(idStr)
		err = c.Get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return ""
		}
	}
	//json
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}
	err = json.Unmarshal(requestBody, &c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}
	//create or update
	if c.ID > 0 {
		err = c.Update()
	} else {
		err = c.Create()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(c))
}
func DeleteChannel(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c video.Channel
	var err error
	idStr := params["id"]

	c.ID, err = strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	err = c.Delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(c))
}

func SaveCdn(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c video.CdnFile
	var err error
	idStr := params["id"]
	if idStr != "" {
		c.ID, err = strconv.Atoi(idStr)
		err = c.Get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return ""
		}
	}
	//json
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}
	err = json.Unmarshal(requestBody, &c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}
	//create or update
	if c.ID > 0 {
		err = c.Update()
	} else {
		err = c.Create()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(c))
}
func DeleteCdn(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c video.CdnFile
	var err error
	idStr := params["id"]

	c.ID, err = strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	err = c.Delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(c))
}

func SaveVideoType(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c video.VideoType
	var err error
	idStr := params["id"]
	if idStr != "" {
		c.ID, err = strconv.Atoi(idStr)
		err = c.Get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return ""
		}
	}
	//json
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}
	err = json.Unmarshal(requestBody, &c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}
	//create or update
	if c.ID > 0 {
		err = c.Update()
	} else {
		err = c.Create()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(c))
}
func DeleteVideoType(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c video.VideoType
	var err error
	idStr := params["id"]

	c.ID, err = strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	err = c.Delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(c))
}

func GetCdnType(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var v video.CdnFileType
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
func GetAllCdnTypes(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error

	ct, err := video.GetAllCdnFileTypes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(ct))
}
func SaveCdnType(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c video.CdnFileType
	var err error
	idStr := params["id"]
	if idStr != "" {
		c.ID, err = strconv.Atoi(idStr)
		err = c.Get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return ""
		}
	}
	//json
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}
	err = json.Unmarshal(requestBody, &c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}
	//create or update
	if c.ID > 0 {
		err = c.Update()
	} else {
		err = c.Create()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(c))
}
func DeleteCdnType(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c video.CdnFileType
	var err error
	idStr := params["id"]

	c.ID, err = strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	err = c.Delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(c))
}

func GetChannelType(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var v video.ChannelType
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
func GetAllChannelTypes(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error

	cs, err := video.GetAllChannelTypes()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(cs))
}
func SaveChannelType(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c video.ChannelType
	var err error
	idStr := params["id"]
	if idStr != "" {
		c.ID, err = strconv.Atoi(idStr)
		err = c.Get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return ""
		}
	}
	//json
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}
	err = json.Unmarshal(requestBody, &c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}
	//create or update
	if c.ID > 0 {
		err = c.Update()
	} else {
		err = c.Create()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(c))
}
func DeleteChannelType(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c video.ChannelType
	var err error
	idStr := params["id"]

	c.ID, err = strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	err = c.Delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(c))
}
