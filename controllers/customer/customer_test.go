package customer_ctlr

import (
	"bytes"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/curt-labs/GoAPI/models/customer/content"
	. "github.com/smartystreets/goconvey/convey"

	"flag"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestCustomer(t *testing.T) {
	var err error
	var c customer.Customer
	var cu customer.CustomerUser
	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}

	cu.Id = dtx.UserID
	err = cu.Get(dtx.APIKey)
	if err != nil {
		t.Log(err)
	}
	c.Users = append(c.Users, cu)

	Convey("Testing customer/Customer", t, func() {
		//test create customer
		c.Name = "Jason Voorhees"
		c.Email = "jason@crystal.lake"
		bodyBytes, _ := json.Marshal(c)
		bodyJson := bytes.NewReader(bodyBytes)
		thyme := time.Now()
		testThatHttp.Request("put", "/customer", "", "?key=", SaveCustomer, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, customer.Customer{})
		So(c.Id, ShouldBeGreaterThan, 0)

		//test update customer
		c.Fax = "666-1313"
		c.State.Id = 1
		bodyBytes, _ = json.Marshal(c)
		bodyJson = bytes.NewReader(bodyBytes)
		thyme = time.Now()
		testThatHttp.Request("put", "/customer/", ":id", strconv.Itoa(c.Id)+"?key=", SaveCustomer, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, customer.Customer{})
		So(c.Id, ShouldBeGreaterThan, 0)

		thyme = time.Now()
		testThatHttp.Request("get", "/customer/", ":id", strconv.Itoa(c.Id)+"?key=", GetCustomer, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, customer.Customer{})
		So(c.Id, ShouldBeGreaterThan, 0)

		// get customer locations
		thyme = time.Now()
		testThatHttp.RequestWithDtx("get", "/customer/locations", "", "?key=", GetLocations, nil, "", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()*6)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c.Locations)
		So(err, ShouldBeNil)
		So(c.Locations, ShouldHaveSameTypeAs, []customer.CustomerLocation{})

		// //get user
		thyme = time.Now()
		testThatHttp.RequestWithDtx("post", "/customer/user", "", "?key=", GetUser, nil, "", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cu)
		So(err, ShouldBeNil)
		So(cu, ShouldHaveSameTypeAs, customer.CustomerUser{})

		//get users
		thyme = time.Now()
		testThatHttp.RequestWithDtx("get", "/customer/users", "", "?key=", GetUsers, nil, "", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()*5)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var cus []customer.CustomerUser
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cus)
		So(err, ShouldBeNil)
		So(cus, ShouldHaveSameTypeAs, []customer.CustomerUser{})

		// get customer price
		// price.CustID = c.Id
		// price.Create()
		// thyme = time.Now()
		// testThatHttp.Request("get", "/new/customer/price/", ":id", strconv.Itoa(p.ID)+"?key="+apiKey, GetCustomerPrice, nil, "")
		// So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		// So(testThatHttp.Response.Code, ShouldEqual, 200)
		// var price float64
		// err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &price)
		// So(err, ShouldBeNil)
		// So(price, ShouldHaveSameTypeAs, 7.1)

		// //get customer cart reference
		// ci.CustID = c.Id
		// ci.Create()
		// thyme = time.Now()
		// testThatHttp.Request("get", "/new/customer/cartRef/", ":id", strconv.Itoa(p.ID)+"?key="+apiKey, GetCustomerCartReference, nil, "")
		// So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		// So(testThatHttp.Response.Code, ShouldEqual, 200)
		// var reference int
		// err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &reference)
		// So(err, ShouldBeNil)
		// So(reference, ShouldHaveSameTypeAs, 7)

		//test delete customer
		thyme = time.Now()
		testThatHttp.Request("delete", "/customer/", ":id", strconv.Itoa(c.Id)+"?key=", DeleteCustomer, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, customer.Customer{})
		So(c.Id, ShouldBeGreaterThan, 0)

	})
	//cleanup
	err = cu.Delete()
	err = c.Delete()
	err = apicontextmock.DeMock(dtx)
	if err != nil {
		t.Log(err)
	}
}

