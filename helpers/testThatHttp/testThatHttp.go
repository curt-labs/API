package testThatHttp

import (
	"github.com/codegangsta/martini-contrib/render"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/models/customer_new"
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
	//create customerUser and key
	var cu customer_new.CustomerUser
	cu.Name = "test cust user"
	cu.Email = "pretend@test.com"
	cu.Password = "test"
	cu.Sudo = true
	cu.Create()
	var apiKey string
	for _, key := range cu.Keys {
		if strings.ToLower(key.Type) == "public" {
			apiKey = key.Key
		}
	}

	//opts
	opts := httprunner.ReqOpts{
		Body:    body,
		Handler: handler,
		URL:     route + "?key=" + apiKey,
		Method:  method,
	}

	//run
	(&httprunner.Runner{
		Req: &opts,
		N:   runs,
		C:   1,
	}).Run()

	//teardown
	cu.Delete()

}
