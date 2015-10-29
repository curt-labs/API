package category_ctlr

// import (
// 	"bytes"
// 	"encoding/json"
// 	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
// 	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
// 	"github.com/curt-labs/GoAPI/models/products"
// 	. "github.com/smartystreets/goconvey/convey"
// 	// "net/url"
// 	"strconv"
// 	"testing"
// 	"time"
// )

// func TestCategory(t *testing.T) {
// 	var c products.Category
// 	var cs []products.Category
// 	var parts []products.Part
// 	var err error

// 	dtx, err := apicontextmock.Mock()
// 	if err != nil {
// 		t.Log(err)
// 	}

// 	//setup
// 	var cat products.Category
// 	cat.Title = "test cat"
// 	cat.Create()

// 	var sub products.Category
// 	sub.Title = "sub cat"
// 	sub.ParentID = cat.ID
// 	sub.Create()

// 	var p products.Part
// 	p.Categories = append(p.Categories, cat)
// 	p.Create(dtx)

// 	Convey("Testing Category", t, func() {
// 		//test get parents
// 		thyme := time.Now()
// 		testThatHttp.Request("get", "/category", "", "", Parents, nil, "")
// 		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cs)
// 		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
// 		So(testThatHttp.Response.Code, ShouldEqual, 200)
// 		So(err, ShouldBeNil)
// 		So(cs, ShouldHaveSameTypeAs, []products.Category{})
// 		So(len(cs), ShouldBeGreaterThanOrEqualTo, 0)

// 		//test get category
// 		var filterSpecs FilterSpecifications
// 		var filterSpecsArray []FilterSpecifications
// 		filterSpecs.Key = "Weight"
// 		filterSpecs.Values = []string{"33"}
// 		filterSpecsArray = append(filterSpecsArray, filterSpecs)
// 		bodyBytes, _ := json.Marshal(filterSpecsArray)
// 		bodyJson := bytes.NewReader(bodyBytes)

// 		thyme = time.Now()
// 		testThatHttp.Request("post", "/category/", ":id", strconv.Itoa(cat.ID), GetCategory, bodyJson, "application/json")
// 		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
// 		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()*8)
// 		So(err, ShouldBeNil)
// 		So(testThatHttp.Response.Code, ShouldEqual, 200)
// 		So(c, ShouldHaveSameTypeAs, products.Category{})

// 		//test get subcategories
// 		thyme = time.Now()
// 		testThatHttp.Request("get", "/category/", ":id/subs", strconv.Itoa(cat.ID)+"/subs", SubCategories, nil, "application/json")
// 		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cs)

// 		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
// 		So(err, ShouldBeNil)
// 		So(testThatHttp.Response.Code, ShouldEqual, 200)
// 		So(cs, ShouldHaveSameTypeAs, []products.Category{})
// 		So(len(cs), ShouldBeGreaterThan, 0)

// 		//TODO - test hangs at line 670 in parts/category model; same with curl request
// 		testThatHttp.Request("get", "/category/", ":id/parts", strconv.Itoa(cat.ID)+"/parts", GetParts, bodyJson, "")
// 		So(testThatHttp.Response.Code, ShouldEqual, 200)
// 		t.Log(testThatHttp.Response.Body)
// 		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cat)
// 		So(err, ShouldBeNil)
// 		So(cat, ShouldHaveSameTypeAs, products.Category{})
// 		So(len(parts), ShouldBeGreaterThanOrEqualTo, 0)

// 	})
// 	//teardown
// 	cat.Delete(dtx)
// 	sub.Delete(dtx)
// 	p.Delete(dtx)
// }

// func BenchmarkBrands(b *testing.B) {
// 	testThatHttp.RequestBenchmark(b.N, "GET", "/category", nil, Parents)
// 	testThatHttp.RequestBenchmark(b.N, "GET", "/category/1", nil, GetCategory)
// 	testThatHttp.RequestBenchmark(b.N, "GET", "/category/1/subs", nil, SubCategories)
// 	testThatHttp.RequestBenchmark(b.N, "GET", "/category/5/parts", nil, GetParts)
// }
