package cart_ctlr

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/models/cart"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"testing"
	"time"
)

func Test_AddAccount(t *testing.T) {
	Convey("with shop identifier", t, func() {
		shopID := cart.InsertTestData()
		So(shopID, ShouldNotBeNil)
		val := shopID.Hex()
		qs := make(url.Values, 0)
		qs.Add("shop", val)

		response = httprunner.Req(AddAccount, "POST", "", "/shopify/account", &qs, nil, cart.Shop{})
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		cust := cart.Customer{
			ShopId:    *shopID,
			FirstName: "Alex",
			LastName:  "Ninneman",
		}

		addr := cart.CustomerAddress{}
		addr.Address1 = "Test"

		response = httprunner.Req(AddAccount, "POST", "", "/shopify/account", &qs, nil, addr)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		response = httprunner.Req(AddAccount, "POST", "", "/shopify/account", &qs, nil, cust)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		cust.Email = time.Now().Format(time.RFC3339Nano) + "@gmail.com"

		response = httprunner.JsonRequest("POST", "/shopify/account", &qs, nil, AddAccount)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		cust.Password = "password"
		response = httprunner.JsonRequest("POST", "/shopify/account", &qs, cust, AddAccount)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)
	})
}

func Test_AccountLogin(t *testing.T) {

	Convey("with shop identifier", t, func() {
		shopID := cart.InsertTestData()
		So(shopID, ShouldNotBeNil)
		val := shopID.Hex()
		qs := make(url.Values, 0)
		qs.Add("shop", val)

		cust := cart.Customer{
			ShopId:    *shopID,
			FirstName: "Alex",
			LastName:  "Ninneman",
			Email:     time.Now().Format(time.RFC3339Nano) + "@gmail.com",
		}

		cust.Password = "password"
		response = httprunner.JsonRequest("POST", "/shopify/account", &qs, cust, AddAccount)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)

		cust.Password = ""
		response = httprunner.Req(AccountLogin, "POST", "", "/shopify/account/login", &qs, nil, cust)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		cust.Password = "password"
		response = httprunner.Req(AccountLogin, "POST", "", "/shopify/account/login", &qs, nil, cust)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)
	})
}

func Test_GetAccount(t *testing.T) {

	Convey("with shop identifier", t, func() {
		shopID := cart.InsertTestData()
		So(shopID, ShouldNotBeNil)
		val := shopID.Hex()
		qs := make(url.Values, 0)
		qs.Add("shop", val)

		cust := cart.Customer{
			ShopId:    *shopID,
			FirstName: "Alex",
			LastName:  "Ninneman",
			Email:     time.Now().Format(time.RFC3339Nano) + "@gmail.com",
		}

		cust.Password = "password"
		response = httprunner.JsonRequest("POST", "/shopify/account", &qs, cust, AddAccount)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)

		response = httprunner.Req(GetAccount, "GET", "", "/shopify/account", &qs, nil, cust)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		header := map[string]interface{}{
			"Authorization": "Bearer ",
		}
		response = httprunner.Req(GetAccount, "GET", "", "/shopify/account", &qs, nil, nil, header)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		header = map[string]interface{}{
			"Authorization": "Bearer " + cust.Token,
		}
		response = httprunner.Req(GetAccount, "GET", "", "/shopify/account", &qs, nil, nil, header)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)
	})
}

func Test_EditAccount(t *testing.T) {

	Convey("with shop identifier", t, func() {
		shopID := cart.InsertTestData()
		So(shopID, ShouldNotBeNil)
		val := shopID.Hex()
		qs := make(url.Values, 0)
		qs.Add("shop", val)

		cust := cart.Customer{
			ShopId:    *shopID,
			FirstName: "Alex",
			LastName:  "Ninneman",
			Email:     time.Now().Format(time.RFC3339Nano) + "@gmail.com",
		}

		cust.Password = "password"
		response = httprunner.JsonRequest("POST", "/shopify/account", &qs, cust, AddAccount)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)

		cust.Email = time.Now().Format(time.RFC3339Nano) + "@gmail.com"
		header := map[string]interface{}{
			"Authorization": "Bearer as;ldskfja;lfdj",
		}
		response = httprunner.Req(EditAccount, "PUT", "", "/shopify/account", &qs, nil, cust, header)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		header = map[string]interface{}{
			"Authorization": "Bearer " + cust.Token,
		}
		cust.FirstName = ""
		response = httprunner.Req(EditAccount, "PUT", "", "/shopify/account", &qs, nil, cust, header)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)

		cust.FirstName = "Alex"
		response = httprunner.Req(EditAccount, "PUT", "", "/shopify/account", &qs, nil, cust, header)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &cust), ShouldBeNil)
	})
}
