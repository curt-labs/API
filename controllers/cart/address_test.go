package cart_ctlr

import (
	"encoding/json"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/helpers/httprunner"
	"github.com/curt-labs/API/models/cart"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"testing"
	"time"
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
			response = httprunner.ParameterizedRequest("GET", "/shopify/customers/:id/addresses", "/shopify/customers/1234/addresses", &qs, nil, GetAddresses)
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
			t.Log(response.Body.String())
			So(response.Code, ShouldEqual, 200)
			So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)

			response = httprunner.ParameterizedRequest("GET", "/shopify/customers/:id/addresses", "/shopify/customers/"+shopID.Hex()+"/addresses", &qs, nil, GetAddresses)
			So(response.Code, ShouldEqual, 500)
			So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

			response = httprunner.ParameterizedRequest("GET", "/shopify/customers/:id/addresses", "/shopify/customers/"+cust.Id.Hex()+"/addresses", &qs, nil, GetAddresses)
			So(response.Code, ShouldEqual, 200)
			So(json.Unmarshal(response.Body.Bytes(), &[]cart.CustomerAddress{}), ShouldBeNil)

			qs.Add("limit", "10")
			qs.Add("page", "1")
			response = httprunner.ParameterizedRequest("GET", "/shopify/customers/:id/addresses", "/shopify/customers/"+cust.Id.Hex()+"/addresses", &qs, nil, GetAddresses)
			So(response.Code, ShouldEqual, 200)
			So(json.Unmarshal(response.Body.Bytes(), &[]cart.CustomerAddress{}), ShouldBeNil)
		})
	})
}

