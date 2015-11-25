package geography

import (
	"github.com/curt-labs/API/helpers/apicontextmock"
	"github.com/curt-labs/API/helpers/testThatHttp"
	"github.com/curt-labs/API/models/geography"
	. "github.com/smartystreets/goconvey/convey"

	"encoding/json"
	"testing"
	"time"
)

func TestGeography(t *testing.T) {
	var err error
	var s geography.States
	var c geography.Countries

	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}

	Convey("Testing Geography", t, func() {
		//test get states
		thyme := time.Now()
		testThatHttp.Request("get", "/geography/states", "", "?key="+dtx.APIKey, GetAllStates, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &s)
		So(err, ShouldBeNil)
		So(s, ShouldHaveSameTypeAs, geography.States{})

		//test get countries
		thyme = time.Now()
		testThatHttp.Request("get", "/geography/countries", "", "?key="+dtx.APIKey, GetAllCountries, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, geography.Countries{})

		//test get countries and states
		thyme = time.Now()
		testThatHttp.Request("get", "/geography/countrystates", "", "?key="+dtx.APIKey, GetAllCountriesAndStates, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, geography.Countries{})

	})
	_ = apicontextmock.DeMock(dtx)
}
