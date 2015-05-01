package part_ctlr

import (
	"bytes"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/apiKeyType"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/curt-labs/GoAPI/models/customer/content"
	"github.com/curt-labs/GoAPI/models/products"
	"github.com/curt-labs/GoAPI/models/video"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestParts(t *testing.T) {
	var err error
	var p products.Part
	p.ID = 10999 //set part number here for use in creating related objects
	var price products.Price
	// var cu customer.CustomerUser
	var cat products.Category
	cat.Create()

	//create install sheet content type
	var contentType custcontent.ContentType
	contentType.Type = "InstallationSheet"
	err = contentType.Create()

	//create install sheet content
	var installSheetContent products.Content
	installSheetContent.Text = "https://www.curtmfg.com/masterlibrary/13201/installsheet/CM_13201_INS.PDF"
	installSheetContent.ContentType.Id = contentType.Id
	err = installSheetContent.Create()

	//create video type -- used in creating video during video test
	var vt video.VideoType
	vt.Name = "test type"
	vt.Create()

	//key types
	var pub, pri, auth apiKeyType.ApiKeyType
	if database.GetCleanDBFlag() != "" {
		//setup apiKeyTypes

		pub.Type = "public"
		pri.Type = "private"
		auth.Type = "authentication"
		pub.Create()
		pri.Create()
		auth.Create()
	}

	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}

	Convey("TestingParts", t, func() {
		//test create part
		p.Categories = append(p.Categories, cat)
		p.OldPartNumber = "8675309"
		p.ShortDesc = "test part"
		p.Content = append(p.Content, installSheetContent)
		bodyBytes, _ := json.Marshal(p)
		bodyJson := bytes.NewReader(bodyBytes)
		thyme := time.Now()
		testThatHttp.RequestWithDtx("post", "/part", "", "?key="+dtx.APIKey, CreatePart, bodyJson, "application/json", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &p)
		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, products.Part{})
		So(p.ID, ShouldEqual, 10999)

		p.BindCustomer(dtx)

		var custPrice customer.Price
		custPrice.CustID = dtx.CustomerID
		custPrice.PartID = p.ID
		err = custPrice.Create()

		//test create price
		price.Price = 987
		price.PartId = p.ID
		price.Type = "test"
		bodyBytes, _ = json.Marshal(price)
		bodyJson = bytes.NewReader(bodyBytes)
		thyme = time.Now()
		testThatHttp.RequestWithDtx("post", "/price", "", "?key="+dtx.APIKey, SavePrice, bodyJson, "application/json", dtx)
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
		testThatHttp.RequestWithDtx("post", "/price/", ":id", strconv.Itoa(price.Id)+"?key="+dtx.APIKey, SavePrice, bodyJson, "application/json", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &price)
		So(err, ShouldBeNil)
		So(price, ShouldHaveSameTypeAs, products.Price{})
		So(price.Type, ShouldNotEqual, "test")

		// //test get part prices
		thyme = time.Now()
		testThatHttp.RequestWithDtx("get", "/part/", ":part/prices", strconv.Itoa(p.ID)+"/prices?key="+dtx.APIKey, Prices, nil, "", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var prices []products.Price
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &prices)
		So(err, ShouldBeNil)
		So(prices, ShouldHaveSameTypeAs, []products.Price{})

		// //test get part categories
		thyme = time.Now()
		testThatHttp.Request("get", "/part/", ":part/categories", strconv.Itoa(p.ID)+"/categories?key="+dtx.APIKey, Categories, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var cats []products.Category
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &cats)
		So(err, ShouldBeNil)
		So(cats, ShouldHaveSameTypeAs, []products.Category{})

		//test get part install sheet
		thyme = time.Now()
		testThatHttp.Request("get", "/part/", ":part", strconv.Itoa(p.ID)+".pdf?key="+dtx.APIKey, InstallSheet, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()*6) //three seconds
		So(testThatHttp.Response.Code, ShouldEqual, 200)

		//test get videos
		//create part video
		var partVid video.Video

		err = partVid.Create(dtx)

		err = partVid.CreateJoinPart(p.ID)

		thyme = time.Now()
		testThatHttp.Request("get", "/part/videos/", ":part", strconv.Itoa(p.ID)+"?key="+dtx.APIKey, Videos, nil, "")
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
		err = review.Create(dtx)
		thyme = time.Now()
		testThatHttp.Request("get", "/part/reviews/", ":part", strconv.Itoa(p.ID)+"?key="+dtx.APIKey, ActiveApprovedReviews, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var reviews products.Reviews
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &reviews)
		So(err, ShouldBeNil)
		So(len(reviews), ShouldBeGreaterThan, 0)
		review.Delete(dtx) //teardown - part has FK constraint on review.partID

		//get packaging - no package created in test
		thyme = time.Now()
		testThatHttp.Request("get", "/part/packages/", ":part", strconv.Itoa(p.ID)+"?key="+dtx.APIKey, Packaging, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var pack []products.Package
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &pack)
		So(err, ShouldBeNil)
		So(pack, ShouldHaveSameTypeAs, []products.Package{})

		//get content
		thyme = time.Now()
		testThatHttp.Request("get", "/part/content/", ":part", strconv.Itoa(p.ID)+"?key="+dtx.APIKey, GetContent, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var content products.Content
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &content)
		So(err, ShouldBeNil)
		So(content, ShouldHaveSameTypeAs, products.Content{})

		//get attributes
		thyme = time.Now()
		testThatHttp.Request("get", "/part/attributes/", ":part", strconv.Itoa(p.ID)+"?key="+dtx.APIKey, Attributes, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var attributes []products.Attribute
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &attributes)
		So(err, ShouldBeNil)
		So(attributes, ShouldHaveSameTypeAs, []products.Attribute{})

		//test get images
		thyme = time.Now()
		testThatHttp.Request("get", "/part/images/", ":part", strconv.Itoa(p.ID)+"?key="+dtx.APIKey, Images, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var images []products.Image
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &images)
		So(err, ShouldBeNil)
		So(images, ShouldHaveSameTypeAs, []products.Image{})

		//test get vehicles
		thyme = time.Now()
		testThatHttp.Request("get", "/part/vehicles/", ":part", strconv.Itoa(p.ID)+"?key="+dtx.APIKey, Vehicles, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var vs []products.Vehicle
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &vs)
		So(err, ShouldBeNil)
		So(vs, ShouldHaveSameTypeAs, []products.Vehicle{})

		//test get related
		thyme = time.Now()
		testThatHttp.Request("get", "/part/related/", ":part", strconv.Itoa(p.ID)+"?key="+dtx.APIKey, GetRelated, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var parts []products.Part
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &parts)
		So(err, ShouldBeNil)
		So(parts, ShouldHaveSameTypeAs, []products.Part{})

		//test get
		thyme = time.Now()
		t.Log(strconv.Itoa(p.ID))
		testThatHttp.Request("get", "/part/", ":part", strconv.Itoa(p.ID)+"?key="+dtx.APIKey, Get, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds())
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var part products.Part
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &part)
		So(err, ShouldBeNil)
		So(part, ShouldHaveSameTypeAs, products.Part{})

		//test latest
		thyme = time.Now()
		testThatHttp.Request("get", "/part/latest", "", "?key="+dtx.APIKey, Latest, nil, "")
		t.Log("Get latest parts benchmark: ", time.Since(thyme))
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &parts)
		So(err, ShouldBeNil)
		So(parts, ShouldHaveSameTypeAs, []products.Part{})

		// test all
		thyme = time.Now()
		testThatHttp.Request("get", "/part/all", "", "?key="+dtx.APIKey, AllBasics, nil, "")
		t.Log("Get all basic parts benchmark: ", time.Since(thyme))
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &parts)
		So(err, ShouldBeNil)
		So(parts, ShouldHaveSameTypeAs, []products.Part{})

		//test featured
		thyme = time.Now()
		testThatHttp.Request("get", "/part/featured", "", "?key="+dtx.APIKey, Featured, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()*3) //3 seconds!
		t.Log("Get featured parts benchmark: ", time.Since(thyme))
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &parts)
		So(err, ShouldBeNil)
		So(parts, ShouldHaveSameTypeAs, []products.Part{})

		//test get all parts
		thyme = time.Now()
		testThatHttp.Request("get", "/part", "", "?key="+dtx.APIKey, All, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()*6) //6 seconds
		t.Log("Get all parts benchmark: ", time.Since(thyme))
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &parts)
		So(err, ShouldBeNil)
		So(parts, ShouldHaveSameTypeAs, []products.Part{})

		//test get price
		thyme = time.Now()
		testThatHttp.Request("get", "/price/", ":id", strconv.Itoa(price.Id)+"?key="+dtx.APIKey, GetPrice, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &price)
		So(err, ShouldBeNil)
		So(price, ShouldHaveSameTypeAs, products.Price{})

		//test get old part Number
		thyme = time.Now()
		testThatHttp.Request("get", "/part/old/", ":part", p.OldPartNumber+"?key="+dtx.APIKey, OldPartNumber, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &p)
		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, products.Part{})
		So(p.OldPartNumber, ShouldEqual, "8675309")
		So(p.ID, ShouldEqual, 10999)

		//test delete price
		thyme = time.Now()
		testThatHttp.Request("delete", "/price/", ":id", strconv.Itoa(price.Id)+"?key="+dtx.APIKey, DeletePrice, nil, "")
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
		testThatHttp.Request("put", "/part/", ":id", strconv.Itoa(p.ID)+"?key="+dtx.APIKey, UpdatePart, bodyJson, "application/json")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds())
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &p)
		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, products.Part{})
		So(p.OldPartNumber, ShouldEqual, "8675309")
		So(p.ID, ShouldEqual, 10999)

		//test delete part
		thyme = time.Now()
		testThatHttp.Request("delete", "/part/", ":id", strconv.Itoa(p.ID)+"?key="+dtx.APIKey, DeletePart, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &p)
		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, products.Part{})

		// //teardown
		custPrice.Delete()
		partVid.Delete(dtx)
		partVid.DeleteJoinPart(p.ID)

	})
	//teardown
	p.Delete(dtx)
	cat.Delete(dtx)
	if database.GetCleanDBFlag() != "" {
		pub.Delete()
		pri.Delete()
		auth.Delete()
	}
	contentType.Delete()
	installSheetContent.Delete()
	vt.Delete()

	_ = apicontextmock.DeMock(dtx)

}
