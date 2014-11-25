package customer_ctlr_new

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/apiKeyType"
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

	//setup apiKeyTypes
	var pub, pri, auth apiKeyType.ApiKeyType
	pub.Type = "public"
	pri.Type = "private"
	auth.Type = "authentication"
	pub.Create()
	pri.Create()
	auth.Create()

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
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/auth", "", "?key="+apiKey, KeyedUserAuthentication, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, customer_new.Customer{})

		//test get user by id
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/", ":id", cu.Id+"?key="+apiKey, GetUserById, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cu)
		So(err, ShouldBeNil)
		So(cu, ShouldHaveSameTypeAs, customer_new.CustomerUser{})

		//test change user password
		form = url.Values{"email": {"magic@underpants.com"}, "oldPass": {"robthepoor"}, "newPass": {"prolife"}}
		v = form.Encode()
		body = strings.NewReader(v)
		thyme = time.Now()
		testThatHttp.Request("post", "/new/customer/user/changePassword", "", "?key="+apiKey, ChangePassword, body, "application/x-www-form-urlencoded")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var result string
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &result)
		So(err, ShouldBeNil)
		So(result, ShouldHaveSameTypeAs, "Success")

		//test reset  user password
		form = url.Values{"email": {"magic@underpants.com"}, "customerID": {strconv.Itoa(c.CustomerId)}}
		v = form.Encode()
		body = strings.NewReader(v)
		thyme = time.Now()
		testThatHttp.Request("post", "/new/customer/user/resetPassword", "", "?key="+apiKey, ResetPassword, body, "application/x-www-form-urlencoded")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &result)
		So(err, ShouldBeNil)
		So(result, ShouldHaveSameTypeAs, "Success")

		//test generate api key
		thyme = time.Now()
		testThatHttp.Request("post", "/new/customer/user/", ":id/key/:type", cu.Id+"/key/PRIVATE?key="+apiKey, GenerateApiKey, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var newKey customer_new.ApiCredentials
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &newKey)
		So(err, ShouldBeNil)
		So(newKey.Key, ShouldHaveSameTypeAs, "string")

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
	pub.Delete()
	pri.Delete()
	auth.Delete()

	//incase
	cu.Delete()

}

func BenchmarkCRUDCustomerUser(b *testing.B) {
	var cu customer_new.CustomerUser
	var c customer_new.Customer
	c.Name = "Mick Mattleson"
	c.Create()

	qs := make(url.Values, 0)

	Convey("CustomerUser", b, func() {
		form := url.Values{"name": {"Mitt Romney"}, "email": {"magic@underpants.com"}, "pass": {"robthepoor"}, "customerID": {strconv.Itoa(c.Id)}, "isActive": {"true"}, "locationID": {"1"}, "isSudo": {"true"}, "cust_ID": {"1"}}
		//create
		(&httprunner.BenchmarkOptions{
			Method:             "POST",
			Route:              "/new/customer/user/register",
			ParameterizedRoute: "/new/customer/user/register",
			Handler:            RegisterUser,
			QueryString:        &qs,
			JsonBody:           cu,
			FormBody:           form,
			Runs:               b.N,
		}).RequestBenchmark()

		//delete
		(&httprunner.BenchmarkOptions{
			Method:             "DELETE",
			Route:              "/new/customer/user/register",
			ParameterizedRoute: "/new/customer/user/register/" + cu.Id,
			Handler:            DeleteCustomerUser,
			QueryString:        &qs,
			JsonBody:           cu,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
	})
	cu.Delete()
	c.Delete()
}
