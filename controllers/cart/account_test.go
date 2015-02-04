package cart_ctlr

// import (
// 	"encoding/json"
// 	"github.com/curt-labs/GoAPI/helpers/error"
// 	"github.com/curt-labs/GoAPI/helpers/httprunner"
// 	"github.com/curt-labs/GoAPI/models/cart"
// 	. "github.com/smartystreets/goconvey/convey"
// 	"net/url"
// 	"testing"
// )

// func Test_AddAccount(t *testing.T) {
// 	Convey("with no shop identifier", t, func() {
// 		qs := make(url.Values, 0)

// 		resp := httprunner.JsonRequest("POST", "/shopify/account", &qs, cart.Shop{}, AddAccount)
// 		So(resp.Code, ShouldEqual, 500)
// 		So(json.Unmarshal(resp.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
// 	})

// 	Convey("with shop identifier", t, func() {
// 		shopID := cart.InsertTestData()
// 		So(shopID, ShouldNotBeNil)
// 		val := shopID.Hex()
// 		qs := make(url.Values, 0)
// 		qs.Add("shop", val)

// 		response = httprunner.JsonRequest("POST", "/shopify/account", &qs, cart.Shop{}, AddAccount)
// 		So(response.Code, ShouldEqual, 500)
// 		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

// 		cust := cart.Customer{
// 			ShopId:    *shopID,
// 			FirstName: "Alex",
// 			LastName:  "Ninneman",
// 		}

// 		response = httprunner.JsonRequest("POST", "/shopify/account", &qs, cust, AddAccount)
// 		So(response.Code, ShouldEqual, 500)
// 		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

// 		cust.Email = "ninnemana@gmail.com"
// 		response = httprunner.JsonRequest("POST", "/shopify/account", &qs, cust, AddAccount)
// 		So(response.Code, ShouldEqual, 500)
// 		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

// 		cust.Password = "password"
// 		response = httprunner.JsonRequest("POST", "/shopify/account", &qs, cust, AddAccount)
// 		So(response.Code, ShouldEqual, 200)
// 		So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)
// 	})
// }

// func Test_GetAccount(t *testing.T) {
// 	Convey("no shop identifier", t, func() {
// 		response = httprunner.Request("GET", "/shopify/account", nil, GetAccount)
// 		So(response.Code, ShouldEqual, 500)
// 		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

// 		vals := make(url.Values, 0)
// 		vals.Add("shop", "testing")
// 		response = httprunner.Request("GET", "/shopify/account", &vals, GetAccount)
// 		So(response.Code, ShouldEqual, 500)
// 		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
// 	})
// 	Convey("with shop identifier", t, func() {
// 		shopID := cart.InsertTestData()
// 		So(shopID, ShouldNotBeNil)

// 		val := shopID.Hex()
// 		qs := make(url.Values, 0)
// 		qs.Add("shop", val)

// 		Convey("with bad accound reference", func() {
// 			response = httprunner.ParameterizedRequest("GET", "/shopify/account", "/shopify/account", &qs, nil, GetAccount)
// 			So(response.Code, ShouldEqual, 500)
// 			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
// 		})

// 		Convey("with good account reference", func() {
// 			cust := cart.Customer{
// 				ShopId:    *shopID,
// 				FirstName: "Alex",
// 				LastName:  "Ninneman",
// 				Email:     "ninnemana@gmail.com",
// 				Password:  "password",
// 			}
// 			response = httprunner.Req(AddAccount, "GET", "", "/shopify/account", &qs, nil, cust)
// 			// So(response.Code, ShouldEqual, 200)
// 			// So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)

// 			// response = httprunner.JsonRequest("POST", "/shopify/account", &qs, cust, AddAccount)
// 			// So(response.Code, ShouldEqual, 200)
// 			// So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)

// 			// response = httprunner.ParameterizedRequest("GET", "/shopify/account", "/shopify/customers/"+shopID.Hex(), &qs, nil, GetCustomer)
// 			// So(response.Code, ShouldEqual, 500)
// 			// So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

// 			// response = httprunner.ParameterizedRequest("GET", "/shopify/customers/:id", "/shopify/customers/"+cust.Id.Hex(), &qs, nil, GetCustomer)
// 			// So(response.Code, ShouldEqual, 200)
// 			// So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)
// 		})
// 	})
// }

// func Test_AccountLogin(t *testing.T) {
// 	Convey("no shop identifier", t, func() {
// 		response = httprunner.Request("GET", "/shopify/customers/1234", nil, GetCustomer)
// 		So(response.Code, ShouldEqual, 500)
// 		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

