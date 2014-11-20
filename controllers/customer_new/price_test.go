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

func TestCustomerPrice(t *testing.T) {
	var err error
	var p customer_new.Price
	var ps customer_new.Prices
	var cu customer_new.CustomerUser
	var c customer_new.Customer
	c.Name = "Dog Bountyhunter"
	c.Create()

	// //setup
	cu.Name = "test cust user"
	cu.Email = "pretend@test.com"
	cu.Password = "test"
	cu.Sudo = true
	cu.Create()
	var apiKey string
	for _, key := range cu.Keys {
		if strings.ToLower(key.Type) == "public" {
			apiKey = key.Key
		}
	}

	Convey("Testing Customer_New/Price", t, func() {
		//test create customer price
		form := url.Values{"custID": {strconv.Itoa(c.Id)}, "partID": {"11000"}, "price": {"123456"}}
		v := form.Encode()
		body := strings.NewReader(v)
		thyme := time.Now()
		testThatHttp.Request("post", "/new/customer/prices", "", "?key="+apiKey, CreateUpdatePrice, body, "application/x-www-form-urlencoded")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &p)
		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, customer_new.Price{})
		So(p.ID, ShouldBeGreaterThan, 0)

		//test update customer price
		form = url.Values{"isSale": {"true"}, "saleStart": {"01/01/2001"}, "saleEnd": {"01/01/2015"}}
		v = form.Encode()
		body = strings.NewReader(v)
		thyme = time.Now()
		testThatHttp.Request("post", "/new/customer/prices/", ":id", strconv.Itoa(p.ID)+"?key="+apiKey, CreateUpdatePrice, body, "application/x-www-form-urlencoded")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &p)
		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, customer_new.Price{})
		So(p.IsSale, ShouldEqual, 1)
		start, _ := time.Parse(inputTimeFormat, "01/01/2001")
		So(p.SaleStart, ShouldResemble, start)

		//test get customer price
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/prices/", ":id", strconv.Itoa(p.ID)+"?key="+apiKey, GetPrice, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &p)
		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, customer_new.Price{})

		//test get all customer price
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/prices", "", "?key="+apiKey, GetAllPrices, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &ps)
		So(err, ShouldBeNil)
		So(ps, ShouldHaveSameTypeAs, customer_new.Prices{})

		//test get customer price by part
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/prices/part/", ":id", strconv.Itoa(p.ID)+"?key="+apiKey, GetPricesByPart, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &ps)
		So(err, ShouldBeNil)
		So(ps, ShouldHaveSameTypeAs, customer_new.Prices{})

		//test get customer price by customer
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/pricesByCustomer/", ":id", strconv.Itoa(c.Id)+"?key="+apiKey, GetPriceByCustomer, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &p)
		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, customer_new.Price{})

		//test get sales
		form = url.Values{"id": {strconv.Itoa(c.Id)}, "start": {"01/01/2000"}, "end": {"01/01/2016"}}
		v = form.Encode()
		body = strings.NewReader(v)
		thyme = time.Now()
		testThatHttp.Request("post", "/new/customer/prices/sale", "", "?key="+apiKey, GetSales, body, "application/x-www-form-urlencoded")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &ps)
		So(err, ShouldBeNil)
		So(ps, ShouldHaveSameTypeAs, customer_new.Prices{})

		//test delete customer price
		thyme = time.Now()
		testThatHttp.Request("delete", "/new/customer/prices/", ":id", strconv.Itoa(p.ID)+"?key="+apiKey, DeletePrice, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &p)
		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, customer_new.Price{})
	})
	//teardown
	cu.Delete()
	c.Delete()
}
