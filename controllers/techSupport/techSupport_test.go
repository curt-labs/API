package techSupport

import (
	"github.com/curt-labs/API/helpers/apicontextmock"
	"github.com/curt-labs/API/helpers/httprunner"
	"github.com/curt-labs/API/models/techSupport"
	. "github.com/smartystreets/goconvey/convey"

	"encoding/json"
	"net/url"
	"strconv"
	"testing"
)

func TestTechSupport(te *testing.T) {

	dtx, err := apicontextmock.Mock()
	if err != nil {
		te.Log(err)
	}

	Convey("TechSupport", te, func() {
		var t techSupport.TechSupport
		var ts []techSupport.TechSupport

		qs := make(url.Values, 0)
		qs.Add("key", dtx.APIKey)

		t.Contact.LastName = "test contact"
		t.Contact.FirstName = "test contact first"
		t.Contact.Email = "test@emai.com"
		t.Contact.Type = "test"
		t.Contact.Subject = "test"
		t.Contact.Message = "test"
		t.Contact.Brand.ID = dtx.BrandID
		t.BrandID = dtx.BrandID

		err = t.Contact.Add(dtx)
		So(err, ShouldBeNil)

		response := httprunner.ParameterizedJsonRequest("POST", "/techSupport/:contactReceiverTypeID/:sendEmail", "/techSupport/1/false", &qs, t, CreateTechSupport)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &t), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/techSupport/:id", "/techSupport/"+strconv.Itoa(t.ID), &qs, nil, GetTechSupport)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &t), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/techSupport/contact/:id", "/techSupport/contact/"+strconv.Itoa(t.Contact.ID), &qs, nil, GetTechSupportByContact)
		te.Log(response.Body)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &ts), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/techSupport", "/techSupport", &qs, nil, GetAllTechSupport)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &ts), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("DELETE", "/techSupport/:id", "/techSupport/"+strconv.Itoa(t.ID), &qs, nil, DeleteTechSupport)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &t), ShouldEqual, nil)

		err = t.Contact.Delete()
		So(err, ShouldBeNil)

	})
	_ = apicontextmock.DeMock(dtx)
}