func Test_AddAddress(t *testing.T) {
	Convey("with no shop identifier", t, func() {
		qs := make(url.Values, 0)

		response = httprunner.ParameterizedJsonRequest("POST", "/shopify/customers/:id/addresses", "/shopify/customers/1234/addresses", &qs, nil, AddAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
	})

	Convey("with shop identifier", t, func() {
		shopID := cart.InsertTestData()
		So(shopID, ShouldNotBeNil)
		qs := make(url.Values, 0)
		qs.Add("shop", shopID.Hex())

		response = httprunner.ParameterizedJsonRequest("POST", "/shopify/customers/:id/addresses", "/shopify/customers/1234/addresses", &qs, nil, AddAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		response = httprunner.ParameterizedJsonRequest("POST", "/shopify/customers/:id/addresses", "/shopify/customers/"+shopID.Hex()+"/addresses", &qs, nil, AddAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

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

		response = httprunner.ParameterizedJsonRequest("POST", "/shopify/customers/:id/addresses", "/shopify/customers/"+cust.Id.Hex()+"/addresses", &qs, &cart.Shop{}, AddAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		response = httprunner.ParameterizedJsonRequest("POST", "/shopify/customers/:id/addresses", "/shopify/customers/"+cust.Id.Hex()+"/addresses", &qs, &cart.Shop{}, AddAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		addr := cart.CustomerAddress{
			Address1:     "1119 Sunset Lane",
			City:         "Altoona",
			Company:      "AN & Co.",
			FirstName:    "Alex",
			LastName:     "Ninneman",
			Phone:        "7153082604",
			Province:     "Wisconsin",
			ProvinceCode: "WI",
			Country:      "US",
			CountryCode:  "US",
			CountryName:  "United States",
		}

		response = httprunner.ParameterizedJsonRequest("POST", "/shopify/customers/:id/addresses", "/shopify/customers/"+cust.Id.Hex()+"/addresses", &qs, &addr, AddAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		addr.Zip = "54720"
		response = httprunner.ParameterizedJsonRequest("POST", "/shopify/customers/:id/addresses", "/shopify/customers/"+cust.Id.Hex()+"/addresses", &qs, &addr, AddAddress)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &cart.CustomerAddress{}), ShouldBeNil)
	})
}

func Test_GetAddress(t *testing.T) {

	Convey("with shop identifier", t, func() {
		shopID := cart.InsertTestData()
		So(shopID, ShouldNotBeNil)
		qs := make(url.Values, 0)
		qs.Add("shop", shopID.Hex())

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

		addr := cart.CustomerAddress{
			Address1:     "1119 Sunset Lane",
			City:         "Altoona",
			Company:      "AN & Co.",
			FirstName:    "Alex",
			LastName:     "Ninneman",
			Phone:        "7153082604",
			Province:     "Wisconsin",
			ProvinceCode: "WI",
			Country:      "US",
			CountryCode:  "US",
			CountryName:  "United States",
			Zip:          "54720",
		}

		response = httprunner.ParameterizedJsonRequest("POST", "/shopify/customers/:id/addresses", "/shopify/customers/"+cust.Id.Hex()+"/addresses", &qs, &addr, AddAddress)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &addr), ShouldBeNil)

		response = httprunner.ParameterizedJsonRequest("GET", "/shopify/customers/:id/addresses/:address", "/shopify/customers/1234/addresses/1235", &qs, nil, GetAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		response = httprunner.ParameterizedJsonRequest("GET", "/shopify/customers/:id/addresses/:address", "/shopify/customers/"+cust.Id.Hex()+"/addresses/1235", &qs, nil, GetAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		response = httprunner.ParameterizedJsonRequest("GET", "/shopify/customers/:id/addresses/:address", "/shopify/customers/"+shopID.Hex()+"/addresses/"+addr.Id.Hex(), &qs, nil, GetAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		response = httprunner.ParameterizedJsonRequest("GET", "/shopify/customers/:id/addresses/:address", "/shopify/customers/"+cust.Id.Hex()+"/addresses/"+shopID.Hex(), &qs, nil, GetAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		response = httprunner.ParameterizedJsonRequest("GET", "/shopify/customers/:id/addresses/:address", "/shopify/customers/"+cust.Id.Hex()+"/addresses/"+addr.Id.Hex(), &qs, nil, GetAddress)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &addr), ShouldBeNil)
	})
}

func Test_EditAddress(t *testing.T) {

	Convey("with shop identifier", t, func() {
		shopID := cart.InsertTestData()
		So(shopID, ShouldNotBeNil)
		qs := make(url.Values, 0)
		qs.Add("shop", shopID.Hex())

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

		addr := cart.CustomerAddress{
			Address1:     "1119 Sunset Lane",
			City:         "Altoona",
			Company:      "AN & Co.",
			FirstName:    "Alex",
			LastName:     "Ninneman",
			Phone:        "7153082604",
			Province:     "Wisconsin",
			ProvinceCode: "WI",
			Country:      "US",
			CountryCode:  "US",
			CountryName:  "United States",
			Zip:          "54720",
		}

		response = httprunner.ParameterizedJsonRequest("POST", "/shopify/customers/:id/addresses", "/shopify/customers/"+cust.Id.Hex()+"/addresses", &qs, &addr, AddAddress)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &addr), ShouldBeNil)

		response = httprunner.ParameterizedJsonRequest("PUT", "/shopify/customers/:id/addresses/:address", "/shopify/customers/1234/addresses/1234", &qs, &addr, EditAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		response = httprunner.ParameterizedJsonRequest("PUT", "/shopify/customers/:id/addresses/:address", "/shopify/customers/"+shopID.Hex()+"/addresses/1234", &qs, &cust, EditAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		response = httprunner.ParameterizedJsonRequest("PUT", "/shopify/customers/:id/addresses/:address", "/shopify/customers/"+cust.Id.Hex()+"/addresses/1234", &qs, &cust, EditAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		response = httprunner.ParameterizedJsonRequest("PUT", "/shopify/customers/:id/addresses/:address", "/shopify/customers/"+cust.Id.Hex()+"/addresses/"+addr.Id.Hex(), &qs, &cart.Shop{}, EditAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		response = httprunner.ParameterizedJsonRequest("PUT", "/shopify/customers/:id/addresses/:address", "/shopify/customers/"+cust.Id.Hex()+"/addresses/"+addr.Id.Hex(), &qs, &cart.CustomerAddress{}, EditAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		response = httprunner.ParameterizedJsonRequest("PUT", "/shopify/customers/:id/addresses/:address", "/shopify/customers/"+cust.Id.Hex()+"/addresses/"+cust.Id.Hex(), &qs, &addr, EditAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		addr.Name = "Test Address"
		response = httprunner.ParameterizedJsonRequest("PUT", "/shopify/customers/:id/addresses/:address", "/shopify/customers/"+cust.Id.Hex()+"/addresses/"+addr.Id.Hex(), &qs, &addr, EditAddress)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &addr), ShouldBeNil)
	})
}

func Test_SetDefaultAddress(t *testing.T) {

	Convey("with shop identifier", t, func() {
		shopID := cart.InsertTestData()
		So(shopID, ShouldNotBeNil)
		qs := make(url.Values, 0)
		qs.Add("shop", shopID.Hex())

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

		addr := cart.CustomerAddress{
			Address1:     "1119 Sunset Lane",
			City:         "Altoona",
			Company:      "AN & Co.",
			FirstName:    "Alex",
			LastName:     "Ninneman",
			Phone:        "7153082604",
			Province:     "Wisconsin",
			ProvinceCode: "WI",
			Country:      "US",
			CountryCode:  "US",
			CountryName:  "United States",
			Zip:          "54720",
		}

		response = httprunner.ParameterizedJsonRequest("POST", "/shopify/customers/:id/addresses", "/shopify/customers/"+cust.Id.Hex()+"/addresses", &qs, &addr, AddAddress)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &addr), ShouldBeNil)

		response = httprunner.ParameterizedJsonRequest("PUT", "/shopify/customers/:id/addresses/:address/default", "/shopify/customers/1234/addresses/1234/default", &qs, &addr, SetDefaultAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		response = httprunner.ParameterizedJsonRequest("PUT", "/shopify/customers/:id/addresses/:address/default", "/shopify/customers/"+shopID.Hex()+"/addresses/1234/default", &qs, &cust, SetDefaultAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		response = httprunner.ParameterizedJsonRequest("PUT", "/shopify/customers/:id/addresses/:address/default", "/shopify/customers/"+cust.Id.Hex()+"/addresses/1234/default", &qs, &cust, SetDefaultAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		response = httprunner.ParameterizedJsonRequest("PUT", "/shopify/customers/:id/addresses/:address/default", "/shopify/customers/"+cust.Id.Hex()+"/addresses/"+cust.Id.Hex()+"/default", &qs, &cust, SetDefaultAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		response = httprunner.ParameterizedJsonRequest("PUT", "/shopify/customers/:id/addresses/:address/default", "/shopify/customers/"+cust.Id.Hex()+"/addresses/"+addr.Id.Hex()+"/default", &qs, &addr, SetDefaultAddress)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &addr), ShouldBeNil)
	})
}

func Test_DeleteAddress(t *testing.T) {

	Convey("with shop identifier", t, func() {
		shopID := cart.InsertTestData()
		So(shopID, ShouldNotBeNil)
		qs := make(url.Values, 0)
		qs.Add("shop", shopID.Hex())

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

		addr := cart.CustomerAddress{
			Address1:     "1119 Sunset Lane",
			City:         "Altoona",
			Company:      "AN & Co.",
			FirstName:    "Alex",
			LastName:     "Ninneman",
			Phone:        "7153082604",
			Province:     "Wisconsin",
			ProvinceCode: "WI",
			Country:      "US",
			CountryCode:  "US",
			CountryName:  "United States",
			Zip:          "54720",
		}

		response = httprunner.ParameterizedJsonRequest("POST", "/shopify/customers/:id/addresses", "/shopify/customers/"+cust.Id.Hex()+"/addresses", &qs, &addr, AddAddress)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &addr), ShouldBeNil)

		response = httprunner.ParameterizedRequest("DELETE", "/shopify/customers/:id/addresses/:address", "/shopify/customers/1234/addresses/1234", &qs, nil, DeleteAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		response = httprunner.ParameterizedRequest("DELETE", "/shopify/customers/:id/addresses/:address", "/shopify/customers/"+shopID.Hex()+"/addresses/1234", &qs, nil, DeleteAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		response = httprunner.ParameterizedRequest("DELETE", "/shopify/customers/:id/addresses/:address", "/shopify/customers/"+cust.Id.Hex()+"/addresses/1234", &qs, nil, DeleteAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		response = httprunner.ParameterizedRequest("DELETE", "/shopify/customers/:id/addresses/:address", "/shopify/customers/"+cust.Id.Hex()+"/addresses/"+cust.Id.Hex(), &qs, nil, DeleteAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		response = httprunner.ParameterizedJsonRequest("PUT", "/shopify/customers/:id/addresses/:address/default", "/shopify/customers/"+cust.Id.Hex()+"/addresses/"+addr.Id.Hex()+"/default", &qs, &addr, SetDefaultAddress)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)

		response = httprunner.ParameterizedRequest("DELETE", "/shopify/customers/:id/addresses/:address", "/shopify/customers/"+cust.Id.Hex()+"/addresses/"+addr.Id.Hex(), &qs, nil, DeleteAddress)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		response = httprunner.ParameterizedJsonRequest("POST", "/shopify/customers/:id/addresses", "/shopify/customers/"+cust.Id.Hex()+"/addresses", &qs, &addr, AddAddress)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &addr), ShouldBeNil)

		response = httprunner.ParameterizedRequest("DELETE", "/shopify/customers/:id/addresses/:address", "/shopify/customers/"+cust.Id.Hex()+"/addresses/"+addr.Id.Hex(), &qs, nil, DeleteAddress)
		So(response.Code, ShouldEqual, 200)
		So(response.Body.String(), ShouldEqual, "")
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
		Email:     time.Now().Format(time.RFC3339Nano) + "@gmail.com",
		Password:  "password",
	}
	if err := cust.Insert("http://www.example.com"); err != nil {
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
