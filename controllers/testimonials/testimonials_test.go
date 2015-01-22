package testimonials

import (
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/models/testimonials"
	. "github.com/smartystreets/goconvey/convey"

	"encoding/json"
	"net/url"
	"strconv"
	"testing"
)

func TestTestTestimonials(t *testing.T) {

	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}

	Convey("Testimonials", t, func() {
		var test testimonials.Testimonial
		var tests testimonials.Testimonials

		qs := make(url.Values, 0)
		qs.Add("key", dtx.APIKey)

		test.BrandID = dtx.BrandID
		test.Content = "test content - controller test"

		response := httprunner.ParameterizedJsonRequest("POST", "/testimonials", "/testimonials", &qs, test, Save)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &test), ShouldEqual, nil)

		test.FirstName = "test name"
		response = httprunner.ParameterizedJsonRequest("PUT", "/testimonials", "/testimonials", &qs, test, Save)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &test), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/testimonials/:id", "/testimonials/"+strconv.Itoa(test.ID), &qs, nil, GetTestimonial)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &test), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/testimonials", "/testimonials", &qs, nil, GetAllTestimonials)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &tests), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("DELETE", "/testimonials/:id", "/testimonials/"+strconv.Itoa(test.ID), &qs, nil, Delete)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &test), ShouldEqual, nil)

	})
	_ = apicontextmock.DeMock(dtx)
}
