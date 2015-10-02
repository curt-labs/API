package showcase

import (
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/models/showcase"
	. "github.com/smartystreets/goconvey/convey"

	"encoding/json"
	"net/url"
	"strconv"
	"testing"
)

func TestShowcases(t *testing.T) {

	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}

	Convey("Showcases", t, func() {
		var test showcase.Showcase
		var tests []showcase.Showcase

		qs := make(url.Values, 0)
		qs.Add("key", dtx.APIKey)

		test.BrandID = dtx.BrandID
		test.Text = "test content - controller test"

		response := httprunner.ParameterizedJsonRequest("POST", "/showcase", "/showcase", &qs, test, Save)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &test), ShouldEqual, nil)

		test.FirstName = "test name"
		response = httprunner.ParameterizedJsonRequest("PUT", "/showcase", "/showcase", &qs, test, Save)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &test), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/showcase/:id", "/showcase/"+strconv.Itoa(test.ID), &qs, nil, GetShowcase)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &test), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/showcase", "/showcase", &qs, nil, GetAllShowcases)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &tests), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("DELETE", "/showcase/:id", "/showcase/"+strconv.Itoa(test.ID), &qs, nil, Delete)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &test), ShouldEqual, nil)

	})
	_ = apicontextmock.DeMock(dtx)
}
