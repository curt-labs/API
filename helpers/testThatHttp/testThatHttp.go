package testThatHttp

import (
	"github.com/codegangsta/martini-contrib/render"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/go-martini/martini"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	dc := &apicontext.DataContext{}
	m.Map(dc)
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
	if Response.Code != 200 {
		log.Print("Response Error: ", Response)
	}
}

func RequestBenchmark(runs int, method, route string, body *url.Values, handler martini.Handler) {

	(&httprunner.BenchmarkOptions{
		Method:      method,
		Route:       route,
		Handler:     handler,
		Runs:        runs,
		QueryString: body,
	}).RequestBenchmark()

}
