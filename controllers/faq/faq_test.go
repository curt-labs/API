package faq_controller

//reference
//https://groups.google.com/forum/#!topic/golang-nuts/DARY7HY-pbY
//http://blog.wercker.com/2014/02/06/RethinkDB-Gingko-Martini-Golang.html
//http://golang.org/pkg/net/http/httptest/#example_ResponseRecorder
//https://github.com/mies/martini-rethink

import (
	"encoding/json"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/faq"
	"github.com/go-martini/martini"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// func TestFaq(t *testing.T) {
// 	Convey("", t, func() {
// 		m := martini.Classic()
// 		m.Get("/faqs", GetAll)
// 		m.Use(render.Renderer())
// 		m.Use(encoding.MapEncoder)

// 		request, _ := http.NewRequest("GET", "/faqs", nil)
// 		response := httptest.NewRecorder()
// 		m.ServeHTTP(response, request)
// 		// t.Log(response.Code)
// 		// t.Log(response.Body)
// 	})
// }

func TestGetAllFaq(t *testing.T) {
	Convey("Faq Getall", t, func() {
		m := martini.Classic()
		m.Get("/faqs", GetAll)
		m.Use(render.Renderer())
		m.Use(encoding.MapEncoder)

		request, _ := http.NewRequest("GET", "/faqs", nil)
		response := httptest.NewRecorder()
		m.ServeHTTP(response, request)
		So(response.Code, ShouldEqual, 200)

		body := response.Body.Bytes()
		var fs []faq_model.Faq

		err := json.Unmarshal(body, &fs)
		So(err, ShouldBeNil)

	})
}

func TestGetFaq(t *testing.T) {
	Convey("Faq Get", t, func() {
		m := martini.Classic()
		m.Get("/faqs/:id", Get)
		m.Use(render.Renderer())
		m.Use(encoding.MapEncoder)

		request, _ := http.NewRequest("GET", "/faqs/1", nil)
		response := httptest.NewRecorder()
		m.ServeHTTP(response, request)
		So(response.Code, ShouldEqual, 200)
		// t.Log(response.Code)
		// t.Log(response.Body)

		body := response.Body.Bytes()
		var f faq_model.Faq

		err := json.Unmarshal(body, &f)
		So(err, ShouldBeNil)

	})
}

//TODO - set form values right
func TestSearch(t *testing.T) {
	Convey("Faq Search", t, func() {
		m := martini.Classic()
		m.Get("/faqs/search", Search)
		m.Use(render.Renderer())
		m.Use(encoding.MapEncoder)

		//Form values
		form := url.Values{"question": {"test2"}, "answer": {"testan2"}}
		v := form.Encode()
		s := strings.NewReader(v)

		request, _ := http.NewRequest("GET", "/faqs/search", s)

		response := httptest.NewRecorder()
		m.ServeHTTP(response, request)
		So(response.Code, ShouldEqual, 200)
		// t.Log(response.Code)
		// t.Log(response.Body)

		body := response.Body.Bytes()
		var f faq_model.Faq

		err := json.Unmarshal(body, &f)
		So(err, ShouldBeNil)
	})
}

//TODO - set form values right
func TestCreate(t *testing.T) {
	Convey("Faq Search", t, func() {
		m := martini.Classic()
		m.Put("/faqs", Create)
		m.Use(render.Renderer())
		m.Use(encoding.MapEncoder)

		//Form values
		form := url.Values{"question": {"test2"}, "answer": {"testan2"}}
		v := form.Encode()
		s := strings.NewReader(v)

		request, _ := http.NewRequest("PUT", "/faqs", s)

		response := httptest.NewRecorder()
		m.ServeHTTP(response, request)
		So(response.Code, ShouldEqual, 200)
		t.Log(response.Code)
		t.Log(response.Body)

		body := response.Body.Bytes()
		var f faq_model.Faq

		err := json.Unmarshal(body, &f)
		So(err, ShouldBeNil)

		f.Delete()
	})
}
