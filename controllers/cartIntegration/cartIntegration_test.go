package cartIntegration

import (
	"bytes"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/cartIntegration"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/curt-labs/GoAPI/models/products"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestCartIntegration(t *testing.T) {
	//setup
	var err error
	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}

	var c cartIntegration.CartIntegration
	var cs []cartIntegration.CartIntegration
	var p products.Part
	var price customer.Price

	var cust customer.Customer
	cust.CustomerId = dtx.CustomerID
	cust.Create()
	t.Log(cust)

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
		thyme := time.Now()
		testThatHttp.Request("post", "/cart", "", "", SaveCI, bodyJson, "application/json")
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)

		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(err, ShouldBeNil)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		So(c, ShouldHaveSameTypeAs, cartIntegration.CartIntegration{})

		//test update CartIntegration
		c.PartID = p.ID
		bodyBytes, _ = json.Marshal(c)
		bodyJson = bytes.NewReader(bodyBytes)
		thyme = time.Now()
		testThatHttp.Request("post", "/cart/", ":id", strconv.Itoa(c.ID), SaveCI, bodyJson, "application/json")
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)

		So(err, ShouldBeNil)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		So(c, ShouldHaveSameTypeAs, cartIntegration.CartIntegration{})

		//test get CartIntegration
		thyme = time.Now()
		testThatHttp.Request("get", "/cart/", ":id", strconv.Itoa(c.ID), GetCI, nil, "")
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)

		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(err, ShouldBeNil)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		So(c.CustID, ShouldEqual, cust.Id)
		So(c, ShouldHaveSameTypeAs, cartIntegration.CartIntegration{})

		//test get CartIntegration by part
		thyme = time.Now()
		testThatHttp.Request("get", "/cart/part/", ":id", strconv.Itoa(c.PartID), GetCIbyPart, nil, "")
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cs)

		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		So(err, ShouldBeNil)
		So(cs, ShouldHaveSameTypeAs, []cartIntegration.CartIntegration{})
		So(len(cs), ShouldBeGreaterThan, 0)

		//test get CartIntegration by customer
		thyme = time.Now()
		testThatHttp.RequestWithDtx("get", "/cart/customer/", ":id", strconv.Itoa(cust.CustomerId), GetCIbyCustomer, nil, "", dtx)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cs)

		So(err, ShouldBeNil)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		So(cs, ShouldHaveSameTypeAs, []cartIntegration.CartIntegration{})
		log.Println("cs is: ")
		log.Println(cs)
		So(len(cs), ShouldBeGreaterThan, 0)

		// //test get CartIntegration by customer
		thyme = time.Now()
		testThatHttp.Request("get", "/cart/customer/count/", ":custID", strconv.Itoa(cust.CustomerId), GetCustomerPricingCount, nil, "")
		var count int
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &count)

		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		So(err, ShouldBeNil)
		So(count, ShouldBeGreaterThan, 0)

		//test get CustomerPricing
		thyme = time.Now()
		testThatHttp.Request("get", "/cart/customer/pricing/", ":custID/:page/:count", strconv.Itoa(cust.CustomerId)+"/1/1", GetCustomerPricingPaged, nil, "")
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cs)

		So(err, ShouldBeNil)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		So(cs, ShouldHaveSameTypeAs, []cartIntegration.CartIntegration{})
		So(len(cs), ShouldEqual, 1)

		//test get CustomerPricingPaged
		thyme = time.Now()
		testThatHttp.Request("get", "/cart/customer/pricing/", ":custID", strconv.Itoa(cust.CustomerId), GetCustomerPricing, nil, "")
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cs)

		So(err, ShouldBeNil)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds())
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		So(cs, ShouldHaveSameTypeAs, []cartIntegration.CartIntegration{})
		So(len(cs), ShouldBeGreaterThan, 0)

		//test delete CartIntegration
		thyme = time.Now()
		testThatHttp.Request("delete", "/cart/", ":id", strconv.Itoa(c.ID), DeleteCI, nil, "")
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)

		So(testThatHttp.Response.Code, ShouldEqual, 200)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(err, ShouldBeNil)

	})
	//cleanup
	cust.Delete()
	p.Delete()
	price.Delete()
}

func BenchmarkBrands(b *testing.B) {
	testThatHttp.RequestBenchmark(b.N, "GET", "/cart/1", nil, GetCI)
	testThatHttp.RequestBenchmark(b.N, "GET", "/cart/11000", nil, GetCIbyPart)
	testThatHttp.RequestBenchmark(b.N, "GET", "/cart/customer/2", nil, GetCIbyCustomer)
	testThatHttp.RequestBenchmark(b.N, "GET", "/cart//customer/count/2", nil, GetCustomerPricingCount)
	testThatHttp.RequestBenchmark(b.N, "GET", "/cart/customer/pricing/2/1/1", nil, GetCustomerPricingPaged)
	testThatHttp.RequestBenchmark(b.N, "GET", "/cart/customer/pricing/2", nil, GetCustomerPricing)
}
