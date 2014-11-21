package dealers_ctlr_new

import (
	"encoding/json"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/apiKeyType"
	"github.com/curt-labs/GoAPI/models/customer_new"
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestDealers_New(t *testing.T) {
	var err error
	var c customer_new.Customer
	var cs customer_new.Customers
	var cu customer_new.CustomerUser
	var dt customer_new.DealerType
	var loc customer_new.CustomerLocation
	var locs customer_new.CustomerLocations

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

	//setup apiKeyTypes
	var pub, pri, auth apiKeyType.ApiKeyType
	pub.Type = "public"
	pri.Type = "private"
	auth.Type = "authentication"
	pub.Create()
	pri.Create()
	auth.Create()

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
		So(cs, ShouldHaveSameTypeAs, customer_new.Customers{})
		So(len(cs), ShouldBeGreaterThan, 0)

		//test get local dealers
		thyme = time.Now()
		testThatHttp.Request("get", "/new/dealers/local", "", "?key="+apiKey+"&latlng=43.853282,-95.571675,45.800981,-90.468526&center=44.83536,-93.0201", GetLocalDealers, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var dls customer_new.DealerLocations
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &dls)
		So(err, ShouldBeNil)
		So(dls, ShouldHaveSameTypeAs, customer_new.DealerLocations{})

		//test get local dealers
		thyme = time.Now()
		testThatHttp.Request("get", "/new/dealers/local/regions", "", "?key="+apiKey, GetLocalRegions, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()*5) //Long test
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var stateRegions []customer_new.StateRegion
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &stateRegions)
		So(err, ShouldBeNil)
		So(stateRegions, ShouldHaveSameTypeAs, []customer_new.StateRegion{})

		//test get dealerTypes
		thyme = time.Now()
		testThatHttp.Request("get", "/new/dealers/local/types", "", "?key="+apiKey, GetLocalDealerTypes, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var types []customer_new.DealerType
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &types)
		So(err, ShouldBeNil)
		So(types, ShouldHaveSameTypeAs, []customer_new.DealerType{})

		//test get dealerTiers
		thyme = time.Now()
		testThatHttp.Request("get", "/new/dealers/local/tiers", "", "?key="+apiKey, GetLocalDealerTiers, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var tiers []customer_new.DealerTier
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &tiers)
		So(err, ShouldBeNil)
		So(tiers, ShouldHaveSameTypeAs, []customer_new.DealerTier{})

		//test get dealerTiers
		thyme = time.Now()
		testThatHttp.Request("get", "/new/dealers/etailer/platinum", "", "?key="+apiKey, PlatinumEtailers, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cs)
		So(err, ShouldBeNil)
		So(cs, ShouldHaveSameTypeAs, customer_new.Customers{})

		//test get location by id
		thyme = time.Now()
		testThatHttp.Request("get", "/new/dealers/location/", ":id", strconv.Itoa(loc.Id)+"?key="+apiKey, GetLocationById, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &loc)
		So(err, ShouldBeNil)
		So(loc, ShouldHaveSameTypeAs, customer_new.CustomerLocation{})

		//test get all business classes
		thyme = time.Now()
		testThatHttp.Request("get", "/new/dealers/business/classes", "", "?key="+apiKey, GetAllBusinessClasses, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var bcs customer_new.BusinessClasses
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &bcs)
		So(err, ShouldBeNil)
		So(bcs, ShouldHaveSameTypeAs, customer_new.BusinessClasses{})

		//test search Locations
		thyme = time.Now()
		testThatHttp.Request("get", "/new/dealers/search/", ":search", "test?key="+apiKey, SearchLocations, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &locs)
		So(err, ShouldBeNil)
		So(locs, ShouldHaveSameTypeAs, customer_new.CustomerLocations{})

		//test search Locations by type
		thyme = time.Now()
		testThatHttp.Request("get", "/new/dealers/search/type/", ":search", "test?key="+apiKey, SearchLocationsByType, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &locs)
		So(err, ShouldBeNil)
		So(locs, ShouldHaveSameTypeAs, customer_new.CustomerLocations{})

		//test search Locations by lat long
		thyme = time.Now()
		testThatHttp.Request("get", "/new/dealers/search/geo/", ":latitude/:longitude", fmt.Sprint(c.Latitude)+"/"+fmt.Sprint(c.Longitude)+"?key="+apiKey, SearchLocationsByType, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &dls)
		So(err, ShouldBeNil)
		So(dls, ShouldHaveSameTypeAs, customer_new.DealerLocations{})

	})
	c.Delete()
	cu.Delete()
	err = dt.Delete()
	loc.Delete()

	pub.Delete()
	pri.Delete()
	auth.Delete()
}
