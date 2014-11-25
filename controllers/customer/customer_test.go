package customer_ctlr

import (
	// "bytes"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/curt-labs/GoAPI/models/customer/content"
	"github.com/curt-labs/GoAPI/models/customer_new"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestCustomer(t *testing.T) {
	//setup - populate table with cust
	//old types
	var cust customer.Customer
	var user customer.CustomerUser
	var locs []customer.GeoLocation
	var users []customer.CustomerUser
	var cts []custcontent.ContentType
	var con custcontent.PartContent
	var custCon custcontent.CustomerContent

	//customer_new - for db setup only
	var c customer_new.Customer
	var cu customer_new.CustomerUser

	c.Name = "test custe"
	c.Create()
	cu.CustomerID = c.Id
	cu.Name = "test cust user"
	cu.Email = "pretend@test.com"
	cu.Password = "test"
	cu.Sudo = true
	cu.Create()
	var err error
	var apiKey string
	for _, key := range cu.Keys {
		if strings.ToLower(key.Type) == "public" {
			apiKey = key.Key
		}
	}
	t.Log("APIKEY", apiKey)

	custCon.Save(11000, 1, apiKey)

	Convey("Testing Customer", t, func() {
		//test gets
		form := url.Values{"key": {apiKey}}
		v := form.Encode()
		body := strings.NewReader(v)

		//getCustomer
		testThatHttp.Request("post", "/customer", "", "", GetCustomer, body, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cust)
		So(err, ShouldBeNil)
		So(cust, ShouldHaveSameTypeAs, customer.Customer{})

		//getUser
		testThatHttp.Request("post", "/customer/user", "", "?key="+apiKey, GetUser, body, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &user)
		So(err, ShouldBeNil)
		So(user, ShouldHaveSameTypeAs, customer.CustomerUser{})

		//getLocations
		testThatHttp.Request("post", "/customer/locations", "", "?key="+apiKey, GetLocations, body, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &locs)
		So(err, ShouldBeNil)
		So(locs, ShouldHaveSameTypeAs, []customer.GeoLocation{})

		//getUsers
		testThatHttp.Request("post", "/customer/users", "", "?key="+apiKey, GetUsers, nil, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &users)
		So(err, ShouldBeNil)
		So(users, ShouldHaveSameTypeAs, []customer.CustomerUser{})

		//get content types
		testThatHttp.Request("get", "/cms/content-types", "", "?key="+apiKey, GetAllContentTypes, nil, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cts)
		So(err, ShouldBeNil)
		So(cts, ShouldHaveSameTypeAs, []custcontent.ContentType{})

		//get content types
		testThatHttp.Request("get", "/cms", "", "?key="+apiKey, GetAllContent, nil, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &con)
		So(err, ShouldBeNil)
		So(con, ShouldHaveSameTypeAs, custcontent.PartContent{})

		//get content types
		testThatHttp.Request("get", "/cms/part", "", "?key="+apiKey, AllPartContent, nil, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &con)
		So(err, ShouldBeNil)
		So(con, ShouldHaveSameTypeAs, custcontent.PartContent{})

		//get part content
		testThatHttp.Request("get", "/cms/part/", ":id", strconv.Itoa(11000)+"?key="+apiKey, UniquePartContent, nil, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &con)
		So(err, ShouldBeNil)
		So(con, ShouldHaveSameTypeAs, custcontent.PartContent{})

		//TODO - finish, assuming we ever use this old customer controller

	})
	//teardown
	c.Delete()
	cu.Delete()

	custCon.Delete(11000, 1, apiKey)
}

func BenchmarkCustomer(b *testing.B) {
	//create customer
	var cu customer_new.CustomerUser
	var apiKey string
	cu.Name = "test cust user"
	cu.Email = "pretend@test.com"
	cu.Password = "test"
	cu.Sudo = true
	cu.Create()

	for _, key := range cu.Keys {
		if strings.ToLower(key.Type) == "public" {
			apiKey = key.Key
		}
	}

	qs := make(url.Values, 0)
	qs.Add("key", apiKey)

	testThatHttp.RequestBenchmark(b.N, "POST", "/customer", &qs, GetCustomer) //TODO form
	testThatHttp.RequestBenchmark(b.N, "POST", "/customer/user", &qs, GetUser)
	testThatHttp.RequestBenchmark(b.N, "POST", "/customer/locations", &qs, GetLocations)
	testThatHttp.RequestBenchmark(b.N, "POST", "/customer/users", &qs, GetUsers)
	testThatHttp.RequestBenchmark(b.N, "GET", "/cms/content-types", nil, GetAllContentTypes)
	testThatHttp.RequestBenchmark(b.N, "GET", "/cms/", nil, GetAllContent)
	testThatHttp.RequestBenchmark(b.N, "GET", "/cms/part", nil, AllPartContent)
	testThatHttp.RequestBenchmark(b.N, "GET", "/cms/part/11000", nil, UniquePartContent)

	cu.Delete()
}
