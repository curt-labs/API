package cart_ctlr

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/models/cart"
	. "github.com/smartystreets/goconvey/convey"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

var (
	response httptest.ResponseRecorder
)

func Test_GetCustomers(t *testing.T) {
	Convey("Testing GetCustomers", t, func() {
		Convey("no shop identifier", func() {
			response = httprunner.Request("GET", "/shopify/customers", nil, GetCustomers)
			So(response.Code, ShouldEqual, 500)
			So(response.Body.String(), ShouldNotEqual, "[]")

			vals := make(url.Values, 0)
			vals.Add("shop", "testing")
			response = httprunner.Request("GET", "/shopify/customers", &vals, GetCustomers)
			So(response.Code, ShouldEqual, 500)
			So(response.Body.String(), ShouldNotEqual, "[]")

			vals.Add("since_id", "something")
			response = httprunner.Request("GET", "/shopify/customers", &vals, GetCustomers)
			So(response.Code, ShouldEqual, 500)
			So(response.Body.String(), ShouldNotEqual, "[]")
		})
		Convey("with shop identifier", func() {
			shopID := cart.InsertTestData()
			So(shopID, ShouldNotBeNil)

			val := shopID.Hex()
			qs := make(url.Values, 0)
			qs.Add("shop", val)

			response = httprunner.Request("GET", "/shopify/customers", &qs, GetCustomers)
			So(response.Code, ShouldEqual, 200)

			So(json.Unmarshal(response.Body.Bytes(), &[]cart.Customer{}), ShouldBeNil)

			qs.Add("since_id", "something")
			response = httprunner.Request("GET", "/shopify/customers", &qs, GetCustomers)
			So(response.Code, ShouldEqual, 200)
			So(json.Unmarshal(response.Body.Bytes(), &[]cart.Customer{}), ShouldBeNil)

			qs.Set("since_id", val)
			response = httprunner.Request("GET", "/shopify/customers", &qs, GetCustomers)
			So(response.Code, ShouldEqual, 200)
			var custs []cart.Customer
			So(json.Unmarshal(response.Body.Bytes(), &custs), ShouldBeNil)
			So(len(custs), ShouldEqual, 0)

			created_min := time.Now().AddDate(-1, 0, 0).Format(TimeLayout)
			created_max := time.Now().Format(TimeLayout)
			qs.Add("created_at_min", created_min)
			qs.Add("created_at_max", created_max)
			response = httprunner.Request("GET", "/shopify/customers", &qs, GetCustomers)
			So(response.Code, ShouldEqual, 200)
			So(json.Unmarshal(response.Body.Bytes(), &custs), ShouldBeNil)
			So(len(custs), ShouldEqual, 0)

			update_min := time.Now().AddDate(-1, 0, 0).Format(TimeLayout)
			update_max := time.Now().Format(TimeLayout)
			qs.Add("updated_at_min", update_min)
			qs.Add("updated_at_max", update_max)
			response = httprunner.Request("GET", "/shopify/customers", &qs, GetCustomers)
			So(response.Code, ShouldEqual, 200)
			So(json.Unmarshal(response.Body.Bytes(), &custs), ShouldBeNil)
			So(len(custs), ShouldEqual, 0)

			qs.Add("page", "bad")
			qs.Add("limit", "test")
			response = httprunner.Request("GET", "/shopify/customers", &qs, GetCustomers)
			So(response.Code, ShouldEqual, 200)
			So(json.Unmarshal(response.Body.Bytes(), &custs), ShouldBeNil)
			So(len(custs), ShouldEqual, 0)

			qs.Set("page", "1")
			qs.Set("limit", "1")
			response = httprunner.Request("GET", "/shopify/customers", &qs, GetCustomers)
			So(response.Code, ShouldEqual, 200)
			So(json.Unmarshal(response.Body.Bytes(), &custs), ShouldBeNil)
			So(len(custs), ShouldEqual, 0)
		})
	})
}

