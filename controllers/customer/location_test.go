package customer_ctlr

import (
	"bytes"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/customer"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestCustomerLocation(t *testing.T) {
	var err error
	var loc customer.CustomerLocation

	Convey("Testing customer/Location", t, func() {
		//test create customer location
		form := url.Values{"name": {"Dave Grohl"}, "address": {"404 S. Barstow St."}, "city": {"Eau Claire"}}
		v := form.Encode()
		body := strings.NewReader(v)
		thyme := time.Now()
		testThatHttp.Request("post", "/customer/location", "", "", SaveLocation, body, "application/x-www-form-urlencoded")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()*2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &loc)
		So(err, ShouldBeNil)
		So(loc, ShouldHaveSameTypeAs, customer.CustomerLocation{})
		So(loc.Id, ShouldBeGreaterThan, 0)

		//test update location with json
		loc.Fax = "715-839-0000"
		bodyBytes, _ := json.Marshal(loc)
		bodyJson := bytes.NewReader(bodyBytes)
		thyme = time.Now()
		testThatHttp.Request("put", "/customer/location/", ":id", strconv.Itoa(loc.Id), SaveLocation, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &loc)
		So(err, ShouldBeNil)
		So(loc, ShouldHaveSameTypeAs, customer.CustomerLocation{})

		//test get location
		thyme = time.Now()
		testThatHttp.Request("get", "/customer/location/", ":id", strconv.Itoa(loc.Id), GetLocation, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &loc)
		So(err, ShouldBeNil)
		So(loc, ShouldHaveSameTypeAs, customer.CustomerLocation{})

		//test get all locations
		thyme = time.Now()
		testThatHttp.Request("get", "/customer/location", "", "", GetAllLocations, bodyJson, "application/json")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var locs customer.CustomerLocations
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &locs)
		So(err, ShouldBeNil)
		So(locs, ShouldHaveSameTypeAs, customer.CustomerLocations{})

		//test delete location
		thyme = time.Now()
		testThatHttp.Request("delete", "/new/customer/location/", ":id", strconv.Itoa(loc.Id), DeleteLocation, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &loc)
		So(err, ShouldBeNil)
		So(loc, ShouldHaveSameTypeAs, customer.CustomerLocation{})

	})
}

func BenchmarkCRUDCustomerLocation(b *testing.B) {

	qs := make(url.Values, 0)
	var loc customer.CustomerLocation

	Convey("CustomerLocation", b, func() {
		form := url.Values{"name": {"Dave Grohl"}, "address": {"404 S. Barstow St."}, "city": {"Eau Claire"}}
		//create
		(&httprunner.BenchmarkOptions{
			Method:             "POST",
			Route:              "/customer/location",
			ParameterizedRoute: "/customer/location",
			Handler:            SaveLocation,
			QueryString:        &qs,
			JsonBody:           loc,
			FormBody:           form,
			Runs:               b.N,
		}).RequestBenchmark()

		//get
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/customer/location",
			ParameterizedRoute: "/customer/location/" + strconv.Itoa(loc.Id),
			Handler:            GetLocation,
			QueryString:        &qs,
			JsonBody:           loc,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()

		//get all
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/customer/location",
			ParameterizedRoute: "/customer/location",
			Handler:            GetLocations,
			QueryString:        &qs,
			JsonBody:           loc,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()

		//delete
		(&httprunner.BenchmarkOptions{
			Method:             "DELETE",
			Route:              "/customer/location",
			ParameterizedRoute: "/customer/location/" + strconv.Itoa(loc.Id),
			Handler:            DeleteLocation,
			QueryString:        &qs,
			JsonBody:           loc,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
	})
}
