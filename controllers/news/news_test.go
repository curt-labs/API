package news_controller

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	"github.com/curt-labs/GoAPI/helpers/pagination"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/news"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestNews(t *testing.T) {
	var n news_model.News
	var err error
	dtx, err := apicontextmock.Mock()
	Convey("Test News", t, func() {
		//test create
		//Form values
		form := url.Values{"title": {"test"}, "lead": {"testLead"}}
		v := form.Encode()
		body := strings.NewReader(v)
		testThatHttp.RequestWithDtx("post", "/news", "", "", Create, body, "application/x-www-form-urlencoded", dtx)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &n)
		So(n, ShouldHaveSameTypeAs, news_model.News{})

		//test update
		//Form values
		form = url.Values{"title": {"test new"}, "lead": {"testLead new"}}
		v = form.Encode()
		body = strings.NewReader(v)
		testThatHttp.RequestWithDtx("put", "/news/", ":id", strconv.Itoa(n.ID), Update, body, "application/x-www-form-urlencoded", dtx)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &n)
		So(n, ShouldHaveSameTypeAs, news_model.News{})

		//test get
		testThatHttp.Request("get", "/news/", ":id", strconv.Itoa(n.ID), Get, nil, "application/x-www-form-urlencoded")
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &n)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		So(n, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(n.Title, ShouldEqual, "test new")

		//test getall
		testThatHttp.RequestWithDtx("get", "/news", "", "", GetAll, nil, "application/x-www-form-urlencoded", dtx)
		var ns news_model.Newses
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &ns)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		So(len(ns), ShouldBeGreaterThan, 0)
		So(err, ShouldBeNil)

		//test getleads
		testThatHttp.RequestWithDtx("get", "/news/leads", "", "", GetLeads, nil, "application/x-www-form-urlencoded", dtx)
		var l pagination.Objects
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &l)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		So(len(l.Objects), ShouldBeGreaterThan, 0)
		So(err, ShouldBeNil)

		//test gettitles
		testThatHttp.RequestWithDtx("get", "/news/titles", "", "", GetTitles, nil, "application/x-www-form-urlencoded", dtx)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &l)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		So(len(l.Objects), ShouldBeGreaterThan, 0)
		So(err, ShouldBeNil)

		//test search
		testThatHttp.RequestWithDtx("get", "/news/search", "", "?title=test", Search, nil, "application/x-www-form-urlencoded", dtx)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &l)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		So(len(l.Objects), ShouldBeGreaterThan, 0)
		So(err, ShouldBeNil)

		//test delete
		testThatHttp.Request("delete", "/news/", ":id", strconv.Itoa(n.ID), Delete, nil, "application/x-www-form-urlencoded")
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &n)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		So(n, ShouldNotBeNil)
		So(err, ShouldBeNil)

	})
	_ = apicontextmock.DeMock(dtx)
}
