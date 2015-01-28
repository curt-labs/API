package dealers_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/customer"
	. "github.com/smartystreets/goconvey/convey"

	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestDealers_New(t *testing.T) {
	var err error
	var c customer.Customer
	var cs customer.Customers
	var loc customer.CustomerLocation
	var locs customer.CustomerLocations

	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}
	//setup lcoation for getById
	loc.Address = "123 Test Ave."
	loc.Create(dtx)

	Convey("Testing Dealers_New", t, func() {
		//test get etailers
		thyme := time.Now()
		testThatHttp.RequestWithDtx("get", "/dealers/etailer", "", "?key="+dtx.APIKey, GetEtailers, nil, "", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds())
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cs)
		So(err, ShouldBeNil)
		So(cs, ShouldHaveSameTypeAs, customer.Customers{})
		So(len(cs), ShouldBeGreaterThanOrEqualTo, 0)

		//test get local dealers
		thyme = time.Now()
		testThatHttp.Request("get", "/dealers/local", "", "?key="+dtx.APIKey+"&latlng=43.853282,-95.571675,45.800981,-90.468526&center=44.83536,-93.0201", GetLocalDealers, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var dls customer.DealerLocations
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &dls)
		So(err, ShouldBeNil)
		So(dls, ShouldHaveSameTypeAs, customer.DealerLocations{})

		//test get local dealers
		thyme = time.Now()
		testThatHttp.Request("get", "/dealers/local/regions", "", "?key="+dtx.APIKey, GetLocalRegions, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()*5) //Long test
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var stateRegions []customer.StateRegion
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &stateRegions)
		So(err, ShouldBeNil)
		So(stateRegions, ShouldHaveSameTypeAs, []customer.StateRegion{})

		//test get dealerTypes
		thyme = time.Now()
		testThatHttp.Request("get", "/dealers/local/types", "", "?key="+dtx.APIKey, GetLocalDealerTypes, nil, "")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var types []customer.DealerType
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &types)
		So(err, ShouldBeNil)
		So(types, ShouldHaveSameTypeAs, []customer.DealerType{})

		//test get dealerTiers
		thyme = time.Now()
		testThatHttp.Request("get", "/dealers/local/tiers", "", "?key="+dtx.APIKey, GetLocalDealerTiers, nil, "")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var tiers []customer.DealerTier
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &tiers)
		So(err, ShouldBeNil)
		So(tiers, ShouldHaveSameTypeAs, []customer.DealerTier{})

		//test get dealerTiers
		thyme = time.Now()
		testThatHttp.Request("get", "/dealers/etailer/platinum", "", "?key="+dtx.APIKey, PlatinumEtailers, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cs)
		So(err, ShouldBeNil)
		So(cs, ShouldHaveSameTypeAs, customer.Customers{})

		//test get location by id
		thyme = time.Now()
		testThatHttp.Request("get", "/dealers/location/", ":id", strconv.Itoa(loc.Id)+"?key="+dtx.APIKey, GetLocationById, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &loc)
		So(err, ShouldBeNil)
		So(loc, ShouldHaveSameTypeAs, customer.CustomerLocation{})

		//test get all business classes
		thyme = time.Now()
		testThatHttp.Request("get", "/dealers/business/classes", "", "?key="+dtx.APIKey, GetAllBusinessClasses, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var bcs customer.BusinessClasses
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &bcs)
		So(err, ShouldBeNil)
		So(bcs, ShouldHaveSameTypeAs, customer.BusinessClasses{})

		//test search Locations
		thyme = time.Now()
		testThatHttp.Request("get", "/dealers/search/", ":search", "test?key="+dtx.APIKey, SearchLocations, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &locs)
		So(err, ShouldBeNil)
		So(locs, ShouldHaveSameTypeAs, customer.CustomerLocations{})

		//test search Locations by type
		thyme = time.Now()
		testThatHttp.Request("get", "/dealers/search/type/", ":search", "test?key="+dtx.APIKey, SearchLocationsByType, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &locs)
		So(err, ShouldBeNil)
		So(locs, ShouldHaveSameTypeAs, customer.CustomerLocations{})

		//test search Locations by lat long
		thyme = time.Now()
		testThatHttp.Request("get", "/dealers/search/geo/", ":latitude/:longitude", fmt.Sprint(c.Latitude)+"/"+fmt.Sprint(c.Longitude)+"?key="+dtx.APIKey, SearchLocationsByType, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()*2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &dls)
		So(err, ShouldBeNil)
		So(dls, ShouldHaveSameTypeAs, customer.DealerLocations{})

	})
	//teardown
	loc.Delete(dtx)
	_ = apicontextmock.DeMock(dtx)
}

func BenchmarkCRUDDealers(b *testing.B) {

	dtx, err := apicontextmock.Mock()
	if err != nil {
		b.Log(err)
	}

	qs := make(url.Values, 0)
	qs.Set("key", dtx.APIKey)

	qs2 := make(url.Values, 0)
	qs2.Set("key", dtx.APIKey)
	qs2.Set("latlng", "43.853282,-95.571675,45.800981,-90.468526&center=44.83536,-93.0201")

	Convey("CustomerDealers", b, func() {
		//get etailers
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/dealers/etailer",
			ParameterizedRoute: "/dealers/etailer",
			Handler:            GetEtailers,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()

		//get local dealers
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/dealers/local",
			ParameterizedRoute: "/dealers/local",
			Handler:            GetLocalDealers,
			QueryString:        &qs2,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
		//get local dealers
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/dealers/local/regions",
			ParameterizedRoute: "/dealers/local/regions",
			Handler:            GetLocalRegions,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
		//get local dealers
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/dealers/local/types",
			ParameterizedRoute: "/dealers/local/types",
			Handler:            GetLocalDealerTypes,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
		//get local dealers
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/dealers/local/tiers",
			ParameterizedRoute: "/dealers/local/tiers",
			Handler:            GetLocalDealerTiers,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
		//get local dealers
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/dealers/etailer/platinum",
			ParameterizedRoute: "/dealers/etailer/platinum",
			Handler:            PlatinumEtailers,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
		//get local dealers
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/dealers/location/",
			ParameterizedRoute: "/dealers/location/1",
			Handler:            GetLocationById,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
		//get local dealers
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/dealers/business/classes",
			ParameterizedRoute: "/dealers/business/classes",
			Handler:            GetAllBusinessClasses,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
		//get local dealers
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/dealers/search/",
			ParameterizedRoute: "/dealers/search/hitch",
			Handler:            SearchLocations,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
		//get local dealers
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/dealers/search/type/",
			ParameterizedRoute: "/dealers/search/type/installer",
			Handler:            SearchLocationsByType,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
		//get local dealers
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/dealers/search/geo/",
			ParameterizedRoute: "/dealers/search/geo/44.83536/-93.0201",
			Handler:            SearchLocationsByLatLng,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
	})
	_ = apicontextmock.DeMock(dtx)
}
