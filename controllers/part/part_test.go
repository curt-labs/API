package part_ctlr

import (
	"bytes"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/apiKeyType"
	"github.com/curt-labs/GoAPI/models/customer_new"
	"github.com/curt-labs/GoAPI/models/products"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestParts(t *testing.T) {
	var err error
	var p products.Part
	var price products.Price
	var cu customer_new.CustomerUser
	var cat products.Category
	cat.Create()

	//setup apiKeyTypes
	var pub, pri, auth apiKeyType.ApiKeyType
	pub.Type = "public"
	pri.Type = "private"
	auth.Type = "authentication"
	pub.Create()
	pri.Create()
	auth.Create()

	var c customer_new.Customer
	c.Name = "test man"
	c.Create()

	cu.CustomerID = c.Id
	cu.Name = "test cust user"
	cu.Email = "pretend@test.com"
	cu.Password = "test"
	cu.Sudo = true
	err = cu.Create()
	t.Log(err)
	var apiKey string
	for _, key := range cu.Keys {
		if strings.ToLower(key.Type) == "public" {
			apiKey = key.Key
		}
	}
	t.Log("APIKEY", apiKey)

	Convey("TestingParts", t, func() {

		//test create part
		p.ID = 10999
		p.Categories = append(p.Categories, cat)
		p.OldPartNumber = "8675309"
		p.ShortDesc = "test part"
		bodyBytes, _ := json.Marshal(p)
		bodyJson := bytes.NewReader(bodyBytes)
		thyme := time.Now()
		testThatHttp.Request("post", "/part", "", "?key="+apiKey, CreatePart, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &p)
		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, products.Part{})
		So(p.ID, ShouldEqual, 10999)

		err = p.BindCustomer(apiKey) //setup
		So(err, ShouldBeNil)

		//test create price
		price.Price = 987
		price.PartId = p.ID
		price.Type = "test"
		bodyBytes, _ = json.Marshal(price)
		bodyJson = bytes.NewReader(bodyBytes)
		thyme = time.Now()
		testThatHttp.Request("post", "/price", "", "?key="+apiKey, SavePrice, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &price)
		So(err, ShouldBeNil)
		So(price, ShouldHaveSameTypeAs, products.Price{})
		So(price.Id, ShouldBeGreaterThan, 0)

		//test update price
		price.Type = "tester"
		bodyBytes, _ = json.Marshal(price)
		bodyJson = bytes.NewReader(bodyBytes)
		thyme = time.Now()
		testThatHttp.Request("post", "/price/", ":id", strconv.Itoa(price.Id)+"?key="+apiKey, SavePrice, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &price)
		So(err, ShouldBeNil)
		So(price, ShouldHaveSameTypeAs, products.Price{})
		So(price.Type, ShouldNotEqual, "test")

		//test get part prices
		thyme = time.Now()
		testThatHttp.Request("get", "/part/", ":id/prices", strconv.Itoa(p.ID)+"/prices?key="+apiKey, Prices, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var prices []products.Price
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &prices)
		So(err, ShouldBeNil)
		So(prices, ShouldHaveSameTypeAs, []products.Price{})

		//test get part categories
		thyme = time.Now()
		testThatHttp.Request("get", "/part/", ":part/categories", strconv.Itoa(p.ID)+"/categories?key="+apiKey, Categories, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var cats []products.Category
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cats)
		So(err, ShouldBeNil)
		So(cats, ShouldHaveSameTypeAs, []products.Category{})

		//test get price
		thyme = time.Now()
		testThatHttp.Request("get", "/price/", ":id", strconv.Itoa(price.Id)+"?key="+apiKey, GetPrice, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &price)
		So(err, ShouldBeNil)
		So(price, ShouldHaveSameTypeAs, products.Price{})

		//test get old part Number
		thyme = time.Now()
		testThatHttp.Request("get", "/part/old/", ":part", p.OldPartNumber+"?key="+apiKey, OldPartNumber, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &p)
		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, products.Part{})
		So(p.OldPartNumber, ShouldEqual, "8675309")
		So(p.ID, ShouldEqual, 10999)

		//test delete price
		thyme = time.Now()
		testThatHttp.Request("delete", "/price/", ":id", strconv.Itoa(price.Id)+"?key="+apiKey, DeletePrice, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &price)
		So(err, ShouldBeNil)
		So(price, ShouldHaveSameTypeAs, products.Price{})

		//test update part
		p.OldPartNumber = "8675309"
		p.InstallSheet, err = url.Parse("www.sheetsrus.com")
		bodyBytes, _ = json.Marshal(p)
		bodyJson = bytes.NewReader(bodyBytes)
		thyme = time.Now()
		testThatHttp.Request("put", "/part/", ":id", strconv.Itoa(p.ID)+"?key="+apiKey, UpdatePart, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &p)
		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, products.Part{})
		So(p.OldPartNumber, ShouldEqual, "8675309")
		So(p.ID, ShouldEqual, 10999)

		//test delete part
		thyme = time.Now()
		testThatHttp.Request("delete", "/part/", ":id", strconv.Itoa(p.ID)+"?key="+apiKey, DeletePart, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &p)
		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, products.Part{})

	})
	// cu.Delete()
	p.Delete()
	cat.Delete()
	pub.Delete()
	pri.Delete()
	auth.Delete()
}
