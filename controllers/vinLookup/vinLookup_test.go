package vinLookup

import (
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/models/products"
	. "github.com/smartystreets/goconvey/convey"

	"encoding/json"
	"net/url"
	"testing"
)

func TestVinLookup(t *testing.T) {

	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}

	Convey("VinLookup", t, func() {
		var l products.Lookup
		vins := []string{"1GTJK34131E957990", "1GJHG39R2X1371269"}

		clean := database.GetCleanDBFlag()

		for _, vin := range vins {
			qs := make(url.Values, 0)
			qs.Add("key", dtx.APIKey)

			//test GetParts
			response := httprunner.ParameterizedRequest("GET", "/:vin", "/"+vin, &qs, nil, GetParts)
			if clean == "" && response.Code == 200 {
				So(response.Code, ShouldEqual, 200)
				So(json.Unmarshal(response.Body.Bytes(), &l), ShouldEqual, nil)
			} else {
				So(response.Code, ShouldEqual, 500)
			}

			//tet GetConfigs
			response = httprunner.ParameterizedRequest("GET", "/configs/:vin", "/configs/"+vin, &qs, nil, GetConfigs)
			So(response.Code, ShouldEqual, 200)
			So(json.Unmarshal(response.Body.Bytes(), &l), ShouldEqual, nil)

		}

	})

	_ = apicontextmock.DeMock(dtx)
}
