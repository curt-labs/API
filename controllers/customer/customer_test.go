package customer_ctlr

import (
	"bytes"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/apiKeyType"
	"github.com/curt-labs/GoAPI/models/cartIntegration"
	"github.com/curt-labs/GoAPI/models/customer"
	//"github.com/curt-labs/GoAPI/models/products"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestCustomer(t *testing.T) {
	var c customer.Customer
	var cu customer.CustomerUser
	//var p products.Part
	//p.ID = 123
	//p.Create()
	//var price customer.Price
	//price.PartID = p.ID
	//price.Price = 1000000

	//var ci cartIntegration.CartIntegration
	//ci.PartID = p.ID
	//ci.CustPartID = 987654321

	var pub, pri, auth apiKeyType.ApiKeyType
	if database.EmptyDb != nil {
		t.Log("clean db")
		//setup apiKeyTypes
		pub.Type = "Public"
		pri.Type = "Private"
		auth.Type = "Authentication"
		pub.Create()
		pri.Create()
		auth.Create()
	}

	//setup user
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

	Convey("Testing customer/Customer", t, func() {
		//test create customer
		c.Name = "Jason Voorhees"
		c.Email = "jason@crystal.lake"
		bodyBytes, _ := json.Marshal(c)
		bodyJson := bytes.NewReader(bodyBytes)
		thyme := time.Now()
		testThatHttp.Request("put", "/new/customer", "", "?key="+apiKey, SaveCustomer, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, customer.Customer{})
		So(c.Id, ShouldBeGreaterThan, 0)

		//test update customer
		c.Fax = "666-1313"
		c.State.Id = 1
		bodyBytes, _ = json.Marshal(c)
		bodyJson = bytes.NewReader(bodyBytes)
		thyme = time.Now()
		testThatHttp.Request("put", "/new/customer/", ":id", strconv.Itoa(c.Id)+"?key="+apiKey, SaveCustomer, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, customer.Customer{})
		So(c.Id, ShouldBeGreaterThan, 0)

		//test get customer
		c.JoinUser(cu)
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/", ":id", strconv.Itoa(c.Id)+"?key="+apiKey, GetCustomer, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, customer.Customer{})
		So(c.Id, ShouldBeGreaterThan, 0)
		//same get customer as a post request
		thyme = time.Now()
		testThatHttp.Request("post", "/new/customer/", ":id", strconv.Itoa(c.Id)+"?key="+apiKey, GetCustomer, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, customer.Customer{})
		So(c.Id, ShouldBeGreaterThan, 0)

		// get customer locations
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/locations", "", "?key="+apiKey, GetLocations, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c.Locations)
		So(err, ShouldBeNil)
		So(c.Locations, ShouldHaveSameTypeAs, []customer.CustomerLocation{})

		// get customer locations via post
		thyme = time.Now()
		testThatHttp.Request("post", "/new/customer/locations", "", "?key="+apiKey, GetLocations, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c.Locations)
		So(err, ShouldBeNil)
		So(c.Locations, ShouldHaveSameTypeAs, []customer.CustomerLocation{})

		//get user
		thyme = time.Now()
		testThatHttp.Request("post", "/new/customer/user", "", "?key="+apiKey, GetUser, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cu)
		So(err, ShouldBeNil)
		So(cu, ShouldHaveSameTypeAs, customer.CustomerUser{})

		//get users
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/users", "", "?key="+apiKey, GetUsers, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var cus []customer.CustomerUser
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cus)
		So(err, ShouldBeNil)
		So(cus, ShouldHaveSameTypeAs, []customer.CustomerUser{})

		//get customer price
		price.CustID = c.Id
		price.Create()
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/price/", ":id", strconv.Itoa(p.ID)+"?key="+apiKey, GetCustomerPrice, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var price float64
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &price)
		So(err, ShouldBeNil)
		So(price, ShouldHaveSameTypeAs, 7.1)

		//get customer cart reference
		ci.CustID = c.Id
		ci.Create()
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/cartRef/", ":id", strconv.Itoa(p.ID)+"?key="+apiKey, GetCustomerCartReference, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var reference int
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &reference)
		So(err, ShouldBeNil)
		So(reference, ShouldHaveSameTypeAs, 7)

		//test delete customer
		thyme = time.Now()
		testThatHttp.Request("delete", "/new/customer/", ":id", strconv.Itoa(c.Id)+"?key="+apiKey, DeleteCustomer, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, customer.Customer{})
		So(c.Id, ShouldBeGreaterThan, 0)

	})
	//cleanup
	err = cu.Delete()
	t.Log("CU", err)
	p.Delete()
	price.Delete()
	ci.Delete()

	if database.EmptyDb != nil {
		err = pub.Delete()

		err = pri.Delete()

		err = auth.Delete()
	}

}

func BenchmarkCRUDCustomer(b *testing.B) {
	//get apiKey by creating customeruser
	var cu customer.CustomerUser
	var apiKey string
	cu.Name = "test cust benchmark user"
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

	Convey("Customer", b, func() {
		var c customer.Customer
		c.Name = "Freddy Krueger"
		c.Email = "freddy@elm.st"
		// var locs customer.CustomerLocations

		//create
		(&httprunner.BenchmarkOptions{
			Method:             "POST",
			Route:              "/new/customer",
			ParameterizedRoute: "/new/customer",
			Handler:            SaveCustomer,
			QueryString:        &qs,
			JsonBody:           c,
			Runs:               b.N,
		}).RequestBenchmark()

		//get
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/new/customer",
			ParameterizedRoute: "/new/customer/" + strconv.Itoa(c.Id),
			Handler:            GetCustomer,
			QueryString:        &qs,
			JsonBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()

		//get locations
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/new/customer",
			ParameterizedRoute: "/new/customer",
			Handler:            GetLocations,
			QueryString:        &qs,
			JsonBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()

		//get user
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/new/customer/user",
			ParameterizedRoute: "/new/customer/user",
			Handler:            GetUser,
			QueryString:        &qs,
			JsonBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()

		//get users
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/new/customer/users",
			ParameterizedRoute: "/new/customer/users",
			Handler:            GetUsers,
			QueryString:        &qs,
			JsonBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()

		//set up price & part
		var p products.Part
		p.ID = 123
		p.Create()
		var price customer.Price
		price.CustID = c.Id
		price.Create()

		//get price
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/new/customer/price",
			ParameterizedRoute: "/new/customer/price/" + strconv.Itoa(p.ID),
			Handler:            GetUser,
			QueryString:        &qs,
			JsonBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()

		//get cart ref
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/new/customer/cartRef",
			ParameterizedRoute: "/new/customer/cartRef/" + strconv.Itoa(p.ID),
			Handler:            GetUser,
			QueryString:        &qs,
			JsonBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()

		//delete
		(&httprunner.BenchmarkOptions{
			Method:             "DELETE",
			Route:              "/new/customer",
			ParameterizedRoute: "/new/customer/" + strconv.Itoa(c.Id),
			Handler:            DeleteCustomer,
			QueryString:        &qs,
			JsonBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()

		//tear down price and part
		price.Delete()
		p.Delete()
	})
	//teardown customer user
	cu.Delete()

}
