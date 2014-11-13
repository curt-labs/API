package main

//reference
//https://groups.google.com/forum/#!topic/golang-nuts/DARY7HY-pbY
//http://blog.wercker.com/2014/02/06/RethinkDB-Gingko-Martini-Golang.html
//http://golang.org/pkg/net/http/httptest/#example_ResponseRecorder

import (
	"github.com/go-martini/martini"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStuff(t *testing.T) {
	Convey("test", t, func() {
		req, _ := http.NewRequest("GET", "/faqs", nil)
		w := httptest.NewRecorder()
		m := martini.New()
		m.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, 200)

	})

}
