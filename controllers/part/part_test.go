package part_ctlr

import (
	"bytes"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/apiKeyType"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/curt-labs/GoAPI/models/customer/content"
	"github.com/curt-labs/GoAPI/models/products"
	"github.com/curt-labs/GoAPI/models/video"
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
	p.ID = 10999 //set part number here for use in creating related objects
	var price products.Price
	var cu customer.CustomerUser
	var cat products.Category
	cat.Create()

	//create install sheet content type
	var contentType custcontent.ContentType
	contentType.Type = "installationSheet"
	contentType.Create()

	//create install sheet content
	var installSheetContent products.Content
	installSheetContent.Text = "https://www.curtmfg.com/masterlibrary/16047/installsheet/CM_16021_INS.PDF"
	installSheetContent.ContentType = contentType
	installSheetContent.Create()

	//create video type -- used in creating video during video test
	var vt video.VideoType
	vt.Name = "test type"
	vt.Create()

	//setup apiKeyTypes
	var pub, pri, auth apiKeyType.ApiKeyType
	pub.Type = "public"
	pri.Type = "private"
	auth.Type = "authentication"
	pub.Create()
	pri.Create()
	auth.Create()

	//create customer
	var c customer.Customer
	c.Name = "test man"
	c.Create()

	//creat customer User
	cu.CustomerID = c.Id
	cu.Name = "test cust user"
	cu.Email = "pretend@test.com"
	cu.Password = "test"
	cu.Sudo = true
	err = cu.Create()
	var apiKey string
	for _, key := range cu.Keys {
		if strings.ToLower(key.Type) == "public" {
			apiKey = key.Key
		}
	}
	t.Log("APIKEY", apiKey)

	Convey("TestingParts", t, func() {
		//test create part
		p.Categories = append(p.Categories, cat)
		p.OldPartNumber = "8675309"
		p.ShortDesc = "test part"
		p.Content = append(p.Content, installSheetContent)
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

		var custPrice customer.Price
		custPrice.CustID = c.Id
		custPrice.PartID = p.ID
		custPrice.Create()

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
		testThatHttp.Request("get", "/part/", ":part/prices", strconv.Itoa(p.ID)+"/prices?key="+apiKey, Prices, nil, "")
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

		//test get part install sheet
		thyme = time.Now()
		testThatHttp.Request("get", "/part/", ":part", strconv.Itoa(p.ID)+".pdf?key="+apiKey, InstallSheet, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()*3) //three seconds
		t.Log("Get install sheet benchmark: ", time.Since(thyme))
		So(testThatHttp.Response.Code, ShouldEqual, 200)

		//test get videos
		//create part video
		var partVid products.PartVideo
		partVid.YouTubeVideoId = "11122333XYZ"
		partVid.PartID = p.ID
		partVid.VideoType.ID = vt.ID
		err = partVid.CreatePartVideo()

		thyme = time.Now()
		testThatHttp.Request("get", "/part/videos/", ":part", strconv.Itoa(p.ID)+"?key="+apiKey, Videos, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var ps []products.PartVideo
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &ps)
		So(err, ShouldBeNil)
		So(len(ps), ShouldBeGreaterThan, 0)

		//get active approved reviews
		//create active approved review
		var review products.Review
		review.PartID = p.ID
		review.Rating = 1
		review.Active = true
		review.Approved = true
		err = review.Create()
		thyme = time.Now()
		testThatHttp.Request("get", "/part/reviews/", ":part", strconv.Itoa(p.ID)+"?key="+apiKey, ActiveApprovedReviews, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var reviews products.Reviews
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &reviews)
		So(err, ShouldBeNil)
		So(len(reviews), ShouldBeGreaterThan, 0)
		review.Delete() //teardown - part has FK constraint on review.partID

		//get packaging - no package created in test
		thyme = time.Now()
		testThatHttp.Request("get", "/part/packages/", ":part", strconv.Itoa(p.ID)+"?key="+apiKey, Packaging, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var pack []products.Package
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &pack)
		So(err, ShouldBeNil)
		So(pack, ShouldHaveSameTypeAs, []products.Package{})

		//get content
		thyme = time.Now()
		testThatHttp.Request("get", "/part/content/", ":part", strconv.Itoa(p.ID)+"?key="+apiKey, GetContent, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var content products.Content
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &content)
		So(err, ShouldBeNil)
		So(content, ShouldHaveSameTypeAs, products.Content{})

		//get attributes
		thyme = time.Now()
		testThatHttp.Request("get", "/part/attributes/", ":part", strconv.Itoa(p.ID)+"?key="+apiKey, Attributes, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var attributes []products.Attribute
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &attributes)
		So(err, ShouldBeNil)
		So(attributes, ShouldHaveSameTypeAs, []products.Attribute{})

		//test get images
		thyme = time.Now()
		testThatHttp.Request("get", "/part/images/", ":part", strconv.Itoa(p.ID)+"?key="+apiKey, Images, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var images []products.Image
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &images)
		So(err, ShouldBeNil)
		So(images, ShouldHaveSameTypeAs, []products.Image{})

		//test get vehicles
		thyme = time.Now()
		testThatHttp.Request("get", "/part/vehicles/", ":part", strconv.Itoa(p.ID)+"?key="+apiKey, Vehicles, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var vs []products.Vehicle
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &vs)
		So(err, ShouldBeNil)
		So(vs, ShouldHaveSameTypeAs, []products.Vehicle{})

		//test get related
		thyme = time.Now()
		testThatHttp.Request("get", "/part/related/", ":part", strconv.Itoa(p.ID)+"?key="+apiKey, GetRelated, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var parts []products.Part
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &parts)
		So(err, ShouldBeNil)
		So(parts, ShouldHaveSameTypeAs, []products.Part{})

		//test get
		thyme = time.Now()
		testThatHttp.Request("get", "/part/", ":part", strconv.Itoa(p.ID)+"?key="+apiKey, Get, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var part products.Part
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &part)
		So(err, ShouldBeNil)
		So(part, ShouldHaveSameTypeAs, products.Part{})

		//test latest
		thyme = time.Now()
		testThatHttp.Request("get", "/part/latest", "", "?key="+apiKey, Latest, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()*3) //3 seconds
		t.Log("Get latest parts benchmark: ", time.Since(thyme))
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &parts)
		So(err, ShouldBeNil)
		So(parts, ShouldHaveSameTypeAs, []products.Part{})

		//test featured
		thyme = time.Now()
		testThatHttp.Request("get", "/part/featured", "", "?key="+apiKey, Featured, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()*3) //3 seconds!
		t.Log("Get featured parts benchmark: ", time.Since(thyme))
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &parts)
		So(err, ShouldBeNil)
		So(parts, ShouldHaveSameTypeAs, []products.Part{})

		//test get all parts
		thyme = time.Now()
		testThatHttp.Request("get", "/part", "", "?key="+apiKey, All, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()*6) //6 seconds
		t.Log("Get all parts benchmark: ", time.Since(thyme))
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &parts)
		So(err, ShouldBeNil)
		So(parts, ShouldHaveSameTypeAs, []products.Part{})

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

		//teardown
		custPrice.Delete()
		partVid.DeleteByPart()

	})
	//teardown
	cu.Delete()
	p.Delete()
	cat.Delete()
	pub.Delete()
	pri.Delete()
	auth.Delete()
	contentType.Delete()
	installSheetContent.Delete()
	vt.Delete()

}
