package customer_ctlr

import (
	"bytes"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"

	// "github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/customer"
	. "github.com/smartystreets/goconvey/convey"

	"strconv"
	"testing"
	"time"
)

func TestCustomer(t *testing.T) {
	var err error
	var c customer.Customer
	var cu customer.CustomerUser
	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}

	Convey("Testing customer/Customer", t, func() {
		//test create customer
		c.Name = "Jason Voorhees"
		c.Email = "jason@crystal.lake"
		bodyBytes, _ := json.Marshal(c)
		bodyJson := bytes.NewReader(bodyBytes)
		thyme := time.Now()
		testThatHttp.Request("put", "/customer", "", "?key=", SaveCustomer, bodyJson, "application/json")
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
		testThatHttp.Request("put", "/customer/", ":id", strconv.Itoa(c.Id)+"?key=", SaveCustomer, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, customer.Customer{})
		So(c.Id, ShouldBeGreaterThan, 0)

		thyme = time.Now()
		testThatHttp.Request("get", "/customer/", ":id", strconv.Itoa(c.Id)+"?key=", GetCustomer, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, customer.Customer{})
		So(c.Id, ShouldBeGreaterThan, 0)

		// get customer locations
		thyme = time.Now()
		testThatHttp.RequestWithDtx("get", "/customer/locations", "", "?key=", GetLocations, nil, "", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds())
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c.Locations)
		So(err, ShouldBeNil)
		So(c.Locations, ShouldHaveSameTypeAs, []customer.CustomerLocation{})

		// //get user
		thyme = time.Now()
		testThatHttp.RequestWithDtx("post", "/customer/user", "", "?key=", GetUser, nil, "", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cu)
		So(err, ShouldBeNil)
		So(cu, ShouldHaveSameTypeAs, customer.CustomerUser{})

		// //get users
		thyme = time.Now()
		testThatHttp.RequestWithDtx("get", "/customer/users", "", "?key=", GetUsers, nil, "", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds())
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var cus []customer.CustomerUser
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cus)
		So(err, ShouldBeNil)
		So(cus, ShouldHaveSameTypeAs, []customer.CustomerUser{})

		//get customer price
		// price.CustID = c.Id
		// price.Create()
		// thyme = time.Now()
		// testThatHttp.Request("get", "/new/customer/price/", ":id", strconv.Itoa(p.ID)+"?key="+apiKey, GetCustomerPrice, nil, "")
		// So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		// So(testThatHttp.Response.Code, ShouldEqual, 200)
		// var price float64
		// err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &price)
		// So(err, ShouldBeNil)
		// So(price, ShouldHaveSameTypeAs, 7.1)

		// //get customer cart reference
		// ci.CustID = c.Id
		// ci.Create()
		// thyme = time.Now()
		// testThatHttp.Request("get", "/new/customer/cartRef/", ":id", strconv.Itoa(p.ID)+"?key="+apiKey, GetCustomerCartReference, nil, "")
		// So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		// So(testThatHttp.Response.Code, ShouldEqual, 200)
		// var reference int
		// err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &reference)
		// So(err, ShouldBeNil)
		// So(reference, ShouldHaveSameTypeAs, 7)

		//test delete customer
		thyme = time.Now()
		testThatHttp.Request("delete", "/customer/", ":id", strconv.Itoa(c.Id)+"?key=", DeleteCustomer, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, customer.Customer{})
		So(c.Id, ShouldBeGreaterThan, 0)

	})
	//cleanup
	err = cu.Delete()
	err = c.Delete()
	err = apicontextmock.DeMock(dtx)
	if err != nil {
		t.Log(err)
	}
}

// func BenchmarkCRUDCustomer(b *testing.B) {
// 	dtx, err := apicontextmock.Mock()
// 	if err != nil {
// 		b.Log(err)
// 	}

// 	qs := make(url.Values, 0)
// 	qs.Add("key", dtx.APIKey)

// 	Convey("Customer", b, func() {
// 		var c customer.Customer
// 		c.Name = "Freddy Krueger"
// 		c.Email = "freddy@elm.st"
// 		// var locs customer.CustomerLocations

// 		//create
// 		(&httprunner.BenchmarkOptions{
// 			Method:             "POST",
// 			Route:              "/customer",
// 			ParameterizedRoute: "/customer",
// 			Handler:            SaveCustomer,
// 			QueryString:        &qs,
// 			JsonBody:           c,
// 			Runs:               b.N,
// 		}).RequestBenchmark()

// 		//get
// 		(&httprunner.BenchmarkOptions{
// 			Method:             "GET",
// 			Route:              "/customer",
// 			ParameterizedRoute: "/customer/" + strconv.Itoa(c.Id),
// 			Handler:            GetCustomer,
// 			QueryString:        &qs,
// 			JsonBody:           nil,
// 			Runs:               b.N,
// 		}).RequestBenchmark()

// 		//get locations
// 		(&httprunner.BenchmarkOptions{
// 			Method:             "GET",
// 			Route:              "/customer",
// 			ParameterizedRoute: "/customer",
// 			Handler:            GetLocations,
// 			QueryString:        &qs,
// 			JsonBody:           nil,
// 			Runs:               b.N,
// 		}).RequestBenchmark()

// 		//get user
// 		(&httprunner.BenchmarkOptions{
// 			Method:             "GET",
// 			Route:              "/customer/user",
// 			ParameterizedRoute: "/customer/user",
// 			Handler:            GetUser,
// 			QueryString:        &qs,
// 			JsonBody:           nil,
// 			Runs:               b.N,
// 		}).RequestBenchmark()

// 		//get users
// 		(&httprunner.BenchmarkOptions{
// 			Method:             "GET",
// 			Route:              "/customer/users",
// 			ParameterizedRoute: "/customer/users",
// 			Handler:            GetUsers,
// 			QueryString:        &qs,
// 			JsonBody:           nil,
// 			Runs:               b.N,
// 		}).RequestBenchmark()

// 		//delete
// 		(&httprunner.BenchmarkOptions{
// 			Method:             "DELETE",
// 			Route:              "/customer",
// 			ParameterizedRoute: "/customer/" + strconv.Itoa(c.Id),
// 			Handler:            DeleteCustomer,
// 			QueryString:        &qs,
// 			JsonBody:           nil,
// 			Runs:               b.N,
// 		}).RequestBenchmark()

// 	})
// 	//teardown customer user

// 	err = apicontextmock.DeMock(dtx)
// 	if err != nil {
// 		b.Log(err)
// 	}
// }
