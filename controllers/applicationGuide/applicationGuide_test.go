package applicationGuide

import (
	"bytes"
	"encoding/json"
	"github.com/curt-labs/API/helpers/testThatHttp"
	"github.com/curt-labs/API/models/applicationGuide"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestApplicationGuide(t *testing.T) {
	var a applicationGuide.ApplicationGuide
	var err error
	Convey("Testing ApplicationGuide", t, func() {
		//test create
		form := url.Values{"url": {"test"}, "fileType": {"pdf"}, "website_id": {"1"}}
		v := form.Encode()
		body := strings.NewReader(v)
		testThatHttp.Request("post", "/applicationGuide", "", "", CreateApplicationGuide, body, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &a)
		So(err, ShouldBeNil)
		So(a, ShouldHaveSameTypeAs, applicationGuide.ApplicationGuide{})

		//test create using json
		var jsonAppGuide applicationGuide.ApplicationGuide
		jsonAppGuide.Url = "www.www.com"
		jsonAppGuide.FileType = "pdf"
		jsonAppGuide.Website.ID = 1
		bodyBytes, _ := json.Marshal(jsonAppGuide)
		bodyJson := bytes.NewReader(bodyBytes)
		testThatHttp.Request("post", "/applicationGuide", "", "", CreateApplicationGuide, bodyJson, "application/json")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &jsonAppGuide)
		So(err, ShouldBeNil)
		So(jsonAppGuide, ShouldHaveSameTypeAs, applicationGuide.ApplicationGuide{})

		//test get
		testThatHttp.Request("get", "/applicationGuide/", ":id", strconv.Itoa(a.ID), GetApplicationGuide, nil, "")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &a)
		So(err, ShouldBeNil)
		So(a, ShouldHaveSameTypeAs, applicationGuide.ApplicationGuide{})
		So(a.Url, ShouldEqual, "test")

		//test get by website
		testThatHttp.Request("get", "/applicationGuide/website/", ":id", strconv.Itoa(a.Website.ID), GetApplicationGuidesByWebsite, nil, "")
		var as []applicationGuide.ApplicationGuide
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &as)
		So(err, ShouldBeNil)
		So(len(as), ShouldBeGreaterThanOrEqualTo, 0)
		So(a.Url, ShouldEqual, "test")

		//test delete
		testThatHttp.Request("delete", "/applicationGuide/", ":id", strconv.Itoa(a.ID), DeleteApplicationGuide, nil, "")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &a)
		So(err, ShouldBeNil)

		//cleanup
		err = jsonAppGuide.Delete()

	})
}

func BenchmarkApplicationGuide(b *testing.B) {
	testThatHttp.RequestBenchmark(b.N, "GET", "/applicationGuide/1", nil, GetApplicationGuide)
	testThatHttp.RequestBenchmark(b.N, "GET", "/applicationGuide/1", nil, GetApplicationGuidesByWebsite)
}
