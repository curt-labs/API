package cartIntegration

import (
	"bytes"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/cartIntegration"
	"github.com/curt-labs/GoAPI/models/customer_new"
	"github.com/curt-labs/GoAPI/models/products"
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"testing"
)

func TestCartIntegration(t *testing.T) {
	//setup
	var c cartIntegration.CartIntegration
	var cs []cartIntegration.CartIntegration
	var p products.Part
	var price customer_new.Price
	var err error
	var cust customer_new.Customer
	cust.CustomerId = 666
	cust.Create()

	p.ShortDesc = "test"
	p.ID = 123456789
	p.Status = 800
	err = p.Create()
	if err != nil {
		err = nil
		err = p.Update()
	}
	price.CustID = cust.Id
	price.PartID = p.ID
	price.Create()

	Convey("Testing CartIntegration", t, func() {
		//test create CartIntegration
		c.PartID = p.ID
		c.CustID = cust.Id
		bodyBytes, _ := json.Marshal(c)
		bodyJson := bytes.NewReader(bodyBytes)
		testThatHttp.Request("post", "/cart", "", "", SaveCI, bodyJson, "application/json")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, cartIntegration.CartIntegration{})

		//test update CartIntegration
		c.PartID = p.ID
		bodyBytes, _ = json.Marshal(c)
		bodyJson = bytes.NewReader(bodyBytes)
		testThatHttp.Request("post", "/cart/", ":id", strconv.Itoa(c.ID), SaveCI, bodyJson, "application/json")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, cartIntegration.CartIntegration{})

		//test get CartIntegration
		testThatHttp.Request("get", "/cart/", ":id", strconv.Itoa(c.ID), GetCI, nil, "")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c.CustID, ShouldEqual, cust.Id)
		So(c, ShouldHaveSameTypeAs, cartIntegration.CartIntegration{})

		//test get CartIntegration by part
		testThatHttp.Request("get", "/cart/part/", ":id", strconv.Itoa(c.PartID), GetCIbyPart, nil, "")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cs)
		So(err, ShouldBeNil)
		So(cs, ShouldHaveSameTypeAs, []cartIntegration.CartIntegration{})
		So(len(cs), ShouldBeGreaterThan, 0)

		//test get CartIntegration by customer
		testThatHttp.Request("get", "/cart/customer/", ":id", strconv.Itoa(cust.CustomerId), GetCIbyCustomer, nil, "")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cs)
		So(err, ShouldBeNil)
		So(cs, ShouldHaveSameTypeAs, []cartIntegration.CartIntegration{})
		So(len(cs), ShouldBeGreaterThan, 0)

		// //test get CartIntegration by customer
		testThatHttp.Request("get", "/cart/customer/count/", ":custID", strconv.Itoa(cust.CustomerId), GetCustomerPricingCount, nil, "")
		var count int
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &count)
		So(err, ShouldBeNil)
		So(count, ShouldBeGreaterThan, 0)

		//test get CustomerPricing
		testThatHttp.Request("get", "/cart/customer/pricing/", ":custID/:page/:count", strconv.Itoa(cust.CustomerId)+"/1/1", GetCustomerPricingPaged, nil, "")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cs)
		So(err, ShouldBeNil)
		So(cs, ShouldHaveSameTypeAs, []cartIntegration.CartIntegration{})
		So(len(cs), ShouldEqual, 1)

		//test get CustomerPricingPaged
		testThatHttp.Request("get", "/cart/customer/pricing/", ":custID", strconv.Itoa(cust.CustomerId), GetCustomerPricing, nil, "")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cs)
		So(err, ShouldBeNil)
		So(cs, ShouldHaveSameTypeAs, []cartIntegration.CartIntegration{})
		So(len(cs), ShouldBeGreaterThan, 0)

		//test delete CartIntegration
		testThatHttp.Request("delete", "/cart/", ":id", strconv.Itoa(c.ID), DeleteCI, nil, "")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)

	})
	//cleanup
	cust.Delete()
	p.Delete()
	price.Delete()
}