// 		vals := make(url.Values, 0)
// 		vals.Add("shop", "testing")
// 		response = httprunner.Request("GET", "/shopify/customers/1234", &vals, GetCustomer)
// 		So(response.Code, ShouldEqual, 500)
// 		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
// 	})
// 	Convey("with shop identifier", t, func() {
// 		shopID := cart.InsertTestData()
// 		So(shopID, ShouldNotBeNil)

// 		val := shopID.Hex()
// 		qs := make(url.Values, 0)
// 		qs.Add("shop", val)

// 		Convey("with bad customer reference", func() {
// 			response = httprunner.ParameterizedRequest("GET", "/shopify/customers/:id", "/shopify/customers/1234", &qs, nil, GetCustomer)
// 			So(response.Code, ShouldEqual, 500)
// 			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
// 		})

// 		Convey("with good customer reference", func() {
// 			cust := cart.Customer{
// 				ShopId:    *shopID,
// 				FirstName: "Alex",
// 				LastName:  "Ninneman",
// 				Email:     "ninnemana@gmail.com",
// 				Password:  "password",
// 			}
// 			response = httprunner.JsonRequest("POST", "/shopify/customers", &qs, cust, AddCustomer)
// 			So(response.Code, ShouldEqual, 200)
// 			So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)

// 			response = httprunner.JsonRequest("POST", "/shopify/customers/login", &qs, cart.CustomerAddress{}, AccountLogin)
// 			So(response.Code, ShouldEqual, 500)
// 			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

// 			cust.Email = ""
// 			cust.Password = ""
// 			response = httprunner.JsonRequest("POST", "/shopify/customers/login", &qs, cust, AccountLogin)
// 			So(response.Code, ShouldEqual, 500)
// 			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

// 			cust.Email = "ninnemana@gmail.com"
// 			response = httprunner.JsonRequest("POST", "/shopify/customers/login", &qs, cust, AccountLogin)
// 			So(response.Code, ShouldEqual, 500)
// 			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

// 			cust.Email = "ninnemana@gmail.com"
// 			cust.Password = "bad_password"
// 			response = httprunner.JsonRequest("POST", "/shopify/customers/login", &qs, cust, AccountLogin)
// 			So(response.Code, ShouldEqual, 500)
// 			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

// 			cust.Email = "ninnemana@gmail.com"
// 			cust.Password = "password"
// 			response = httprunner.JsonRequest("POST", "/shopify/customers/login", &qs, cust, AccountLogin)
// 			So(response.Code, ShouldEqual, 200)
// 			var c cart.Customer
// 			So(json.Unmarshal(response.Body.Bytes(), &c), ShouldBeNil)
// 			So(c.Password, ShouldEqual, "")
// 		})
// 	})
// }

// func Test_EditAccount(t *testing.T) {
// 	Convey("no shop identifier", t, func() {
// 		qs := make(url.Values, 0)
// 		response = httprunner.JsonRequest("PUT", "/shopify/customers/1234", &qs, &cart.Customer{}, EditCustomer)
// 		So(response.Code, ShouldEqual, 500)
// 		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

// 		qs.Add("shop", "testing")
// 		response = httprunner.JsonRequest("PUT", "/shopify/customers/1234", &qs, &cart.Customer{}, EditCustomer)
// 		So(response.Code, ShouldEqual, 500)
// 		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
// 	})
// 	Convey("with shop identifier", t, func() {
// 		shopID := cart.InsertTestData()
// 		So(shopID, ShouldNotBeNil)

// 		val := shopID.Hex()
// 		qs := make(url.Values, 0)
// 		qs.Add("shop", val)

// 		Convey("with bad customer reference", func() {
// 			response = httprunner.ParameterizedJsonRequest("PUT", "/shopify/customers/:id", "/shopify/customers/1234", &qs, &cart.Customer{}, EditCustomer)
// 			So(response.Code, ShouldEqual, 500)
// 			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
// 		})

// 		Convey("with good customer reference", func() {
// 			cust := cart.Customer{
// 				ShopId:    *shopID,
// 				FirstName: "Alex",
// 				LastName:  "Ninneman",
// 				Email:     "ninnemana@gmail.com",
// 				Password:  "password",
// 			}
// 			response = httprunner.JsonRequest("POST", "/shopify/customers", &qs, cust, AddCustomer)
// 			So(response.Code, ShouldEqual, 200)
// 			So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)

