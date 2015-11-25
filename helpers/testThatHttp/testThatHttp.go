package testThatHttp

import (
	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/apicontextmock"
	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/httprunner"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
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
	MockedDTX := &apicontext.DataContext{}

	var apiKey string
	api := strings.SplitAfter(paramVal, "?key=")
	if len(api) > 1 {
		apiKey = api[1]
	}

	cType := "application/x-www-form-urlencoded" //default content type = form
	if contentType != "" {
		cType = contentType
	}
	m := martini.Classic()
	m.Use(render.Renderer())
	dc := &apicontext.DataContext{APIKey: apiKey, BrandID: 1, WebsiteID: 1}
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
	_ = apicontextmock.DeMock(MockedDTX)
}

func RequestWithDtx(reqType string, route string, paramKey string, paramVal string, handler martini.Handler, body io.Reader, contentType string, dtx *apicontext.DataContext) {

	cType := "application/x-www-form-urlencoded" //default content type = form
	if contentType != "" {
		cType = contentType
	}
	m := martini.Classic()
	m.Use(render.Renderer())

	m.Map(dtx)
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
	// _ = apicontextmock.DeMock(dtx)
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
