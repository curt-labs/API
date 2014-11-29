package dealers_ctlr

import (
	"encoding/json"
	"fmt"
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

func TestDealers_New(t *testing.T) {
	var err error
	var c customer.Customer
	var cs customer.Customers
	var cu customer.CustomerUser
	var dt customer.DealerType
	var loc customer.CustomerLocation
	var locs customer.CustomerLocations

	loc.Address = "123 Test Ave."
	loc.Create()

	//need dealer type (online) for empty DB
	dt.Type = "test type"
	dt.Online = true
	err = dt.Create()

	c.Name = "Dog Bountyhunter"
	c.DealerType.Id = dt.Id //is etailer
	c.IsDummy = false
	c.Latitude = 44.83536
	c.Longitude = -93.0201
	err = c.Create()

	cu.CustomerID = c.Id
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
	t.Log("APIKEY", apiKey)

	Convey("Testing Dealers_New", t, func() {
		//test get etailers
		thyme := time.Now()
		testThatHttp.Request("get", "/new/dealers/etailer", "", "?key="+apiKey, GetEtailers, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)

		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cs)
		So(err, ShouldBeNil)
		So(cs, ShouldHaveSameTypeAs, customer.Customers{})
		So(len(cs), ShouldBeGreaterThan, 0)

		//test get local dealers
		thyme = time.Now()
		testThatHttp.Request("get", "/new/dealers/local", "", "?key="+apiKey+"&latlng=43.853282,-95.571675,45.800981,-90.468526&center=44.83536,-93.0201", GetLocalDealers, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var dls customer.DealerLocations
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &dls)
		So(err, ShouldBeNil)
		So(dls, ShouldHaveSameTypeAs, customer.DealerLocations{})

		//test get local dealers
		thyme = time.Now()
		testThatHttp.Request("get", "/new/dealers/local/regions", "", "?key="+apiKey, GetLocalRegions, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()*5) //Long test
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var stateRegions []customer.StateRegion
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &stateRegions)
		So(err, ShouldBeNil)
		So(stateRegions, ShouldHaveSameTypeAs, []customer.StateRegion{})

		//test get dealerTypes
		thyme = time.Now()
		testThatHttp.Request("get", "/new/dealers/local/types", "", "?key="+apiKey, GetLocalDealerTypes, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var types []customer.DealerType
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &types)
		So(err, ShouldBeNil)
		So(types, ShouldHaveSameTypeAs, []customer.DealerType{})

		//test get dealerTiers
		thyme = time.Now()
		testThatHttp.Request("get", "/new/dealers/local/tiers", "", "?key="+apiKey, GetLocalDealerTiers, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var tiers []customer.DealerTier
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &tiers)
		So(err, ShouldBeNil)
		So(tiers, ShouldHaveSameTypeAs, []customer.DealerTier{})

		//test get dealerTiers
		thyme = time.Now()
		testThatHttp.Request("get", "/new/dealers/etailer/platinum", "", "?key="+apiKey, PlatinumEtailers, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cs)
		So(err, ShouldBeNil)
		So(cs, ShouldHaveSameTypeAs, customer.Customers{})

		//test get location by id
		thyme = time.Now()
		testThatHttp.Request("get", "/new/dealers/location/", ":id", strconv.Itoa(loc.Id)+"?key="+apiKey, GetLocationById, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &loc)
		So(err, ShouldBeNil)
		So(loc, ShouldHaveSameTypeAs, customer.CustomerLocation{})

		//test get all business classes
		thyme = time.Now()
		testThatHttp.Request("get", "/new/dealers/business/classes", "", "?key="+apiKey, GetAllBusinessClasses, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var bcs customer.BusinessClasses
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &bcs)
		So(err, ShouldBeNil)
		So(bcs, ShouldHaveSameTypeAs, customer.BusinessClasses{})

		//test search Locations
		thyme = time.Now()
		testThatHttp.Request("get", "/new/dealers/search/", ":search", "test?key="+apiKey, SearchLocations, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &locs)
		So(err, ShouldBeNil)
		So(locs, ShouldHaveSameTypeAs, customer.CustomerLocations{})

		//test search Locations by type
		thyme = time.Now()
		testThatHttp.Request("get", "/new/dealers/search/type/", ":search", "test?key="+apiKey, SearchLocationsByType, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &locs)
		So(err, ShouldBeNil)
		So(locs, ShouldHaveSameTypeAs, customer.CustomerLocations{})

		//test search Locations by lat long
		thyme = time.Now()
		testThatHttp.Request("get", "/new/dealers/search/geo/", ":latitude/:longitude", fmt.Sprint(c.Latitude)+"/"+fmt.Sprint(c.Longitude)+"?key="+apiKey, SearchLocationsByType, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &dls)
		So(err, ShouldBeNil)
		So(dls, ShouldHaveSameTypeAs, customer.DealerLocations{})

	})
	c.Delete()
	cu.Delete()
	err = dt.Delete()
	loc.Delete()

}

func BenchmarkCRUDDealers(b *testing.B) {
	var cu customer.CustomerUser
	cu.Name = "test dealer cust user"
	cu.Email = "dealer@test.com"
	cu.Password = "test"
	cu.Sudo = true
	cu.Create()
	var apiKey string
	for _, key := range cu.Keys {
		if strings.ToLower(key.Type) == "public" {
			apiKey = key.Key
		}
	}

	qs := make(url.Values, 0)
	qs.Set("key", apiKey)

	qs2 := make(url.Values, 0)
	qs2.Set("key", apiKey)
	qs2.Set("latlng", "43.853282,-95.571675,45.800981,-90.468526&center=44.83536,-93.0201")

	Convey("CustomerDealers", b, func() {
		//get etailers
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/new/dealers/etailer",
			ParameterizedRoute: "/new/dealers/etailer",
			Handler:            GetEtailers,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()

		//get local dealers
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/new/dealers/local",
			ParameterizedRoute: "/new/dealers/local",
			Handler:            GetLocalDealers,
			QueryString:        &qs2,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
		//get local dealers
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/new/dealers/local/regions",
			ParameterizedRoute: "/new/dealers/local/regions",
			Handler:            GetLocalRegions,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
		//get local dealers
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/new/dealers/local/types",
			ParameterizedRoute: "/new/dealers/local/types",
			Handler:            GetLocalDealerTypes,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
		//get local dealers
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/new/dealers/local/tiers",
			ParameterizedRoute: "/new/dealers/local/tiers",
			Handler:            GetLocalDealerTiers,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
		//get local dealers
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/new/dealers/etailer/platinum",
			ParameterizedRoute: "/new/dealers/etailer/platinum",
			Handler:            PlatinumEtailers,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
		//get local dealers
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/new/dealers/location/",
			ParameterizedRoute: "/new/dealers/location/1",
			Handler:            GetLocationById,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
		//get local dealers
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/new/dealers/business/classes",
			ParameterizedRoute: "/new/dealers/business/classes",
			Handler:            GetAllBusinessClasses,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
		//get local dealers
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/new/dealers/search/",
			ParameterizedRoute: "/new/dealers/search/hitch",
			Handler:            SearchLocations,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
		//get local dealers
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/new/dealers/search/type/",
			ParameterizedRoute: "/new/dealers/search/type/installer",
			Handler:            SearchLocationsByType,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
		//get local dealers
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/new/dealers/search/geo/",
			ParameterizedRoute: "/new/dealers/search/geo/44.83536/-93.0201",
			Handler:            SearchLocationsByLatLng,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
	})
	cu.Delete()
}