func BenchmarkCRUDCustomer(b *testing.B) {
	dtx, err := apicontextmock.Mock()
	if err != nil {
		b.Log(err)
	}

	qs := make(url.Values, 0)
	qs.Add("key", dtx.APIKey)

	Convey("Customer", b, func() {
		var c customer.Customer
		c.Name = "Freddy Krueger"
		c.Email = "freddy@elm.st"
		// var locs customer.CustomerLocations

		//create
		(&httprunner.BenchmarkOptions{
			Method:             "POST",
			Route:              "/customer",
			ParameterizedRoute: "/customer",
			Handler:            SaveCustomer,
			QueryString:        &qs,
			JsonBody:           c,
			Runs:               b.N,
		}).RequestBenchmark()

		//get
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/customer",
			ParameterizedRoute: "/customer/" + strconv.Itoa(c.Id),
			Handler:            GetCustomer,
			QueryString:        &qs,
			JsonBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()

		//get locations
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/customer",
			ParameterizedRoute: "/customer",
			Handler:            GetLocations,
			QueryString:        &qs,
			JsonBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()

		//get user
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/customer/user",
			ParameterizedRoute: "/customer/user",
			Handler:            GetUser,
			QueryString:        &qs,
			JsonBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()

		//get users
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/customer/users",
			ParameterizedRoute: "/customer/users",
			Handler:            GetUsers,
			QueryString:        &qs,
			JsonBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()

		//delete
		(&httprunner.BenchmarkOptions{
			Method:             "DELETE",
			Route:              "/customer",
			ParameterizedRoute: "/customer/" + strconv.Itoa(c.Id),
			Handler:            DeleteCustomer,
			QueryString:        &qs,
			JsonBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()

	})
	//teardown customer user

	err = apicontextmock.DeMock(dtx)
	if err != nil {
		b.Log(err)
	}
}

func TestCustomerContent(t *testing.T) {
	flag.Parse()
	//customer - for db setup only
	var c customer.Customer
	var content custcontent.CustomerContent
	var partContent custcontent.PartContent
	var categoryContent custcontent.CustomerContent
	var ct custcontent.ContentType
	var crs custcontent.CustomerContentRevisions
	var contents []custcontent.CustomerContent
	var catContent custcontent.CategoryContent
	var catContents []custcontent.CategoryContent
	var err error
	var apiKey string

	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}
	catContent.CategoryId = 1

	ct.Type = "test"
	ct.Create()

	c.Name = "test cust"
	c.Create()

	Convey("Testing customer/Customer_content", t, func() {
		//test create part content
		content.Text = "new content"
		content.ContentType.Id = ct.Id
		bodyBytes, _ := json.Marshal(content)
		bodyJson := bytes.NewReader(bodyBytes)
		thyme := time.Now()
		testThatHttp.Request("post", "/customer/cms/part/", ":id", strconv.Itoa(11000)+"?key="+dtx.APIKey, CreatePartContent, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		if testThatHttp.Response.Code == 200 { //returns 500 when ninnemana user doesn't exist
			So(testThatHttp.Response.Code, ShouldEqual, 200)
			err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &content)
			So(err, ShouldBeNil)
			So(content, ShouldHaveSameTypeAs, custcontent.CustomerContent{})
		}

		//create category content
		categoryContent.Text = "new content"
		categoryContent.ContentType.Id = ct.Id
		bodyBytes, _ = json.Marshal(categoryContent)
		bodyJson = bytes.NewReader(bodyBytes)
		thyme = time.Now()
		testThatHttp.Request("post", "/customer/cms/category/", ":id", strconv.Itoa(catContent.CategoryId)+"?key="+apiKey, CreateCategoryContent, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &categoryContent)
		So(err, ShouldBeNil)
		So(categoryContent, ShouldHaveSameTypeAs, custcontent.CustomerContent{})

		//test update part content
		content.Text = "newerer content"
		bodyBytes, _ = json.Marshal(content)
		bodyJson = bytes.NewReader(bodyBytes)
		thyme = time.Now()
		testThatHttp.Request("put", "/customer/cms/part/", ":id", strconv.Itoa(11000)+"?key="+apiKey, UpdatePartContent, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &content)
		So(err, ShouldBeNil)
		So(content, ShouldHaveSameTypeAs, custcontent.CustomerContent{})

		//test update category content
		categoryContent.Text = "newerer content"
		bodyBytes, _ = json.Marshal(categoryContent)
		bodyJson = bytes.NewReader(bodyBytes)
		thyme = time.Now()
		testThatHttp.Request("put", "/customer/cms/part/", ":id", strconv.Itoa(catContent.CategoryId)+"?key="+apiKey, UpdateCategoryContent, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &categoryContent)
		So(err, ShouldBeNil)
		So(categoryContent, ShouldHaveSameTypeAs, custcontent.CustomerContent{})

		//test get part content (unique)
		thyme = time.Now()
		testThatHttp.Request("get", "/customer/cms/part/", ":id", strconv.Itoa(11000)+"?key="+apiKey, UniquePartContent, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &contents)
		So(err, ShouldBeNil)
		So(contents, ShouldHaveSameTypeAs, []custcontent.CustomerContent{})

		//test get all part content
		thyme = time.Now()
		testThatHttp.Request("get", "/customer/cms/part", "", "?key="+apiKey, AllPartContent, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &contents)
		So(err, ShouldBeNil)
		So(contents, ShouldHaveSameTypeAs, []custcontent.CustomerContent{})

		//test get category content (all content)
		thyme = time.Now()
		testThatHttp.Request("get", "/customer/cms/part", "", "?key="+apiKey, AllCategoryContent, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &contents)
		So(err, ShouldBeNil)
		So(contents, ShouldHaveSameTypeAs, []custcontent.CustomerContent{})

		//test get unique category content
		catContent.Content = append(catContent.Content, content) //setup some category Content
		thyme = time.Now()
		testThatHttp.Request("get", "/customer/cms/category/", ":id", strconv.Itoa(catContent.CategoryId)+"?key="+apiKey, UniqueCategoryContent, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &catContents)
		So(err, ShouldBeNil)
		So(catContents, ShouldHaveSameTypeAs, []custcontent.CategoryContent{})

		//test get all content
		thyme = time.Now()
		testThatHttp.Request("get", "/customer/cms", "", "?key="+apiKey, GetAllContent, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &contents)
		So(err, ShouldBeNil)
		So(contents, ShouldHaveSameTypeAs, []custcontent.CustomerContent{})

		//test get content by id
		thyme = time.Now()
		testThatHttp.Request("get", "/customer/cms/", ":id", strconv.Itoa(content.Id)+"?key="+apiKey, GetContentById, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &content)
		So(err, ShouldBeNil)
		So(content, ShouldHaveSameTypeAs, custcontent.CustomerContent{})

		//test get content revisions by id
		thyme = time.Now()
		testThatHttp.Request("get", "/customer/cms/", ":id/revisions", strconv.Itoa(content.Id)+"/revisions?key="+apiKey, GetContentRevisionsById, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &crs)
		So(err, ShouldBeNil)
		So(crs, ShouldHaveSameTypeAs, custcontent.CustomerContentRevisions{})

		//test get all content types
		thyme = time.Now()
		testThatHttp.Request("get", "/customer/cms/content_types", "", "?key="+apiKey, GetAllContentTypes, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var cts []custcontent.ContentType
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cts)
		So(err, ShouldBeNil)
		So(cts, ShouldHaveSameTypeAs, []custcontent.ContentType{})
		So(len(cts), ShouldBeGreaterThan, 0)

		//test delete part content
		bodyBytes, _ = json.Marshal(content)
		bodyJson = bytes.NewReader(bodyBytes)
		thyme = time.Now()
		testThatHttp.Request("delete", "/customer/cms/part/", ":id", strconv.Itoa(11000)+"?key="+apiKey, DeletePartContent, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &partContent)
		So(err, ShouldBeNil)
		So(partContent, ShouldHaveSameTypeAs, custcontent.PartContent{})

		//test delete category content
		bodyBytes, _ = json.Marshal(categoryContent)
		bodyJson = bytes.NewReader(bodyBytes)
		thyme = time.Now()
		testThatHttp.Request("delete", "/customer/cms/category/", ":id", strconv.Itoa(catContent.CategoryId)+"?key="+apiKey, DeleteCategoryContent, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &content)
		So(err, ShouldBeNil)
		So(content, ShouldHaveSameTypeAs, custcontent.CustomerContent{})

	})
	//teardown
	err = content.DeleteById()
	err = categoryContent.DeleteById()

	for _, con := range catContent.Content {
		err = con.DeleteById()
	}
	err = c.Delete()

	err = ct.Delete()

	err = apicontextmock.DeMock(dtx)
	if err != nil {
		t.Log(err)
	}
}

func BenchmarkCRUDContent(b *testing.B) {

	dtx, err := apicontextmock.Mock()
	if err != nil {
		b.Log(err)
	}
	qs := make(url.Values, 0)
	qs.Add("key", dtx.APIKey)

	Convey("Part Content", b, func() {
		var partContent custcontent.CustomerContent
		partContent.Text = "new content"
		partContent.ContentType.Id = 1

		//create
		(&httprunner.BenchmarkOptions{
			Method:             "POST",
			Route:              "/customer/cms/part",
			ParameterizedRoute: "/customer/cms/part/11000",
			Handler:            CreatePartContent,
			QueryString:        &qs,
			JsonBody:           partContent,
			Runs:               b.N,
		}).RequestBenchmark()

		//get all
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/customer/cms/part",
			ParameterizedRoute: "/customer/cms/part",
			Handler:            AllPartContent,
			QueryString:        &qs,
			JsonBody:           partContent,
			Runs:               b.N,
		}).RequestBenchmark()

		//get unique
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/customer/cms/part",
			ParameterizedRoute: "/customer/cms/part/11000",
			Handler:            UniquePartContent,
			QueryString:        &qs,
			JsonBody:           partContent,
			Runs:               b.N,
		}).RequestBenchmark()

		//delete
		(&httprunner.BenchmarkOptions{
			Method:             "DELETE",
			Route:              "/customer/cms/part",
			ParameterizedRoute: "/customer/cms/part/11000",
			Handler:            DeletePartContent,
			QueryString:        &qs,
			JsonBody:           partContent,
			Runs:               b.N,
		}).RequestBenchmark()
	})
	Convey("Category Content", b, func() {
		var categoryContent custcontent.CustomerContent
		categoryContent.Text = "new content"
		categoryContent.ContentType.Id = 1

		//create
		(&httprunner.BenchmarkOptions{
			Method:             "POST",
			Route:              "/customer/cms/category",
			ParameterizedRoute: "/customer/cms/category/1",
			Handler:            CreateCategoryContent,
			QueryString:        &qs,
			JsonBody:           categoryContent,
			Runs:               b.N,
		}).RequestBenchmark()

		//get all
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/customer/cms/category",
			ParameterizedRoute: "/customer/cms/category",
			Handler:            AllCategoryContent,
			QueryString:        &qs,
			JsonBody:           categoryContent,
			Runs:               b.N,
		}).RequestBenchmark()

		//get unique
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/customer/cms/category",
			ParameterizedRoute: "/customer/cms/category/1",
			Handler:            UniqueCategoryContent,
			QueryString:        &qs,
			JsonBody:           categoryContent,
			Runs:               b.N,
		}).RequestBenchmark()

		//delete
		(&httprunner.BenchmarkOptions{
			Method:             "DELETE",
			Route:              "/customer/cms/category",
			ParameterizedRoute: "/customer/cms/category/1",
			Handler:            DeleteCategoryContent,
			QueryString:        &qs,
			JsonBody:           categoryContent,
			Runs:               b.N,
		}).RequestBenchmark()
	})

	err = apicontextmock.DeMock(dtx)
	if err != nil {
		b.Log(err)
	}
}
