package testThatHttp

import (
	"github.com/codegangsta/martini-contrib/render"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/go-martini/martini"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

var (
	Response *httptest.ResponseRecorder
)

func Request(reqType string, route string, paramKey string, paramVal string, handler martini.Handler, body io.Reader, contentType string) {
	cType := "application/x-www-form-urlencoded" //default content type = form
	if contentType != "" {
		cType = contentType
	}
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(encoding.MapEncoder)
	reqType = strings.ToLower(reqType)
	switch {
	case reqType == "get":
		m.Get(route+paramKey, handler)
	case reqType == "post":
		m.Post(route+paramKey, handler)
	case reqType == "put":
		m.Put(route+paramKey, handler)
	case reqType == "delete":
		m.Delete(route+paramKey, handler)
	}

	request, _ := http.NewRequest(strings.ToUpper(reqType), route+paramVal, body)
	request.Header.Set("Content-Type", cType) //content-type=form,json,etc
	Response = httptest.NewRecorder()
	m.ServeHTTP(Response, request)
}
