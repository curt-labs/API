package customer_ctlr_new

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/apiKeyType"
	"github.com/curt-labs/GoAPI/models/customer_new"
	"github.com/curt-labs/GoAPI/models/customer_new/content"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestCustomerContent(t *testing.T) {
	flag.Parse()
	//customer_new - for db setup only
	var c customer_new.Customer
	var cu customer_new.CustomerUser
	var content custcontent.CustomerContent
	var partContent custcontent.PartContent
	var categoryContent custcontent.CustomerContent
	var ct custcontent.ContentType
	var crs custcontent.CustomerContentRevisions
	var contents []custcontent.CustomerContent
	var catContent custcontent.CategoryContent
	var catContents []custcontent.CategoryContent

	catContent.CategoryId = 1

	ct.Type = "test"
	ct.Create()

	c.Name = "test cust"
	c.Create()
	var pub, pri, auth apiKeyType.ApiKeyType
	if database.EmptyDb != nil {
		t.Log("clean db")
		//setup apiKeyTypes
		pub.Type = "Public"
		pri.Type = "Private"
		auth.Type = "Authentication"
		pub.Create()
		pri.Create()
		auth.Create()
	}

	cu.CustomerID = c.Id
	cu.Name = "test cust content user"
	cu.Email = "pretend@test.com"
	cu.Password = "test"
	cu.Sudo = true
	cu.Create()
	var err error
	var apiKey string
	for _, key := range cu.Keys {
		if strings.ToLower(key.Type) == "public" {
			apiKey = key.Key
		}
	}
	t.Log("APIKEY", apiKey)

	// custCon.Save(11000, 1, apiKey)

	Convey("Testing Customer_New/Customer_content", t, func() {
		//test create part content
		content.Text = "new content"
		content.ContentType.Id = 1
		bodyBytes, _ := json.Marshal(content)
		bodyJson := bytes.NewReader(bodyBytes)
		thyme := time.Now()
		testThatHttp.Request("post", "/new/customer/cms/part/", ":id", strconv.Itoa(11000)+"?key="+apiKey, CreatePartContent, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &content)
		So(err, ShouldBeNil)
		So(content, ShouldHaveSameTypeAs, custcontent.CustomerContent{})
		So(content.Id, ShouldBeGreaterThan, 0)

		//create category content
		categoryContent.Text = "new content"
		categoryContent.ContentType.Id = 1
		bodyBytes, _ = json.Marshal(categoryContent)
		bodyJson = bytes.NewReader(bodyBytes)
		thyme = time.Now()
		testThatHttp.Request("post", "/new/customer/cms/category/", ":id", strconv.Itoa(catContent.CategoryId)+"?key="+apiKey, CreateCategoryContent, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &categoryContent)
		So(err, ShouldBeNil)
		So(categoryContent, ShouldHaveSameTypeAs, custcontent.CustomerContent{})
		So(categoryContent.Id, ShouldBeGreaterThan, 0)

		//test update part content
		content.Text = "newerer content"
		bodyBytes, _ = json.Marshal(content)
		bodyJson = bytes.NewReader(bodyBytes)
		thyme = time.Now()
		testThatHttp.Request("put", "/new/customer/cms/part/", ":id", strconv.Itoa(11000)+"?key="+apiKey, UpdatePartContent, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &content)
		So(err, ShouldBeNil)
		So(content, ShouldHaveSameTypeAs, custcontent.CustomerContent{})
		So(content.Id, ShouldBeGreaterThan, 0)

		//test update category content
		categoryContent.Text = "newerer content"
		bodyBytes, _ = json.Marshal(categoryContent)
		bodyJson = bytes.NewReader(bodyBytes)
		thyme = time.Now()
		testThatHttp.Request("put", "/new/customer/cms/part/", ":id", strconv.Itoa(catContent.CategoryId)+"?key="+apiKey, UpdateCategoryContent, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &categoryContent)
		So(err, ShouldBeNil)
		So(categoryContent, ShouldHaveSameTypeAs, custcontent.CustomerContent{})
		So(categoryContent.Id, ShouldBeGreaterThan, 0)

		//test get part content (unique)
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/cms/part/", ":id", strconv.Itoa(11000)+"?key="+apiKey, UniquePartContent, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &contents)
		So(err, ShouldBeNil)
		So(contents, ShouldHaveSameTypeAs, []custcontent.CustomerContent{})

		//test get all part content
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/cms/part", "", "?key="+apiKey, AllPartContent, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &contents)
		So(err, ShouldBeNil)
		So(contents, ShouldHaveSameTypeAs, []custcontent.CustomerContent{})

		//test get category content (all content)
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/cms/part", "", "?key="+apiKey, AllCategoryContent, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &contents)
		So(err, ShouldBeNil)
		So(contents, ShouldHaveSameTypeAs, []custcontent.CustomerContent{})

		//test get unique category content
		catContent.Content = append(catContent.Content, content) //setup some category Content
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/cms/category/", ":id", strconv.Itoa(catContent.CategoryId)+"?key="+apiKey, UniqueCategoryContent, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &catContents)
		So(err, ShouldBeNil)
		So(catContents, ShouldHaveSameTypeAs, []custcontent.CategoryContent{})

		//test get all content
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/cms", "", "?key="+apiKey, GetAllContent, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &contents)
		So(err, ShouldBeNil)
		So(contents, ShouldHaveSameTypeAs, []custcontent.CustomerContent{})

		//test get content by id
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/cms/", ":id", strconv.Itoa(content.Id)+"?key="+apiKey, GetContentById, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &content)
		So(err, ShouldBeNil)
		So(content, ShouldHaveSameTypeAs, custcontent.CustomerContent{})

		//test get content revisions by id
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/cms/", ":id/revisions", strconv.Itoa(content.Id)+"/revisions?key="+apiKey, GetContentRevisionsById, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &crs)
		So(err, ShouldBeNil)
		So(crs, ShouldHaveSameTypeAs, custcontent.CustomerContentRevisions{})

		//test get all content types
		thyme = time.Now()
		testThatHttp.Request("get", "/new/customer/cms/content_types", "", "?key="+apiKey, GetAllContentTypes, nil, "")
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
		testThatHttp.Request("delete", "/new/customer/cms/part/", ":id", strconv.Itoa(11000)+"?key="+apiKey, DeletePartContent, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &partContent)
		So(err, ShouldBeNil)
		So(partContent, ShouldHaveSameTypeAs, custcontent.PartContent{})

		//test delete category content
		bodyBytes, _ = json.Marshal(categoryContent)
		bodyJson = bytes.NewReader(bodyBytes)
		thyme = time.Now()
		testThatHttp.Request("delete", "/new/customer/cms/category/", ":id", strconv.Itoa(catContent.CategoryId)+"?key="+apiKey, DeleteCategoryContent, bodyJson, "application/json")
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

	err = cu.Delete()
	t.Log("cct", err)
	if database.EmptyDb != nil {
		err = pub.Delete()

		err = pri.Delete()

		err = auth.Delete()
	}
}

//using httptestrunner
func TestCreateDeletePartContent(t *testing.T) {
	//get apiKey by creating customeruser
	var cu customer_new.CustomerUser
	var apiKey string
	cu.Name = "test cust content new httprunner user"
	cu.Email = "pretend@test.com"
	cu.Password = "test"
	cu.Sudo = true
	cu.Create()

	for _, key := range cu.Keys {
		if strings.ToLower(key.Type) == "public" {
			apiKey = key.Key
		}
	}
	Convey("Create & Delete Part Content", t, func() {
		var content custcontent.CustomerContent
		content.Text = "new content"
		content.ContentType.Id = 1

		qs := make(url.Values, 0)
		qs.Add("key", apiKey)

		resp := httprunner.JsonRequest("POST", "/new/customer/cms/part/11000", &qs, content, CreatePartContent)
		So(resp.Code, ShouldEqual, 200)
		So(json.Unmarshal(resp.Body.Bytes(), &content), ShouldBeNil)

		resp = httprunner.JsonRequest("DELETE", "/new/customer/cms/part/11000", &qs, content, DeletePartContent)
		So(resp.Code, ShouldEqual, 200)
		So(json.Unmarshal(resp.Body.Bytes(), &content), ShouldBeNil)
	})

	//teardown customer user

	err := cu.Delete()
	t.Log(err)
}

func BenchmarkCRUDContent(b *testing.B) {
	//get apiKey by creating customeruser
	var cu customer_new.CustomerUser
	var apiKey string
	cu.Name = "test cust content benchmark user"
	cu.Email = "pretend@test.com"
	cu.Password = "test"
	cu.Sudo = true
	cu.Create()
	for _, key := range cu.Keys {
		if strings.ToLower(key.Type) == "public" {
			apiKey = key.Key
		}
	}

	qs := make(url.Values, 0)
	qs.Add("key", apiKey)

	Convey("Part Content", b, func() {
		var partContent custcontent.CustomerContent
		partContent.Text = "new content"
		partContent.ContentType.Id = 1

		//create
		(&httprunner.BenchmarkOptions{
			Method:             "POST",
			Route:              "/new/customer/cms/part",
			ParameterizedRoute: "/new/customer/cms/part/11000",
			Handler:            CreatePartContent,
			QueryString:        &qs,
			JsonBody:           partContent,
			Runs:               b.N,
		}).RequestBenchmark()

		//get all
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/new/customer/cms/part",
			ParameterizedRoute: "/new/customer/cms/part",
			Handler:            AllPartContent,
			QueryString:        &qs,
			JsonBody:           partContent,
			Runs:               b.N,
		}).RequestBenchmark()

		//get unique
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/new/customer/cms/part",
			ParameterizedRoute: "/new/customer/cms/part/11000",
			Handler:            UniquePartContent,
			QueryString:        &qs,
			JsonBody:           partContent,
			Runs:               b.N,
		}).RequestBenchmark()

		//delete
		(&httprunner.BenchmarkOptions{
			Method:             "DELETE",
			Route:              "/new/customer/cms/part",
			ParameterizedRoute: "/new/customer/cms/part/11000",
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
			Route:              "/new/customer/cms/category",
			ParameterizedRoute: "/new/customer/cms/category/1",
			Handler:            CreateCategoryContent,
			QueryString:        &qs,
			JsonBody:           categoryContent,
			Runs:               b.N,
		}).RequestBenchmark()

		//get all
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/new/customer/cms/category",
			ParameterizedRoute: "/new/customer/cms/category",
			Handler:            AllCategoryContent,
			QueryString:        &qs,
			JsonBody:           categoryContent,
			Runs:               b.N,
		}).RequestBenchmark()

		//get unique
		(&httprunner.BenchmarkOptions{
			Method:             "GET",
			Route:              "/new/customer/cms/category",
			ParameterizedRoute: "/new/customer/cms/category/1",
			Handler:            UniqueCategoryContent,
			QueryString:        &qs,
			JsonBody:           categoryContent,
			Runs:               b.N,
		}).RequestBenchmark()

		//delete
		(&httprunner.BenchmarkOptions{
			Method:             "DELETE",
			Route:              "/new/customer/cms/category",
			ParameterizedRoute: "/new/customer/cms/category/1",
			Handler:            DeleteCategoryContent,
			QueryString:        &qs,
			JsonBody:           categoryContent,
			Runs:               b.N,
		}).RequestBenchmark()
	})

	//teardown customer user
	cu.Delete()
}
