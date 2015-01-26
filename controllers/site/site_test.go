package site

import (
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/models/site"
	. "github.com/smartystreets/goconvey/convey"

	"encoding/json"
	"net/url"
	"strconv"
	"testing"
)

func TestSite(t *testing.T) {

	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}

	Convey("Site Contents", t, func() {
		var c site.Content
		var contents site.Contents
		var cr site.ContentRevision
		c.WebsiteId = 1
		qs := make(url.Values, 0)
		qs.Add("key", dtx.APIKey)

		c.Title = "test content - controller test"
		cr.Text = "test revision - controller test"
		c.ContentRevisions = append(c.ContentRevisions, cr)

		response := httprunner.ParameterizedJsonRequest("POST", "/site/content", "/site/content", &qs, c, SaveContent)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &c), ShouldEqual, nil)

		c.Active = true
		response = httprunner.ParameterizedJsonRequest("PUT", "/site/content/:id", "/site/content/"+strconv.Itoa(c.Id), &qs, c, SaveContent)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &c), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/site/content/:id", "/site/content/"+strconv.Itoa(c.Id), &qs, nil, GetContent)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &c), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/site/content", "/site/content", &qs, nil, GetAllContents)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &contents), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/site/content/:id/revisions", "/site/content/"+strconv.Itoa(c.Id)+"/revisions", &qs, nil, GetContentRevisions)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &c), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("DELETE", "/site/content/:id", "/site/content/"+strconv.Itoa(c.Id), &qs, nil, DeleteContent)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &c), ShouldEqual, nil)

	})

	Convey("Site Menu", t, func() {
		var m site.Menu
		var menus site.Menus
		qs := make(url.Values, 0)
		qs.Add("key", dtx.APIKey)

		m.Name = "test menu - controller test"

		response := httprunner.ParameterizedJsonRequest("POST", "/site/menu", "/site/menu", &qs, m, SaveMenu)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &m), ShouldEqual, nil)

		m.Active = true
		response = httprunner.ParameterizedJsonRequest("PUT", "/site/menu/:id", "/site/menu/"+strconv.Itoa(m.Id), &qs, m, SaveMenu)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &m), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/site/menu/:id", "/site/menu/"+strconv.Itoa(m.Id), &qs, nil, GetMenu)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &m), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/site/menu", "/site/menu", &qs, nil, GetAllMenus)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &menus), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/site/menu/contents/:id", "/site/menu/contents/"+strconv.Itoa(m.Id), &qs, nil, GetMenuWithContents)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &m), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("DELETE", "/site/menu/:id", "/site/menu/"+strconv.Itoa(m.Id), &qs, nil, DeleteMenu)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &m), ShouldEqual, nil)

	})

	Convey("Site Menu", t, func() {
		var s site.Website
		qs := make(url.Values, 0)
		qs.Add("key", dtx.APIKey)

		s.Description = "test site - controller test"

		response := httprunner.ParameterizedJsonRequest("POST", "/site", "/site", &qs, s, SaveSite)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &s), ShouldEqual, nil)

		s.Url = "www.controllertest.com"

		response = httprunner.ParameterizedJsonRequest("PUT", "/site/:id", "/site/"+strconv.Itoa(s.ID), &qs, s, SaveSite)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &s), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/site/:id/details", "/site/"+strconv.Itoa(s.ID)+"/details", &qs, nil, GetSiteDetails)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &s), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("DELETE", "/site/:id", "/site/"+strconv.Itoa(s.ID), &qs, nil, DeleteSite)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &s), ShouldEqual, nil)

	})

	_ = apicontextmock.DeMock(dtx)
}
