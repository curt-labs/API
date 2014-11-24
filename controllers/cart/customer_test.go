package cart_ctlr

import (
	"bytes"
	"encoding/json"
	"github.com/curt-labs/GoAPI/controllers/middleware"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/models/cart"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var (
	response *httptest.ResponseRecorder
)

func TestGetCustomers(t *testing.T) {
	Convey("Testing GetCustomers", t, func() {
		Convey("no shop identifier", func() {
			Request("GET", "/shopify/customers", nil, GetCustomers)
			So(response.Code, ShouldEqual, 500)
			So(response.Body.String(), ShouldNotEqual, "[]")

			vals := make(url.Values, 0)
			vals.Add("shop", "testing")
			Request("GET", "/shopify/customers", &vals, GetCustomers)
			So(response.Code, ShouldEqual, 500)
			So(response.Body.String(), ShouldNotEqual, "[]")

			vals.Add("since_id", "something")
			Request("GET", "/shopify/customers", &vals, GetCustomers)
			So(response.Code, ShouldEqual, 500)
			So(response.Body.String(), ShouldNotEqual, "[]")
		})
		Convey("with shop identifier", func() {
			shopID := cart.InsertTestData()
			So(shopID, ShouldNotBeNil)

			val := shopID.Hex()
			qs := make(url.Values, 0)
			qs.Add("shop", val)

			Request("GET", "/shopify/customers", &qs, GetCustomers)
			So(response.Code, ShouldEqual, 200)
			So(json.Unmarshal(response.Body.Bytes(), &[]cart.Customer{}), ShouldBeNil)

			qs.Add("since_id", "something")
			Request("GET", "/shopify/customers", &qs, GetCustomers)
			So(response.Code, ShouldEqual, 200)
			So(json.Unmarshal(response.Body.Bytes(), &[]cart.Customer{}), ShouldBeNil)

			qs.Add("since_id", val)
			Request("GET", "/shopify/customers", &qs, GetCustomers)
			So(response.Code, ShouldEqual, 200)
			var custs []cart.Customer
			So(json.Unmarshal(response.Body.Bytes(), &custs), ShouldBeNil)
			So(len(custs), ShouldEqual, 0)
		})
	})
}

func BenchmarkGetCustomers(b *testing.B) {
	shopID := cart.InsertTestData()
	if shopID == nil {
		panic("shopID cannot be nil")
	}

	val := shopID.Hex()
	qs := make(url.Values, 0)
	qs.Add("shop", val)
	qs.Add("since_id", val)
	RequestBenchmark(b.N, "GET", "/shopify/customers", &qs, GetCustomers)
}

func Request(method, route string, body *url.Values, handler martini.Handler) {
	m := martini.Classic()
	switch strings.ToUpper(method) {
	case "GET":
		m.Get(route, handler)
	case "POST":
		m.Post(route, handler)
	case "PUT":
		m.Put(route, handler)
	case "PATCH":
		m.Patch(route, handler)
	case "DELETE":
		m.Delete(route, handler)
	case "HEAD":
		m.Head(route, handler)
	default:
		m.Any(route, handler)
	}

	m.Use(render.Renderer())
	m.Use(encoding.MapEncoder)
	m.Use(middleware.Meddler())

	var request *http.Request
	if body != nil && strings.ToUpper(method) != "GET" {
		request, _ = http.NewRequest(method, route, bytes.NewBufferString(body.Encode()))
	} else if body != nil {
		request, _ = http.NewRequest(method, route+"?"+body.Encode(), nil)
	} else {
		request, _ = http.NewRequest(method, route, nil)
	}

	response = httptest.NewRecorder()
	m.ServeHTTP(response, request)
}

func RequestBenchmark(runs int, method, route string, body *url.Values, handler martini.Handler) {

	opts := httprunner.ReqOpts{
		Body:    body,
		Handler: handler,
		URL:     route,
		Method:  method,
	}

	(&httprunner.Runner{
		Req: &opts,
		N:   runs,
		C:   1,
	}).Run()

}
