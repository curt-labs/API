package brand_ctlr

import (
	// "bytes"
	"encoding/json"
	"github.com/curt-labs/API/helpers/testThatHttp"
	"github.com/curt-labs/API/models/brand"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestBrand(t *testing.T) {
	var b brand.Brand
	var err error
	Convey("Testing Brand", t, func() {
		//test create
		form := url.Values{"name": {"RonCo"}, "code": {"RONCO"}}
		v := form.Encode()
		body := strings.NewReader(v)

		thyme := time.Now()
		testThatHttp.Request("post", "/brand", "", "", CreateBrand, body, "application/x-www-form-urlencoded")
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &b)

		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(err, ShouldBeNil)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		So(b, ShouldHaveSameTypeAs, brand.Brand{})

		//test update
		form = url.Values{"code": {"RONCOandFriends"}}
		v = form.Encode()
		body = strings.NewReader(v)
		thyme = time.Now()
		testThatHttp.Request("put", "/brand/", ":id", strconv.Itoa(b.ID), UpdateBrand, body, "application/x-www-form-urlencoded")
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &b)

		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(err, ShouldBeNil)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		So(b, ShouldHaveSameTypeAs, brand.Brand{})

		//test get
		thyme = time.Now()
		testThatHttp.Request("get", "/brand/", ":id", strconv.Itoa(b.ID), GetBrand, nil, "")
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &b)

		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(err, ShouldBeNil)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		So(b, ShouldHaveSameTypeAs, brand.Brand{})

		//test get all
		thyme = time.Now()
		testThatHttp.Request("get", "/brand", "", "", GetAllBrands, nil, "")
		var bs brand.Brands
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &bs)

		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(err, ShouldBeNil)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		So(bs, ShouldHaveSameTypeAs, brand.Brands{})

		//test delete
		thyme = time.Now()
		testThatHttp.Request("delete", "/brand/", ":id", strconv.Itoa(b.ID), DeleteBrand, nil, "")
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &b)

		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(err, ShouldBeNil)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		So(b, ShouldHaveSameTypeAs, brand.Brand{})
	})
}

func BenchmarkBrands(b *testing.B) {
	testThatHttp.RequestBenchmark(b.N, "GET", "/brand/1", nil, GetBrand)
	testThatHttp.RequestBenchmark(b.N, "GET", "/brand", nil, GetAllBrands)
}
