package cart_ctlr

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/models/cart"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"testing"
)

func Test_GetCustomerAddresses(t *testing.T) {
	Convey("no shop identifier", t, func() {
		response = httprunner.Request("GET", "/shopify/customers/1234/addresses", nil, GetAddresses)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		vals := make(url.Values, 0)
		vals.Add("shop", "testing")
		response = httprunner.Request("GET", "/shopify/customers/1234/addresses", nil, GetAddresses)
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
			response = httprunner.ParameterizedRequest("GET", "/shopify/customers/:id/addresses", "/shopify/customers/1234/addresses", &qs, GetAddresses)
			So(response.Code, ShouldEqual, 500)
			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
		})

		Convey("with good customer reference", func() {
			cust := cart.Customer{
				ShopId:    *shopID,
				FirstName: "Alex",
				LastName:  "Ninneman",
				Email:     "ninnemana@gmail.com",
			}
			response = httprunner.JsonRequest("POST", "/shopify/customers", &qs, cust, AddCustomer)
			So(response.Code, ShouldEqual, 200)
			So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)

			response = httprunner.ParameterizedRequest("GET", "/shopify/customers/:id/addresses", "/shopify/customers/"+shopID.Hex()+"/addresses", &qs, GetAddresses)
			So(response.Code, ShouldEqual, 500)
			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

			response = httprunner.ParameterizedRequest("GET", "/shopify/customers/:id/addresses", "/shopify/customers/"+cust.Id.Hex()+"/addresses", &qs, GetAddresses)
			So(response.Code, ShouldEqual, 200)
			So(json.Unmarshal(response.Body.Bytes(), &[]cart.CustomerAddress{}), ShouldBeNil)

			qs.Add("limit", "10")
			qs.Add("page", "1")
			response = httprunner.ParameterizedRequest("GET", "/shopify/customers/:id/addresses", "/shopify/customers/"+cust.Id.Hex()+"/addresses", &qs, GetAddresses)
			So(response.Code, ShouldEqual, 200)
			So(json.Unmarshal(response.Body.Bytes(), &[]cart.CustomerAddress{}), ShouldBeNil)
		})
	})
}

func BenchmarkGetCustomerAddresses(b *testing.B) {
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
	}
	if err := cust.Insert(); err != nil {
		b.Error(err.Error())
		b.Fail()
	}

	(&httprunner.BenchmarkOptions{
		Method:             "GET",
		Route:              "/shopify/customers/" + cust.Id.Hex() + "/addresses",
		ParameterizedRoute: "/shopify/customers/:id/addresses",
		Handler:            GetAddresses,
		QueryString:        &qs,
		Runs:               b.N,
	}).RequestBenchmark()
}
