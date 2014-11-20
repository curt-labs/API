package customer_ctlr_new

import (
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

func TestCustomerUser(t *testing.T) {
	var err error
	var cu customer_new.CustomerUser
	var c customer_new.Customer
	c.Name = "Dog Bountyhunter"
	c.Create()

	var apiKey string

	Convey("Testing Customer_New/User", t, func() {
		//test create customer user
		form := url.Values{"name": {"Mitt Romney"}, "email": {"magic@underpants.com"}, "pass": {"robthepoor"}, "customerID": {strconv.Itoa(c.Id)}, "isActive": {"true"}, "locationID": {"1"}, "isSudo": {"true"}, "cust_ID": {"1"}}
		v := form.Encode()
		body := strings.NewReader(v)
		thyme := time.Now()
		testThatHttp.Request("post", "/new/customer/user/register", "", "", RegisterUser, body, "application/x-www-form-urlencoded")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cu)
		So(err, ShouldBeNil)
		So(cu, ShouldHaveSameTypeAs, customer_new.CustomerUser{})
		So(cu.Id, ShouldNotBeEmpty)
		t.Log(cu.Id)
		//key stuff - get apiKey
		for _, k := range cu.Keys {
			if strings.ToLower(k.Type) == "public" {
				apiKey = k.Key
			}
		}

		//test update customer user
		form = url.Values{"name": {"Michelle Bachman"}}
		v = form.Encode()
		body = strings.NewReader(v)
		thyme = time.Now()
		testThatHttp.Request("post", "/new/customer/user/", ":id", cu.Id, UpdateCustomerUser, body, "application/x-www-form-urlencoded")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cu)
		So(err, ShouldBeNil)
		So(cu, ShouldHaveSameTypeAs, customer_new.CustomerUser{})
		So(cu.Name, ShouldNotEqual, "Mitt Romney")

		//test authenticateUser
		c.JoinUser(cu)
		form = url.Values{"email": {"magic@underpants.com"}, "password": {"robthepoor"}}
		v = form.Encode()
		body = strings.NewReader(v)
		thyme = time.Now()
		testThatHttp.Request("post", "/new/customer/auth", "", "", AuthenticateUser, body, "application/x-www-form-urlencoded")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, customer_new.Customer{})

		//test keyed user authentication
		// thyme = time.Now()
		// testThatHttp.Request("get", "/new/customer/auth", "", "?key="+apiKey, KeyedUserAuthentication, nil, "")
		// So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		// So(testThatHttp.Response.Code, ShouldEqual, 200)
		// err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		// So(err, ShouldBeNil)
		// So(c, ShouldHaveSameTypeAs, customer_new.Customer{})

		//test delete customer users by customerId
		var cu2 customer_new.CustomerUser
		cu2.Create()
		c.JoinUser(cu2)
		thyme = time.Now()
		testThatHttp.Request("post", "/new/customer/allUsersByCustomerID/", ":id", strconv.Itoa(c.Id), DeleteCustomerUsersByCustomerID, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var response string
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &response)
		So(err, ShouldBeNil)
		So(response, ShouldHaveSameTypeAs, "this is a string")
		// c.GetUsers(apiKey)

		//test delete customer user
		thyme = time.Now()
		testThatHttp.Request("delete", "/new/customer/user/", ":id", cu.Id, DeleteCustomerUser, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cu)
		So(err, ShouldBeNil)
		So(cu, ShouldHaveSameTypeAs, customer_new.CustomerUser{})
		So(cu.Id, ShouldNotBeEmpty)
		cu2.Delete()
	})
	//teardown
	c.Delete()

	//incase
	cu.Delete()

}