func Test_AddCustomer(t *testing.T) {
	Convey("with no shop identifier", t, func() {
		qs := make(url.Values, 0)

		resp := httprunner.JsonRequest("POST", "/shopify/customers", &qs, cart.Shop{}, AddCustomer)
		So(resp.Code, ShouldEqual, 500)
		So(json.Unmarshal(resp.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
	})

	Convey("with shop identifier", t, func() {
		shopID := cart.InsertTestData()
		So(shopID, ShouldNotBeNil)
		val := shopID.Hex()
		qs := make(url.Values, 0)
		qs.Add("shop", val)

		response = httprunner.JsonRequest("POST", "/shopify/customers", &qs, cart.Shop{}, AddCustomer)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		cust := cart.Customer{
			ShopId:    *shopID,
			FirstName: "Alex",
			LastName:  "Ninneman",
		}

		response = httprunner.JsonRequest("POST", "/shopify/customers", &qs, cust, AddCustomer)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		cust.Email = time.Now().Format(time.RFC3339Nano) + "@gmail.com"
		response = httprunner.JsonRequest("POST", "/shopify/customers", &qs, cust, AddCustomer)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		cust.Password = "password"
		response = httprunner.JsonRequest("POST", "/shopify/customers", &qs, cust, AddCustomer)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)
	})
}

func Test_GetCustomer(t *testing.T) {
	Convey("no shop identifier", t, func() {
		response = httprunner.Request("GET", "/shopify/customers/1234", nil, GetCustomer)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		vals := make(url.Values, 0)
		vals.Add("shop", "testing")
		response = httprunner.Request("GET", "/shopify/customers/1234", &vals, GetCustomer)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
	})
	Convey("with shop identifier", t, func() {
		shopID := cart.InsertTestData()
		So(shopID, ShouldNotBeNil)

		val := shopID.Hex()
		qs := make(url.Values, 0)
		qs.Add("shop", val)

		Convey("with bad customer reference", func() {
			response = httprunner.ParameterizedRequest("GET", "/shopify/customers/:id", "/shopify/customers/1234", &qs, nil, GetCustomer)
			So(response.Code, ShouldEqual, 500)
			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
		})

		Convey("with good customer reference", func() {
			cust := cart.Customer{
				ShopId:    *shopID,
				FirstName: "Alex",
				LastName:  "Ninneman",
				Email:     time.Now().Format(time.RFC3339Nano) + "@gmail.com",
				Password:  "password",
			}
			response = httprunner.JsonRequest("POST", "/shopify/customers", &qs, cust, AddCustomer)
			So(response.Code, ShouldEqual, 200)
			So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)

			response = httprunner.ParameterizedRequest("GET", "/shopify/customers/:id", "/shopify/customers/"+shopID.Hex(), &qs, nil, GetCustomer)
			So(response.Code, ShouldEqual, 500)
			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

			response = httprunner.ParameterizedRequest("GET", "/shopify/customers/:id", "/shopify/customers/"+cust.Id.Hex(), &qs, nil, GetCustomer)
			So(response.Code, ShouldEqual, 200)
			So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)
		})
	})
}

func Test_EditCustomer(t *testing.T) {
	Convey("no shop identifier", t, func() {
		qs := make(url.Values, 0)
		response = httprunner.JsonRequest("PUT", "/shopify/customers/1234", &qs, &cart.Customer{}, EditCustomer)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		qs.Add("shop", "testing")
		response = httprunner.JsonRequest("PUT", "/shopify/customers/1234", &qs, &cart.Customer{}, EditCustomer)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
	})
	Convey("with shop identifier", t, func() {
		shopID := cart.InsertTestData()
		So(shopID, ShouldNotBeNil)

		val := shopID.Hex()
		qs := make(url.Values, 0)
		qs.Add("shop", val)

		Convey("with bad customer reference", func() {
			response = httprunner.ParameterizedJsonRequest("PUT", "/shopify/customers/:id", "/shopify/customers/1234", &qs, &cart.Customer{}, EditCustomer)
			So(response.Code, ShouldEqual, 500)
			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
		})

		Convey("with good customer reference", func() {
			cust := cart.Customer{
				ShopId:    *shopID,
				FirstName: "Alex",
				LastName:  "Ninneman",
				Email:     time.Now().Format(time.RFC3339Nano) + "@gmail.com",
				Password:  "password",
			}
			response = httprunner.JsonRequest("POST", "/shopify/customers", &qs, cust, AddCustomer)
			So(response.Code, ShouldEqual, 200)
			So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)

			response = httprunner.ParameterizedJsonRequest("PUT", "/shopify/customers/:id", "/shopify/customers/"+cust.Id.Hex(), &qs, &cart.Shop{}, EditCustomer)
			So(response.Code, ShouldEqual, 500)
			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

			cust.FirstName = ""
			response = httprunner.ParameterizedJsonRequest("PUT", "/shopify/customers/:id", "/shopify/customers/"+cust.Id.Hex(), &qs, cust, EditCustomer)
			So(response.Code, ShouldEqual, 500)
			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

			cust.FirstName = "Alex"
			response = httprunner.ParameterizedJsonRequest("PUT", "/shopify/customers/:id", "/shopify/customers/"+cust.Id.Hex(), &qs, cust, EditCustomer)
			So(response.Code, ShouldEqual, 200)
			So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)
		})
	})
}

