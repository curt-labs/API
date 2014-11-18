package customer_ctlr_new

import (
	"bytes"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/customer_new"
	"github.com/curt-labs/GoAPI/models/customer_new/content"
	. "github.com/smartystreets/goconvey/convey"
	// "net/url"
	"strconv"
	"strings"
	"testing"
)

func TestCustomer(t *testing.T) {

	//customer_new - for db setup only
	var c customer_new.Customer
	var cu customer_new.CustomerUser
	var content custcontent.CustomerContent
	var partContent custcontent.PartContent
	var ct custcontent.ContentType
	ct.Type = "test"
	ct.Create()

	c.Name = "test cust"
	c.Create()
	cu.CustomerID = c.Id
	cu.Name = "test cust user"
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

	Convey("Testing Customer_New", t, func() {
		//test create part content
		content.Text = "new content"
		content.ContentType.Id = 1
		bodyBytes, _ := json.Marshal(content)
		bodyJson := bytes.NewReader(bodyBytes)
		testThatHttp.Request("post", "/new/customer/cms/part/", ":id", strconv.Itoa(11000)+"?key="+apiKey, CreatePartContent, bodyJson, "application/json")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &content)
		So(err, ShouldBeNil)
		So(content, ShouldHaveSameTypeAs, custcontent.CustomerContent{})
		So(content.Id, ShouldBeGreaterThan, 0)

		//test update part content
		content.Text = "newerer content"
		bodyBytes, _ = json.Marshal(content)
		bodyJson = bytes.NewReader(bodyBytes)
		testThatHttp.Request("put", "/new/customer/cms/part/", ":id", strconv.Itoa(11000)+"?key="+apiKey, UpdatePartContent, bodyJson, "application/json")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &content)
		So(err, ShouldBeNil)
		So(content, ShouldHaveSameTypeAs, custcontent.CustomerContent{})
		So(content.Id, ShouldBeGreaterThan, 0)

		//test get part content (unique)
		testThatHttp.Request("get", "/new/customer/cms/part/", ":id", strconv.Itoa(11000)+"?key="+apiKey, UniquePartContent, nil, "")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var contents []custcontent.CustomerContent
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &contents)
		So(err, ShouldBeNil)
		So(contents, ShouldHaveSameTypeAs, []custcontent.CustomerContent{})

		//test get category content (all content)
		testThatHttp.Request("get", "/new/customer/cms/part", "", "?key="+apiKey, AllCategoryContent, nil, "")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &contents)
		So(err, ShouldBeNil)
		So(contents, ShouldHaveSameTypeAs, []custcontent.CustomerContent{})

		//test get all content
		testThatHttp.Request("get", "/new/customer/cms", "", "?key="+apiKey, GetAllContent, nil, "")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &contents)
		So(err, ShouldBeNil)
		So(contents, ShouldHaveSameTypeAs, []custcontent.CustomerContent{})

		//test get all content types
		testThatHttp.Request("get", "/new/customer/cms/content_types", "", "?key="+apiKey, GetAllContentTypes, nil, "")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var cts []custcontent.ContentType
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cts)
		So(err, ShouldBeNil)
		So(cts, ShouldHaveSameTypeAs, []custcontent.ContentType{})
		So(len(cts), ShouldBeGreaterThan, 0)

		//test delete part content
		bodyBytes, _ = json.Marshal(content)
		bodyJson = bytes.NewReader(bodyBytes)
		testThatHttp.Request("delete", "/new/customer/cms/part/", ":id", strconv.Itoa(11000)+"?key="+apiKey, DeletePartContent, bodyJson, "application/json")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &partContent)
		So(err, ShouldBeNil)
		So(partContent, ShouldHaveSameTypeAs, custcontent.PartContent{})

	})
	//teardown
	c.Delete()
	cu.Delete()
	ct.Delete()

	// custCon.Delete(11000, 1, apiKey)
}
