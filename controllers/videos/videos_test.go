package videos_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/models/brand"
	"github.com/curt-labs/GoAPI/models/products"
	"github.com/curt-labs/GoAPI/models/video"
	. "github.com/smartystreets/goconvey/convey"

	"encoding/json"
	"net/url"
	"strconv"
	"testing"
)

func TestVideos(t *testing.T) {

	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}

	Convey("Video Channel Types", t, func() {
		var ct video.ChannelType
		var cts []video.ChannelType

		qs := make(url.Values, 0)
		qs.Add("key", dtx.APIKey)

		ct.Name = "controller test type"

		response := httprunner.ParameterizedJsonRequest("POST", "/channel/type", "/channel/type", &qs, ct, SaveChannelType)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &ct), ShouldEqual, nil)

		ct.Description = "test Desc"

		response = httprunner.ParameterizedJsonRequest("POST", "/channel/type/:id", "/channel/type/"+strconv.Itoa(ct.ID), &qs, ct, SaveChannelType)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &ct), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/channel/type/:id", "/channel/type/"+strconv.Itoa(ct.ID), &qs, nil, GetChannelType)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &ct), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/channel/type", "/channel/type", &qs, nil, GetAllChannelTypes)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &cts), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("DELETE", "/channel/type/:id", "/channel/type/"+strconv.Itoa(ct.ID), &qs, nil, DeleteChannelType)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &ct), ShouldEqual, nil)
	})

	Convey("Video Channels", t, func() {
		var c video.Channel
		var cs video.Channels

		qs := make(url.Values, 0)
		qs.Add("key", dtx.APIKey)

		c.Title = "controller test channel"

		response := httprunner.ParameterizedJsonRequest("POST", "/channel", "/channel", &qs, c, SaveChannel)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &c), ShouldEqual, nil)

		c.Description = "test description"

		response = httprunner.ParameterizedJsonRequest("POST", "/channel/:id", "/channel/"+strconv.Itoa(c.ID), &qs, c, SaveChannel)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &c), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/channel/:id", "/channel/"+strconv.Itoa(c.ID), &qs, nil, GetChannel)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &c), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/channel", "/channel", &qs, nil, GetAllChannels)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &cs), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("DELETE", "/channel/:id", "/channel/"+strconv.Itoa(c.ID), &qs, nil, DeleteChannel)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &c), ShouldEqual, nil)
	})

	Convey("Video Cdn Type", t, func() {
		var ct video.CdnFileType
		var cts []video.CdnFileType

		qs := make(url.Values, 0)
		qs.Add("key", dtx.APIKey)

		ct.Title = "controller test cdn type title"

		response := httprunner.ParameterizedJsonRequest("POST", "/cdn/type", "/cdn/type", &qs, ct, SaveCdnType)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &ct), ShouldEqual, nil)

		ct.Description = "test description"

		response = httprunner.ParameterizedJsonRequest("POST", "/cdn/type/:id", "/cdn/type/"+strconv.Itoa(ct.ID), &qs, ct, SaveCdnType)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &ct), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/cdn/type/:id", "/cdn/type/"+strconv.Itoa(ct.ID), &qs, nil, GetCdnType)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &ct), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/cdn/type", "/cdn/type", &qs, nil, GetAllCdnTypes)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &cts), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("DELETE", "/cdn/type/:id", "/cdn/type/"+strconv.Itoa(ct.ID), &qs, nil, DeleteCdnType)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &ct), ShouldEqual, nil)
	})

	Convey("Video Cdn", t, func() {
		var c video.CdnFile
		var cs video.CdnFiles

		qs := make(url.Values, 0)
		qs.Add("key", dtx.APIKey)

		c.ObjectName = "controller test cdn name"

		response := httprunner.ParameterizedJsonRequest("POST", "/cdn", "/cdn", &qs, c, SaveCdn)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &c), ShouldEqual, nil)

		c.Path = "test path"

		response = httprunner.ParameterizedJsonRequest("POST", "/cdn/:id", "/cdn/"+strconv.Itoa(c.ID), &qs, c, SaveCdn)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &c), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/cdn/:id", "/cdn/"+strconv.Itoa(c.ID), &qs, nil, GetCdn)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &c), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/cdn", "/cdn", &qs, nil, GetAllCdns)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &cs), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("DELETE", "/cdn/:id", "/cdn/"+strconv.Itoa(c.ID), &qs, nil, DeleteCdn)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &c), ShouldEqual, nil)
	})

	Convey("Video Types", t, func() {
		var vt video.VideoType
		var vts []video.VideoType

		qs := make(url.Values, 0)
		qs.Add("key", dtx.APIKey)

		vt.Name = "controller test video type"

		response := httprunner.ParameterizedJsonRequest("POST", "/type", "/type", &qs, vt, SaveVideoType)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &vt), ShouldEqual, nil)

		vt.Icon = "test icon"

		response = httprunner.ParameterizedJsonRequest("POST", "/type/:id", "/type/"+strconv.Itoa(vt.ID), &qs, vt, SaveVideoType)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &vt), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/type/:id", "/type/"+strconv.Itoa(vt.ID), &qs, nil, GetVideoType)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &vt), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/type", "/type", &qs, nil, GetAllVideoTypes)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &vts), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("DELETE", "/type/:id", "/type/"+strconv.Itoa(vt.ID), &qs, nil, DeleteVideoType)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &vt), ShouldEqual, nil)
	})

	Convey("Videos", t, func() {
		var v video.Video
		var vs video.Videos
		var p products.Part
		var b brand.Brand
		b.ID = dtx.BrandID
		p.ShortDesc = "test part"
		p.Create()

		qs := make(url.Values, 0)
		qs.Add("key", dtx.APIKey)

		v.Title = "controller test video title"
		v.Parts = append(v.Parts, p)
		v.Brands = append(v.Brands, b)

		response := httprunner.ParameterizedJsonRequest("POST", "", "", &qs, v, SaveVideo)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &v), ShouldEqual, nil)

		v.Thumbnail = "test thumbnail"

		response = httprunner.ParameterizedJsonRequest("POST", "/:id", "/"+strconv.Itoa(v.ID), &qs, v, SaveVideo)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &v), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/:id", "/"+strconv.Itoa(v.ID), &qs, nil, Get)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &v), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/distinct", "/distinct", &qs, nil, DistinctVideos)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &vs), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/part/:id", "/part/"+strconv.Itoa(p.ID), &qs, nil, GetPartVideos)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &vs), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/details/:id", "/details/"+strconv.Itoa(v.ID), &qs, nil, GetVideoDetails)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &v), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "", "", &qs, nil, GetAllVideos)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &vs), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("DELETE", "/:id", "/"+strconv.Itoa(v.ID), &qs, nil, DeleteVideo)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &v), ShouldEqual, nil)

		p.Delete()
	})

	_ = apicontextmock.DeMock(dtx)
}