// 			response = httprunner.ParameterizedJsonRequest("PUT", "/shopify/customers/:id", "/shopify/customers/"+cust.Id.Hex(), &qs, &cart.Shop{}, EditCustomer)
// 			So(response.Code, ShouldEqual, 500)
// 			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

// 			cust.Email = ""
// 			response = httprunner.ParameterizedJsonRequest("PUT", "/shopify/customers/:id", "/shopify/customers/"+cust.Id.Hex(), &qs, cust, EditCustomer)
// 			So(response.Code, ShouldEqual, 500)
// 			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

// 			cust.Email = "alex@ninneman.org"
// 			response = httprunner.ParameterizedJsonRequest("PUT", "/shopify/customers/:id", "/shopify/customers/"+cust.Id.Hex(), &qs, cust, EditCustomer)
// 			So(response.Code, ShouldEqual, 200)
// 			So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)
// 		})
// 	})
// }

// func Test_DeleteAccount(t *testing.T) {
// 	Convey("no shop identifier", t, func() {
// 		qs := make(url.Values, 0)
// 		response = httprunner.JsonRequest("DELETE", "/shopify/customers/1234", &qs, &cart.Customer{}, DeleteCustomer)
// 		So(response.Code, ShouldEqual, 500)
// 		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

// 		qs.Add("shop", "testing")
// 		response = httprunner.JsonRequest("DELETE", "/shopify/customers/1234", &qs, &cart.Customer{}, DeleteCustomer)
// 		So(response.Code, ShouldEqual, 500)
// 		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
// 	})
// 	Convey("with shop identifier", t, func() {
// 		shopID := cart.InsertTestData()
// 		So(shopID, ShouldNotBeNil)

// 		val := shopID.Hex()
// 		qs := make(url.Values, 0)
// 		qs.Add("shop", val)

// 		Convey("with bad customer reference", func() {
// 			response = httprunner.ParameterizedJsonRequest("DELETE", "/shopify/customers/:id", "/shopify/customers/1234", &qs, &cart.Customer{}, DeleteCustomer)
// 			So(response.Code, ShouldEqual, 500)
// 			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
// 		})

// 		Convey("with good customer reference", func() {
// 			cust := cart.Customer{
// 				ShopId:    *shopID,
// 				FirstName: "Alex",
// 				LastName:  "Ninneman",
// 				Email:     "ninnemana@gmail.com",
// 				Password:  "password",
// 			}
// 			response = httprunner.JsonRequest("POST", "/shopify/customers", &qs, cust, AddCustomer)
// 			So(response.Code, ShouldEqual, 200)
// 			So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)

// 			response = httprunner.ParameterizedJsonRequest("DELETE", "/shopify/customers/:id", "/shopify/customers/"+shopID.Hex(), &qs, cust, DeleteCustomer)
// 			So(response.Code, ShouldEqual, 500)
// 			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

// 			response = httprunner.ParameterizedJsonRequest("DELETE", "/shopify/customers/:id", "/shopify/customers/"+cust.Id.Hex(), &qs, cust, DeleteCustomer)
// 			So(response.Code, ShouldEqual, 200)
// 			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldNotBeNil)
// 		})
// 	})
// }

// func Test_GetAccountOrders(t *testing.T) {
// 	Convey("no shop identifier", t, func() {
// 		response = httprunner.Request("GET", "/shopify/customers/1234/orders", nil, GetCustomerOrders)
// 		So(response.Code, ShouldEqual, 500)
// 		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

// 		vals := make(url.Values, 0)
// 		vals.Add("shop", "testing")
// 		response = httprunner.Request("GET", "/shopify/customers/1234/orders", &vals, GetCustomerOrders)
// 		So(response.Code, ShouldEqual, 500)
// 		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
// 	})
// 	Convey("with shop identifier", t, func() {
// 		shopID := cart.InsertTestData()
// 		So(shopID, ShouldNotBeNil)

// 		val := shopID.Hex()
// 		qs := make(url.Values, 0)
// 		qs.Add("shop", val)

// 		Convey("with bad customer reference", func() {
// 			response = httprunner.ParameterizedRequest("GET", "/shopify/customers/:id/orders", "/shopify/customers/1234/orders", &qs, nil, GetCustomerOrders)
// 			So(response.Code, ShouldEqual, 500)
// 			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
// 		})

