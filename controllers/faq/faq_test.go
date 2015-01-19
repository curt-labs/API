package faq_controller

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/helpers/pagination"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/faq"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestFaqs(t *testing.T) {
	var f faq_model.Faq
	var err error
	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}
	Convey("Test Faqs", t, func() {
		//test create
		form := url.Values{"question": {"test"}, "answer": {"testAnswer"}}
		v := form.Encode()
		body := strings.NewReader(v)
		testThatHttp.Request("post", "/faqs", "", "", Create, body, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &f)
		So(f, ShouldHaveSameTypeAs, faq_model.Faq{})

		//test update
		form = url.Values{"question": {"test new"}, "answer": {"testAnswer new"}}
		v = form.Encode()
		body = strings.NewReader(v)
		testThatHttp.RequestWithDtx("put", "/faqs/", ":id", strconv.Itoa(f.ID), Update, body, "application/x-www-form-urlencoded", dtx)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &f)
		So(f, ShouldHaveSameTypeAs, faq_model.Faq{})

		//test get
		testThatHttp.RequestWithDtx("get", "/faqs/", ":id", strconv.Itoa(f.ID), Get, nil, "application/x-www-form-urlencoded", dtx)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &f)
		So(f, ShouldHaveSameTypeAs, faq_model.Faq{})
		So(f.Question, ShouldEqual, "test new")

		//test getall
		testThatHttp.Request("get", "/faqs", "", "", GetAll, nil, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var fs faq_model.Faqs
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &fs)
		So(len(fs), ShouldBeGreaterThanOrEqualTo, 0)

		//test search - responds w/ horrid pagination object
		testThatHttp.Request("get", "/faqs/search", "", "?question=test new", Search, nil, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var l pagination.Objects
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &l)
		So(len(l.Objects), ShouldBeGreaterThanOrEqualTo, 0)

		//test delete
		testThatHttp.Request("delete", "/faqs/", ":id", strconv.Itoa(f.ID), Delete, nil, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &f)
		So(f, ShouldHaveSameTypeAs, faq_model.Faq{})
	})
	_ = apicontextmock.DeMock(dtx)
}

func BenchmarkCRUDFaqs(b *testing.B) {
	qs := make(url.Values, 0)

	form := url.Values{"question": {"test"}, "answer": {"testAnswer"}}

	Convey("Faqs", b, func() {
		//create faqs
		(&httprunner.BenchmarkOptions{
			Method:             "POST",
			Route:              "/faqs",
			ParameterizedRoute: "/faqs",
			Handler:            Create,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           form,
			Runs:               b.N,
		}).RequestBenchmark()
		//get all
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/faqs",
			ParameterizedRoute: "/faqs",
			Handler:            Get,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
		//get
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/faqs",
			ParameterizedRoute: "/faqs/1",
			Handler:            Get,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
		//delete
		(&httprunner.BenchmarkOptions{
			Method:             "DELETE",
			Route:              "/faqs",
			ParameterizedRoute: "/faqs",
			Handler:            Delete,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()

	})
}
