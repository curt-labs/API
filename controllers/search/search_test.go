package search_ctlr

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	. "github.com/smartystreets/goconvey/convey"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	response httptest.ResponseRecorder
)

func TestSearch(t *testing.T) {
	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}

	qs := make(url.Values, 0)
	qs.Add("key", dtx.APIKey)

	Convey("Testing Search with empty term", t, func() {
		response = httprunner.Req(Search, "GET", "/search", "/search", &qs)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
	})
	Convey("Testing Search with `Hitch`", t, func() {
		response = httprunner.Req(Search, "GET", "/search/:term", "/search/Hitch", &qs)
		So(response.Code, ShouldEqual, 200)
	})
}