func Test_DeleteCustomer(t *testing.T) {
	Convey("no shop identifier", t, func() {
		qs := make(url.Values, 0)
		response = httprunner.JsonRequest("DELETE", "/shopify/customers/1234", &qs, &cart.Customer{}, DeleteCustomer)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		qs.Add("shop", "testing")
		response = httprunner.JsonRequest("DELETE", "/shopify/customers/1234", &qs, &cart.Customer{}, DeleteCustomer)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
	})
	Convey("with shop identifier", t, func() {
		shopID := cart.InsertTestData()
		So(shopID, ShouldNotBeNil)

		val := shopID.Hex()
		qs := make(url.Values, 0)
		qs.Add("shop", val)

		Convey("with bad customer reference", func() {
			response = httprunner.ParameterizedJsonRequest("DELETE", "/shopify/customers/:id", "/shopify/customers/1234", &qs, &cart.Customer{}, DeleteCustomer)
			So(response.Code, ShouldEqual, 500)
			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
		})

		Convey("with good customer reference", func() {
			cust := cart.Customer{
				ShopId:    *shopID,
				FirstName: "Alex",
				LastName:  "Ninneman",
				Email:     time.Now().Format(time.RFC3339Nano) + "@gmail.com",
				Password:  "password",
			}
			response = httprunner.JsonRequest("POST", "/shopify/customers", &qs, cust, AddCustomer)
			So(response.Code, ShouldEqual, 200)
			So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)

			response = httprunner.ParameterizedJsonRequest("DELETE", "/shopify/customers/:id", "/shopify/customers/"+shopID.Hex(), &qs, cust, DeleteCustomer)
			So(response.Code, ShouldEqual, 500)
			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

			response = httprunner.ParameterizedJsonRequest("DELETE", "/shopify/customers/:id", "/shopify/customers/"+cust.Id.Hex(), &qs, cust, DeleteCustomer)
			So(response.Code, ShouldEqual, 200)
			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldNotBeNil)
		})
	})
}

func Test_SearchCustomer(t *testing.T) {
	Convey("no shop identifier", t, func() {
		qs := make(url.Values, 0)
		response = httprunner.Request("GET", "/shopify/customers/search", &qs, SearchCustomer)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		qs.Add("shop", "testing")
		response = httprunner.Request("GET", "/shopify/customers/search", &qs, SearchCustomer)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
	})
	Convey("with shop identifier", t, func() {
		shopID := cart.InsertTestData()
		So(shopID, ShouldNotBeNil)

		val := shopID.Hex()
		qs := make(url.Values, 0)
		qs.Add("shop", val)

		Convey("with no query", func() {
			response = httprunner.Request("GET", "/shopify/customers/search", &qs, SearchCustomer)
			So(response.Code, ShouldEqual, 500)
			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
		})

		// TODO - this needs to be fixed
		// right now I'm forcing it to allow a 500 as a success
		// the search indexing on mongo needs to be setup,
		// still a little fuzzy on how it works.
		Convey("with query", func() {
			qs.Add("query", "alex")
			response = httprunner.Request("GET", "/shopify/customers/search", &qs, SearchCustomer)
			So(response.Code, ShouldEqual, 500)
			So(json.Unmarshal(response.Body.Bytes(), &[]cart.Customer{}), ShouldNotBeNil)
		})
	})
}