// 		Convey("with good customer reference", func() {
// 			cust := cart.Customer{
// 				ShopId:    *shopID,
// 				FirstName: "Alex",
// 				LastName:  "Ninneman",
// 				Email:     "ninnemana@gmail.com",
// 				Password:  "password",
// 			}
// 			response = httprunner.JsonRequest("POST", "/shopify/customers", &qs, cust, AddCustomer)
// 			So(response.Code, ShouldEqual, 200)
// 			So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)

// 			response = httprunner.ParameterizedRequest("GET", "/shopify/customers/:id/orders", "/shopify/customers/"+shopID.Hex()+"/orders", &qs, nil, GetCustomerOrders)
// 			So(response.Code, ShouldEqual, 500)
// 			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

// 			var orders []interface{}
// 			response = httprunner.ParameterizedRequest("GET", "/shopify/customers/:id/orders", "/shopify/customers/"+cust.Id.Hex()+"/orders", &qs, nil, GetCustomerOrders)
// 			So(response.Code, ShouldEqual, 200)
// 			So(json.Unmarshal(response.Body.Bytes(), &orders), ShouldBeNil)
// 		})
// 	})
// }

// func BenchmarkAddAccount(b *testing.B) {
// 	shopID := cart.InsertTestData()
// 	if shopID == nil {
// 		b.Error("failed to create a shop")
// 		b.Fail()
// 	}

// 	val := shopID.Hex()
// 	qs := make(url.Values, 0)
// 	qs.Add("shop", val)

// 	cust := cart.Customer{
// 		ShopId:    *shopID,
// 		FirstName: "Alex",
// 		LastName:  "Ninneman",
// 		Email:     "ninnemana@gmail.com",
// 		Password:  "password",
// 	}

// 	cust.Email = "ninnemana@gmail.com"
// 	(&httprunner.BenchmarkOptions{
// 		Method:             "POST",
// 		Route:              "/shopify/customers",
// 		ParameterizedRoute: "/shopify/customers",
// 		Handler:            AddCustomer,
// 		QueryString:        &qs,
// 		JsonBody:           cust,
// 		Runs:               b.N,
// 	}).RequestBenchmark()
// }

// func BenchmarkGetAcount(b *testing.B) {
// 	shopID := cart.InsertTestData()
// 	if shopID == nil {
// 		b.Error("failed to create a shop")
// 		b.Fail()
// 	}

// 	val := shopID.Hex()
// 	qs := make(url.Values, 0)
// 	qs.Add("shop", val)

// 	cust := cart.Customer{
// 		ShopId:    *shopID,
// 		FirstName: "Alex",
// 		LastName:  "Ninneman",
// 		Email:     "ninnemana@gmail.com",
// 		Password:  "password",
// 	}
// 	if err := cust.Insert("http://www.example.com"); err != nil {
// 		b.Error(err.Error())
// 		b.Fail()
// 	}

// 	(&httprunner.BenchmarkOptions{
// 		Method:             "GET",
// 		Route:              "/shopify/customers/" + cust.Id.Hex(),
// 		ParameterizedRoute: "/shopify/customers/:id",
// 		Handler:            GetCustomers,
// 		QueryString:        &qs,
// 		JsonBody:           nil,
// 		Runs:               b.N,
// 	}).RequestBenchmark()
// }

// func BenchmarkEditAccount(b *testing.B) {
// 	shopID := cart.InsertTestData()
// 	if shopID == nil {
// 		b.Error("failed to create a shop")
// 		b.Fail()
// 	}

// 	val := shopID.Hex()
// 	qs := make(url.Values, 0)
// 	qs.Add("shop", val)

// 	cust := cart.Customer{
// 		ShopId:    *shopID,
// 		FirstName: "Alex",
// 		LastName:  "Ninneman",
// 		Email:     "ninnemana@gmail.com",
// 		Password:  "password",
// 	}

// 	if err := cust.Insert("http://www.example.com"); err != nil {
// 		b.Error(err.Error())
// 		b.Fail()
// 	}

// 	cust.Email = "alex@ninneman.org"
// 	(&httprunner.BenchmarkOptions{
// 		Method:             "PUT",
// 		Route:              "/shopify/customers/" + cust.Id.Hex(),
// 		ParameterizedRoute: "/shopify/customers/" + cust.Id.Hex(),
// 		Handler:            EditCustomer,
// 		QueryString:        &qs,
// 		JsonBody:           cust,
// 		Runs:               b.N,
// 	}).RequestBenchmark()
// }
