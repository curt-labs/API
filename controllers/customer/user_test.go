package customer_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/apiKeyType"
	"github.com/curt-labs/GoAPI/models/customer"
	. "github.com/smartystreets/goconvey/convey"

	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestCustomerUser(t *testing.T) {

	var err error
	var cu customer.CustomerUser
	var c customer.Customer
	c.Name = "Dog Bountyhunter"
	c.Create()

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
	Convey("Testing customer/User", t, func() {
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
		So(cu, ShouldHaveSameTypeAs, customer.CustomerUser{})
		So(cu.Id, ShouldNotBeEmpty)
		//key stuff - get apiKey
		var apiKey string
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
		So(cu, ShouldHaveSameTypeAs, customer.CustomerUser{})
		So(cu.Name, ShouldNotEqual, "Mitt Romney")

		//test authenticateUser
		err = c.JoinUser(cu)
		So(err, ShouldBeNil)
		form = url.Values{"email": {"magic@underpants.com"}, "password": {"robthepoor"}}
		v = form.Encode()
		body = strings.NewReader(v)
		thyme = time.Now()
		testThatHttp.Request("post", "/new/customer/auth", "", "", AuthenticateUser, body, "application/x-www-form-urlencoded")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, customer.Customer{})

		//test keyed user authentication
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/auth", "", "?key="+apiKey, KeyedUserAuthentication, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, customer.Customer{})

		//test get user by id
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/", ":id", cu.Id+"?key="+apiKey, GetUserById, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cu)
		So(err, ShouldBeNil)
		So(cu, ShouldHaveSameTypeAs, customer.CustomerUser{})

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
		var newKey customer.ApiCredentials
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &newKey)
		So(err, ShouldBeNil)
		So(newKey.Key, ShouldHaveSameTypeAs, "string")

		//test delete customer users by customerId
		var cu2 customer.CustomerUser
		cu2.Create([]int{1})
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
		So(cu, ShouldHaveSameTypeAs, customer.CustomerUser{})
		So(cu.Id, ShouldNotBeEmpty)
		cu2.Delete()
	})
	//teardown
	err = c.Delete()
	if err != nil {
		t.Log(err)
	}

	if database.EmptyDb != nil {
		err = pub.Delete()
		if err != nil {
			t.Log(err)
		}
		err = pri.Delete()
		if err != nil {
			t.Log(err)
		}
		err = auth.Delete()
		if err != nil {
			t.Log(err)
		}
	}

	err = cu.Delete()
	if err != nil {
		t.Log(err)
	}

}

func BenchmarkCRUDCustomerUser(b *testing.B) {
	dtx, err := apicontextmock.Mock()
	if err != nil {
		b.Log(err)
	}
	var cu customer.CustomerUser

	cu.Id = dtx.UserID

	qs := make(url.Values, 0)
	qs.Add("key", dtx.APIKey)

	Convey("CustomerUser", b, func() {
		form := url.Values{"name": {"Mitt Romney"}, "email": {"magic@underpants.com"}, "pass": {"robthepoor"}, "customerID": {strconv.Itoa(dtx.CustomerID)}, "isActive": {"true"}, "locationID": {"1"}, "isSudo": {"true"}, "cust_ID": {"1"}}
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

		//authenticate user
		form = url.Values{"email": {"magic@underpants.com"}, "password": {"robthepoor"}}
		(&httprunner.BenchmarkOptions{
			Method:             "POST",
			Route:              "/new/customer/auth",
			ParameterizedRoute: "/new/customer/auth",
			Handler:            AuthenticateUser,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()

		//authenticate user by key
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/new/customer/auth",
			ParameterizedRoute: "/new/customer/auth",
			Handler:            KeyedUserAuthentication,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
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

		form = url.Values{"email": {"magic@underpants.com"}, "password": {"robthepoor"}}

	})
	cu.Delete()

	err = apicontextmock.DeMock(dtx)
	if err != nil {
		b.Log(err)
	}
}
