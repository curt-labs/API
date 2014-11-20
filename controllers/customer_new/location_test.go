package customer_ctlr_new

import (
	"bytes"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/customer_new"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestCustomerLocation(t *testing.T) {
	var err error
	var loc customer_new.CustomerLocation
	var cu customer_new.CustomerUser

	//setup
	cu.Name = "test cust user"
	cu.Email = "pretend@test.com"
	cu.Password = "test"
	cu.Sudo = true
	var apiKey string
	for _, key := range cu.Keys {
		if strings.ToLower(key.Type) == "public" {
			apiKey = key.Key
		}
	}
	t.Log("APIKEY", apiKey)

	Convey("Testing Customer_New/Location", t, func() {
		//test create customer location
		form := url.Values{"name": {"Dave Grohl"}, "address": {"404 S. Barstow St."}, "city": {"Eau Claire"}}
		v := form.Encode()
		body := strings.NewReader(v)
		thyme := time.Now()
		testThatHttp.Request("post", "/new/customer/location", "", "?key="+apiKey, SaveLocation, body, "application/x-www-form-urlencoded")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &loc)
		So(err, ShouldBeNil)
		So(loc, ShouldHaveSameTypeAs, customer_new.CustomerLocation{})
		So(loc.Id, ShouldBeGreaterThan, 0)

		//test update location with json
		loc.Fax = "715-839-0000"
		bodyBytes, _ := json.Marshal(loc)
		bodyJson := bytes.NewReader(bodyBytes)
		thyme = time.Now()
		testThatHttp.Request("put", "/new/customer/location/", ":id", strconv.Itoa(loc.Id)+"?key="+apiKey, SaveLocation, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &loc)
		So(err, ShouldBeNil)
		So(loc, ShouldHaveSameTypeAs, customer_new.CustomerLocation{})

		//test get location
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/location/", ":id", strconv.Itoa(loc.Id)+"?key="+apiKey, GetLocation, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &loc)
		So(err, ShouldBeNil)
		So(loc, ShouldHaveSameTypeAs, customer_new.CustomerLocation{})

		//test get location
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/location", "", "?key="+apiKey, GetAllLocations, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var locs customer_new.CustomerLocations
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &locs)
		So(err, ShouldBeNil)
		So(locs, ShouldHaveSameTypeAs, customer_new.CustomerLocations{})

		//test delete location
		thyme = time.Now()
		testThatHttp.Request("delete", "/new/customer/location", "", "?key="+apiKey, DeleteLocation, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &loc)
		So(err, ShouldBeNil)
		So(loc, ShouldHaveSameTypeAs, customer_new.CustomerLocation{})

	})
	//teardown
	cu.Delete()
}
