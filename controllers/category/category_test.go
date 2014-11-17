package category_ctlr

import (
	"bytes"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/products"
	. "github.com/smartystreets/goconvey/convey"
	// "net/url"
	"strconv"
	// "strings"
	"testing"
)

func TestNews(t *testing.T) {
	var c products.Category
	var cs []products.Category
	var parts []products.Part
	var err error

	//setup
	var cat products.Category
	cat.Title = "test cat"
	cat.Create()

	var sub products.Category
	sub.Title = "sub cat"
	sub.ParentID = cat.ID
	sub.Create()

	var p products.Part
	p.Categories = append(p.Categories, cat)
	p.Create()

	Convey("Testing Category", t, func() {
		//test create
		// form := url.Values{"url": {"test"}, "fileType": {"pdf"}, "website_id": {"1"}}
		// v := form.Encode()
		// body := strings.NewReader(v)
		testThatHttp.Request("get", "/category", "", "", Parents, nil, "")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cs)
		So(err, ShouldBeNil)
		So(cs, ShouldHaveSameTypeAs, []products.Category{})
		So(len(cs), ShouldBeGreaterThan, 0)

		var filterSpecs FilterSpecifications
		filterSpecs.Key = "foo"
		filterSpecs.Values = []string{"bar"}
		bodyBytes, _ := json.Marshal(filterSpecs)
		bodyJson := bytes.NewReader(bodyBytes)
		testThatHttp.Request("post", "/category/", ":id", strconv.Itoa(cat.ID), GetCategory, bodyJson, "application/json")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, products.Category{})

		testThatHttp.Request("get", "/category/", ":id/subs", strconv.Itoa(cat.ID)+"/subs", SubCategories, nil, "application/json")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cs)
		So(err, ShouldBeNil)
		So(cs, ShouldHaveSameTypeAs, []products.Category{})
		So(len(cs), ShouldBeGreaterThan, 0)

		// //TODO - test hangs at 185 on category_ctlr
		// filterSpecs.Key = ""
		// filterSpecs.Values = []string{""}
		// filterArray := make([]FilterSpecifications, 0)
		// filterArray = append(filterArray, filterSpecs)
		// bodyBytes, _ = json.Marshal(filterArray)
		// bodyJson = bytes.NewReader(bodyBytes)
		// testThatHttp.Request("get", "/category/", ":id/parts", strconv.Itoa(cat.ID)+"/parts", GetParts, bodyJson, "application/json")
		// So(testThatHttp.Response.Code, ShouldEqual, 200)
		// err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &parts)
		// So(err, ShouldBeNil)
		// So(parts, ShouldHaveSameTypeAs, []products.Part{})
		// So(len(parts), ShouldBeGreaterThan, 0)

	})
	//teardown
	cat.Delete()
	sub.Delete()
	p.Delete()
}
