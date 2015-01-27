package vehicle

import (
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/models/products"
	. "github.com/smartystreets/goconvey/convey"

	"encoding/json"
	"math/rand"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestVehicles(t *testing.T) {

	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}

	Convey("Vehicls", t, func() {
		var l products.Lookup

		ti := time.Now().Second()
		rand.Seed(int64(ti))

		//get all
		qs := make(url.Values, 0)
		qs.Add("key", dtx.APIKey)
		response := httprunner.ParameterizedJsonRequest("POST", "/vehicle", "/vehicle", &qs, &qs, Query)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &l), ShouldEqual, nil)

		if len(l.Years) < 1 {
			return
		}

		//year
		qs.Add("year", strconv.Itoa(l.Years[rand.Intn(len(l.Years))]))
		response = httprunner.ParameterizedJsonRequest("POST", "/vehicle", "/vehicle", &qs, &qs, Query)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &l), ShouldEqual, nil)

		if len(l.Makes) < 1 {
			return
		}

		//make
		qs.Add("make", l.Makes[rand.Intn(len(l.Makes))])
		response = httprunner.ParameterizedJsonRequest("POST", "/vehicle", "/vehicle", &qs, &qs, Query)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &l), ShouldEqual, nil)

		if len(l.Models) < 1 {
			return
		}

		//model
		qs.Add("model", l.Models[rand.Intn(len(l.Models))])
		response = httprunner.ParameterizedJsonRequest("POST", "/vehicle", "/vehicle", &qs, &qs, Query)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &l), ShouldEqual, nil)

		if len(l.Submodels) < 1 {
			So(len(l.Parts), ShouldBeGreaterThan, 0)
			return
		}
		t.Log(l.Vehicle, len(l.Parts))

		//submodel
		qs.Add("submodel", l.Submodels[rand.Intn(len(l.Submodels))])
		response = httprunner.ParameterizedJsonRequest("POST", "/vehicle", "/vehicle", &qs, &qs, Query)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &l), ShouldEqual, nil)

		if len(l.Configurations) < 1 {
			So(len(l.Parts), ShouldBeGreaterThan, 0)
			return
		}

		//configs
		for _, v := range l.Configurations {
			qs.Add(v.Type, v.Options[rand.Intn(len(v.Options))])
		}
		response = httprunner.ParameterizedJsonRequest("POST", "/vehicle", "/vehicle", &qs, &qs, Query)
		So(response.Code, ShouldEqual, 200)
		So(json.Unmarshal(response.Body.Bytes(), &l), ShouldEqual, nil)
		So(len(l.Parts), ShouldBeGreaterThanOrEqualTo, 0)

	})
	_ = apicontextmock.DeMock(dtx)
}