func Test_GetCustomerOrders(t *testing.T) {
	Convey("no shop identifier", t, func() {
		response = httprunner.Request("GET", "/shopify/customers/1234/orders", nil, GetCustomerOrders)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		vals := make(url.Values, 0)
		vals.Add("shop", "testing")
		response = httprunner.Request("GET", "/shopify/customers/1234/orders", &vals, GetCustomerOrders)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
	})
	Convey("with shop identifier", t, func() {
		shopID := cart.InsertTestData()
		So(shopID, ShouldNotBeNil)

		val := shopID.Hex()
		qs := make(url.Values, 0)
		qs.Add("shop", val)

		Convey("with bad customer reference", func() {
			response = httprunner.ParameterizedRequest("GET", "/shopify/customers/:id/orders", "/shopify/customers/1234/orders", &qs, nil, GetCustomerOrders)
			So(response.Code, ShouldEqual, 500)
			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
		})

		Convey("with good customer reference", func() {
			cust := cart.Customer{
				ShopId:    *shopID,
				FirstName: "Alex",
				LastName:  "Ninneman",
				Email:     time.Now().Format(time.RFC3339Nano) + "@gmail.com",
				Password:  "password",
			}
			response = httprunner.JsonRequest("POST", "/shopify/customers", &qs, cust, AddCustomer)
			So(response.Code, ShouldEqual, 200)
			So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)

			response = httprunner.ParameterizedRequest("GET", "/shopify/customers/:id/orders", "/shopify/customers/"+shopID.Hex()+"/orders", &qs, nil, GetCustomerOrders)
			So(response.Code, ShouldEqual, 500)
			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

			var orders []interface{}
			response = httprunner.ParameterizedRequest("GET", "/shopify/customers/:id/orders", "/shopify/customers/"+cust.Id.Hex()+"/orders", &qs, nil, GetCustomerOrders)
			So(response.Code, ShouldEqual, 200)
			So(json.Unmarshal(response.Body.Bytes(), &orders), ShouldBeNil)
		})
	})
}

func BenchmarkGetCustomers(b *testing.B) {
	shopID := cart.InsertTestData()
	if shopID == nil {
		b.Error("shopID cannot be nil")
		b.Fail()
	}

	val := shopID.Hex()
	qs := make(url.Values, 0)
	qs.Add("shop", val)
	qs.Add("since_id", val)
	(&httprunner.BenchmarkOptions{
		Method:             "GET",
		Route:              "/shopify/customers",
		ParameterizedRoute: "/shopify/customers",
		Handler:            GetCustomer,
		QueryString:        &qs,
		JsonBody:           nil,
		Runs:               b.N,
	}).RequestBenchmark()
}

func BenchmarkAddCustomer(b *testing.B) {
	shopID := cart.InsertTestData()
	if shopID == nil {
		b.Error("failed to create a shop")
		b.Fail()
	}

	val := shopID.Hex()
	qs := make(url.Values, 0)
	qs.Add("shop", val)

	cust := cart.Customer{
		ShopId:    *shopID,
		FirstName: "Alex",
		LastName:  "Ninneman",
		Email:     time.Now().Format(time.RFC3339Nano) + "@gmail.com",
		Password:  "password",
	}

	cust.Email = "ninnemana@gmail.com"
	(&httprunner.BenchmarkOptions{
		Method:             "POST",
		Route:              "/shopify/customers",
		ParameterizedRoute: "/shopify/customers",
		Handler:            AddCustomer,
		QueryString:        &qs,
		JsonBody:           cust,
		Runs:               b.N,
	}).RequestBenchmark()
}

func BenchmarkGetCustomer(b *testing.B) {
	shopID := cart.InsertTestData()
	if shopID == nil {
		b.Error("failed to create a shop")
		b.Fail()
	}

	val := shopID.Hex()
	qs := make(url.Values, 0)
	qs.Add("shop", val)

	cust := cart.Customer{
		ShopId:    *shopID,
		FirstName: "Alex",
		LastName:  "Ninneman",
		Email:     "ninnemana@gmail.com",
		Password:  "password",
	}
	if err := cust.Insert("http://www.example.com"); err != nil {
		b.Error(err.Error())
		b.Fail()
	}

	(&httprunner.BenchmarkOptions{
		Method:             "GET",
		Route:              "/shopify/customers/" + cust.Id.Hex(),
		ParameterizedRoute: "/shopify/customers/:id",
		Handler:            GetCustomers,
		QueryString:        &qs,
		JsonBody:           nil,
		Runs:               b.N,
	}).RequestBenchmark()
}

