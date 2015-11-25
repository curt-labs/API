package warranty

import (
	"github.com/curt-labs/API/helpers/apicontextmock"
	"github.com/curt-labs/API/helpers/httprunner"
	"github.com/curt-labs/API/models/products"
	"github.com/curt-labs/API/models/warranty"
	. "github.com/smartystreets/goconvey/convey"

	"encoding/json"
	"net/url"
	"strconv"
	"testing"
)

func TestWarranties(t *testing.T) {
	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}

	Convey("Video Warranties", t, func() {
		var w warranty.Warranty
		var ws []warranty.Warranty
		var p products.Part

		//part setup
		p.ID = 8675309
		p.BrandID = dtx.BrandID
		p.Create(dtx)

		qs := make(url.Values, 0)
		qs.Add("key", dtx.APIKey)

		w.SerialNumber = "controller test serial num"
		w.Contact.FirstName = "Joseph"
		w.Contact.LastName = "Smith"
		w.Contact.Email = "mormons@aregreat.com"
		w.PartNumber = p.ID

		response := httprunner.ParameterizedJsonRequest("POST", "/warranty/:contactReceiverTypeID/:sendEmail", "/warranty/15/false", &qs, w, CreateWarranty)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &w), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/warranty/:id", "/warranty/"+strconv.Itoa(w.ID), &qs, nil, GetWarranty)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &w), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/warranty/contact/:id", "/warranty/contact/"+strconv.Itoa(w.Contact.ID), &qs, nil, GetWarrantyByContact)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &ws), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/warranty/all", "/warranty/all", &qs, nil, GetAllWarranties)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &ws), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("DELETE", "/warranty/:id", "/warranty/"+strconv.Itoa(w.ID), &qs, nil, DeleteWarranty)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &w), ShouldEqual, nil)

		p.Delete(dtx)
	})

	_ = apicontextmock.DeMock(dtx)
}
