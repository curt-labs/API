package webProperty_controller

import (
	"github.com/curt-labs/API/helpers/apicontextmock"
	"github.com/curt-labs/API/helpers/httprunner"
	"github.com/curt-labs/API/models/webProperty"
	. "github.com/smartystreets/goconvey/convey"

	"encoding/json"
	"net/url"
	"strconv"
	"testing"
)

func TestWebProperty(t *testing.T) {

	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}

	Convey("Test Web Property", t, func() {
		var wt webProperty_model.WebPropertyType
		var wr webProperty_model.WebPropertyRequirement
		var wn webProperty_model.WebPropertyNote
		var w webProperty_model.WebProperty
		var ws webProperty_model.WebProperties
		var wts []webProperty_model.WebPropertyType
		var wrs []webProperty_model.WebPropertyRequirement
		var wns []webProperty_model.WebPropertyNote

		qs := make(url.Values, 0)
		qs.Add("key", dtx.APIKey)

		w.CustID = dtx.CustomerID
		wt.Type = "controller test type"
		w.Name = "controller test name"
		wr.Compliance = true
		wn.Text = "controller test text"

		//POST

		response := httprunner.ParameterizedJsonRequest("POST", "/webProperties/json/type", "/webProperties/json/type", &qs, wt, CreateUpdateWebPropertyType)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &wt), ShouldEqual, nil)

		wt.Type = "999"

		response = httprunner.ParameterizedJsonRequest("POST", "/webProperties/json/type/:id", "/webProperties/json/type/"+strconv.Itoa(wt.ID), &qs, wt, CreateUpdateWebPropertyType)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &wt), ShouldEqual, nil)

		w.WebPropertyType = wt

		response = httprunner.ParameterizedJsonRequest("POST", "/webProperties/json/requirement", "/webProperties/json/requirement", &qs, wr, CreateUpdateWebPropertyRequirement)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &wr), ShouldEqual, nil)

		response = httprunner.ParameterizedJsonRequest("POST", "/webProperties/json/requirement/:id", "/webProperties/json/requirement/"+strconv.Itoa(wr.RequirementID), &qs, wr, CreateUpdateWebPropertyRequirement)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &wr), ShouldEqual, nil)

		w.WebPropertyRequirements = append(w.WebPropertyRequirements, wr)

		response = httprunner.ParameterizedJsonRequest("POST", "/webProperties/json/note", "/webProperties/json/note", &qs, wn, CreateUpdateWebPropertyNote)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &wn), ShouldEqual, nil)

		response = httprunner.ParameterizedJsonRequest("POST", "/webProperties/json/note/:id", "/webProperties/json/note/"+strconv.Itoa(wn.ID), &qs, wn, CreateUpdateWebPropertyNote)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &wn), ShouldEqual, nil)

		w.WebPropertyNotes = append(w.WebPropertyNotes, wn)

		response = httprunner.ParameterizedJsonRequest("POST", "/webProperties/json", "/webProperties/json", &qs, w, CreateUpdateWebProperty)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &w), ShouldEqual, nil)

		response = httprunner.ParameterizedJsonRequest("POST", "/webProperties/json/:id", "/webProperties/json/"+strconv.Itoa(w.ID), &qs, w, CreateUpdateWebProperty)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &w), ShouldEqual, nil)

		//GET
		response = httprunner.ParameterizedRequest("GET", "/webProperties/:id", "/webProperties/"+strconv.Itoa(w.ID), &qs, nil, Get)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &w), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/webProperties", "/webProperties", &qs, nil, GetAll)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &ws), ShouldEqual, nil)
		So(len(ws), ShouldBeGreaterThan, 0)

		response = httprunner.ParameterizedRequest("GET", "/webProperties/type/:id", "/webProperties/type/"+strconv.Itoa(wt.ID), &qs, nil, GetWebPropertyType)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &wt), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/webProperties/type", "/webProperties/type/", &qs, nil, GetAllTypes)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &wts), ShouldEqual, nil)
		So(len(wts), ShouldBeGreaterThan, 0)

		response = httprunner.ParameterizedRequest("GET", "/webProperties/requirement/:id", "/webProperties/requirement/"+strconv.Itoa(wr.RequirementID), &qs, nil, GetWebPropertyRequirement)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &wr), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/webProperties/requirement", "/webProperties/requirement", &qs, nil, GetAllRequirements)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &wrs), ShouldEqual, nil)
		So(len(wrs), ShouldBeGreaterThan, 0)

		response = httprunner.ParameterizedRequest("GET", "/webProperties/note/:id", "/webProperties/note/"+strconv.Itoa(wn.ID), &qs, nil, GetWebPropertyNote)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &wn), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("GET", "/webProperties/note", "/webProperties/note/", &qs, nil, GetAllNotes)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &wns), ShouldEqual, nil)
		So(len(wns), ShouldBeGreaterThan, 0)

		qs.Add("name", "controller test name")

		response = httprunner.ParameterizedRequest("GET", "/webProperties/search", "/webProperties/search/", &qs, nil, Search)
		So(response.Code, ShouldEqual, 200)
		var l interface{}
		So(json.Unmarshal(response.Body.Bytes(), &l), ShouldEqual, nil)

		//DELETE
		response = httprunner.ParameterizedRequest("DELETE", "/webProperties/requirement/:id", "/webProperties/requirement/"+strconv.Itoa(wr.RequirementID), &qs, nil, DeleteWebPropertyRequirement)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &wr), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("DELETE", "/webProperties/note/:id", "/webProperties/note/"+strconv.Itoa(wn.ID), &qs, nil, DeleteWebPropertyNote)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &wn), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("DELETE", "/webProperties/:id", "/webProperties/"+strconv.Itoa(w.ID), &qs, nil, DeleteWebProperty)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &w), ShouldEqual, nil)

		response = httprunner.ParameterizedRequest("DELETE", "/webProperties/type/:id", "/webProperties/type/"+strconv.Itoa(wt.ID), &qs, nil, DeleteWebPropertyType)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &wt), ShouldEqual, nil)

	})

	_ = apicontextmock.DeMock(dtx)
}
