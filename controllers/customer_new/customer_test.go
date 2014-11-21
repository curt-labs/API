package customer_ctlr_new

import (
	"bytes"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/apiKeyType"
	"github.com/curt-labs/GoAPI/models/cartIntegration"
	"github.com/curt-labs/GoAPI/models/customer_new"
	"github.com/curt-labs/GoAPI/models/products"
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestCustomer(t *testing.T) {
	var c customer_new.Customer
	var cu customer_new.CustomerUser
	var p products.Part
	p.ID = 123
	p.Create()
	var price customer_new.Price
	price.PartID = p.ID
	price.Price = 1000000

	var ci cartIntegration.CartIntegration
	ci.PartID = p.ID
	ci.CustPartID = 987654321

	//setup apiKeyTypes
	var pub, pri, auth apiKeyType.ApiKeyType
	pub.Type = "public"
	pri.Type = "private"
	auth.Type = "authentication"
	pub.Create()
	pri.Create()
	auth.Create()

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

	Convey("Testing Customer_New/Customer", t, func() {
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
		So(c, ShouldHaveSameTypeAs, customer_new.Customer{})
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
		So(c, ShouldHaveSameTypeAs, customer_new.Customer{})
		So(c.Id, ShouldBeGreaterThan, 0)

		//test get customer
		c.JoinUser(cu)
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/", ":id", strconv.Itoa(c.Id)+"?key="+apiKey, GetCustomer, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, customer_new.Customer{})
		So(c.Id, ShouldBeGreaterThan, 0)
		//same get customer as a post request
		thyme = time.Now()
		testThatHttp.Request("post", "/new/customer/", ":id", strconv.Itoa(c.Id)+"?key="+apiKey, GetCustomer, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, customer_new.Customer{})
		So(c.Id, ShouldBeGreaterThan, 0)

		// get customer locations
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/locations", "", "?key="+apiKey, GetLocations, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c.Locations)
		So(err, ShouldBeNil)
		So(c.Locations, ShouldHaveSameTypeAs, []customer_new.CustomerLocation{})

		// get customer locations via post
		thyme = time.Now()
		testThatHttp.Request("post", "/new/customer/locations", "", "?key="+apiKey, GetLocations, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c.Locations)
		So(err, ShouldBeNil)
		So(c.Locations, ShouldHaveSameTypeAs, []customer_new.CustomerLocation{})

		//get user
		thyme = time.Now()
		testThatHttp.Request("post", "/new/customer/user", "", "?key="+apiKey, GetUser, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cu)
		So(err, ShouldBeNil)
		So(cu, ShouldHaveSameTypeAs, customer_new.CustomerUser{})

		//get users
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/users", "", "?key="+apiKey, GetUsers, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var cus []customer_new.CustomerUser
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cus)
		So(err, ShouldBeNil)
		So(cus, ShouldHaveSameTypeAs, []customer_new.CustomerUser{})

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
		So(c, ShouldHaveSameTypeAs, customer_new.Customer{})
		So(c.Id, ShouldBeGreaterThan, 0)

	})
	//cleanup
	cu.Delete()
	p.Delete()
	price.Delete()
	ci.Delete()
	pub.Delete()
	pri.Delete()
	auth.Delete()
}