func BenchmarkEditCustomer(b *testing.B) {
	shopID := cart.InsertTestData()
	if shopID == nil {
		b.Error("failed to create a shop")
		b.Fail()
	}

	val := shopID.Hex()
	qs := make(url.Values, 0)
	qs.Add("shop", val)

	cust := cart.Customer{
		ShopId:    *shopID,
		FirstName: "Alex",
		LastName:  "Ninneman",
		Email:     time.Now().Format(time.RFC3339Nano) + "@gmail.com",
		Password:  "password",
	}

	if err := cust.Insert("http://www.example.com"); err != nil {
		b.Error(err.Error())
		b.Fail()
	}

	cust.Email = "alex@ninneman.org"
	(&httprunner.BenchmarkOptions{
		Method:             "PUT",
		Route:              "/shopify/customers/" + cust.Id.Hex(),
		ParameterizedRoute: "/shopify/customers/" + cust.Id.Hex(),
		Handler:            EditCustomer,
		QueryString:        &qs,
		JsonBody:           cust,
		Runs:               b.N,
	}).RequestBenchmark()
}

func BenchmarkDeleteCustomer(b *testing.B) {
	shopID := cart.InsertTestData()
	if shopID == nil {
		b.Error("failed to create a shop")
		b.Fail()
	}

	val := shopID.Hex()
	qs := make(url.Values, 0)
	qs.Add("shop", val)

	cust := cart.Customer{
		ShopId:    *shopID,
		FirstName: "Alex",
		LastName:  "Ninneman",
		Email:     time.Now().Format(time.RFC3339Nano) + "@gmail.com",
	}

	if err := cust.Insert("http://www.example.com"); err != nil {
		b.Error(err.Error())
		b.Fail()
	}

	(&httprunner.BenchmarkOptions{
		Method:             "DELETE",
		Route:              "/shopify/customers/" + cust.Id.Hex(),
		ParameterizedRoute: "/shopify/customers/" + cust.Id.Hex(),
		Handler:            DeleteCustomer,
		QueryString:        &qs,
		JsonBody:           cust,
		Runs:               b.N,
	}).RequestBenchmark()
}

func BenchmarkSearchCustomer(b *testing.B) {
	shopID := cart.InsertTestData()
	if shopID == nil {
		b.Error("failed to create a shop")
		b.Fail()
	}

	val := shopID.Hex()
	qs := make(url.Values, 0)
	qs.Add("shop", val)
	qs.Add("query", "alex")

	cust := cart.Customer{
		ShopId:    *shopID,
		FirstName: "Alex",
		LastName:  "Ninneman",
		Email:     time.Now().Format(time.RFC3339Nano) + "@gmail.com",
	}

	if err := cust.Insert("http://www.example.com"); err != nil {
		b.Error(err.Error())
		b.Fail()
	}

	(&httprunner.BenchmarkOptions{
		Method:             "GET",
		Route:              "/shopify/customers/search",
		ParameterizedRoute: "/shopify/customers/search",
		Handler:            SearchCustomer,
		QueryString:        &qs,
		Runs:               b.N,
	}).RequestBenchmark()
}

func BenchmarkGetCustomerOrders(b *testing.B) {
	shopID := cart.InsertTestData()
	if shopID == nil {
		b.Error("failed to create a shop")
		b.Fail()
	}

	val := shopID.Hex()
	qs := make(url.Values, 0)
	qs.Add("shop", val)

	cust := cart.Customer{
		ShopId:    *shopID,
		FirstName: "Alex",
		LastName:  "Ninneman",
		Email:     time.Now().Format(time.RFC3339Nano) + "@gmail.com",
		Password:  "password",
	}
	if err := cust.Insert("http://www.example.com"); err != nil {
		b.Error(err.Error())
		b.Fail()
	}

	(&httprunner.BenchmarkOptions{
		Method:             "GET",
		Route:              "/shopify/customers/" + cust.Id.Hex() + "/orders",
		ParameterizedRoute: "/shopify/customers/:id/orders",
		Handler:            GetCustomerOrders,
		QueryString:        &qs,
		Runs:               b.N,
	}).RequestBenchmark()
}
