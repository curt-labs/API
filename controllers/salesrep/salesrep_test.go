package salesrep

import (
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/models/salesrep"
	. "github.com/smartystreets/goconvey/convey"

	"encoding/json"
	"net/url"
	"strconv"
	"testing"
)

func TestSalesRep(t *testing.T) {

	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}

	Convey("SalesRep", t, func() {
		var s salesrep.SalesRep
		var ss salesrep.SalesReps
		qs := make(url.Values, 0)
		qs.Add("key", dtx.APIKey)
		qs.Add("name", "Tim")
		qs.Add("code", "red")

		response := httprunner.ParameterizedRequest("POST", "/salesrep", "/salesrep", &qs, &qs, AddSalesRep)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &s), ShouldEqual, nil)

		qs.Set("name", "Bill")
		response = httprunner.ParameterizedRequest("PUT", "/salesrep/:id", "/salesrep/"+strconv.Itoa(s.ID), &qs, &qs, AddSalesRep)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &s), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/salesrep/:id", "/salesrep/"+strconv.Itoa(s.ID), &qs, nil, GetSalesRep)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &s), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/salesrep", "/salesrep", &qs, nil, GetAllSalesReps)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &ss), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("DELETE", "/salesrep/:id", "/salesrep/"+strconv.Itoa(s.ID), &qs, nil, DeleteSalesRep)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &s), ShouldEqual, nil)

	})
	_ = apicontextmock.DeMock(dtx)
}
