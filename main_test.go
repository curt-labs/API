package main

//reference
//https://groups.google.com/forum/#!topic/golang-nuts/DARY7HY-pbY
//http://blog.wercker.com/2014/02/06/RethinkDB-Gingko-Martini-Golang.html
//http://golang.org/pkg/net/http/httptest/#example_ResponseRecorder

import (
	// "encoding/json"
	// "github.com/curt-labs/GoAPI/controllers/faq"
	// "github.com/curt-labs/GoAPI/helpers/encoding"
	"bytes"
	"github.com/go-martini/martini"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCrap(t *testing.T) {
	// Convey("Test it ", t, func() {
	// 	w := httptest.NewRecorder()
	// 	req, err := http.NewRequest("GET", "/", nil)
	// 	So(err, ShouldBeNil)
	// 	res := faq_controller.Test(w, req, nil)
	// 	t.Log(res)
	// 	So(res, ShouldEqual, "Success")

	// })

	Convey("Test GetAll", t, func() {
		body := url.Values{}
		body.Add("id", "1")
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "localhost:8080/faqs", bytes.NewBufferString(body.Encode()))
		So(err, ShouldBeNil)

		http.DefaultServeMux.ServeHTTP(w, req)

		t.Log(w.Body.String())
	})
}

func TestHTTP(t *testing.T) {
	Convey("Test Faqs", t, func() {

		w := httptest.NewRecorder() //implementation of responsewriter
		m := martini.New()

		req, err := http.NewRequest("GET", "/faqs", nil)
		m.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, 200)
		So(err, ShouldBeNil)

		w = httptest.NewRecorder()
		req, err = http.NewRequest("GET", "/faqs/1", nil)
		m.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, 200)
		So(err, ShouldBeNil)

		req, err = http.NewRequest("GET", "/faqs/search?question=hitch", nil)
		m.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, 200)
		So(err, ShouldBeNil)

		body := url.Values{"question": {"test"}, "answer": {"testan"}}
		v := body.Encode()
		s := strings.NewReader(v)
		req, err = http.NewRequest("PUT", "/faqs", s)
		m.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, 200)
		So(err, ShouldBeNil)

		body = url.Values{"question": {"test2"}, "answer": {"testan2"}}
		v = body.Encode()
		s = strings.NewReader(v)
		req, err = http.NewRequest("POST", "/faqs/1", s)
		m.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, 200)
		So(err, ShouldBeNil)

	})

}
