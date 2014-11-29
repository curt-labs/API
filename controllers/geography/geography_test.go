package geography

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/curt-labs/GoAPI/models/geography"
	. "github.com/smartystreets/goconvey/convey"

	"strings"
	"testing"
	"time"
)

func TestGeography(t *testing.T) {
	var err error
	var s geography.States
	var c geography.Countries

	var cu customer.CustomerUser

	cu.Name = "test cust user"
	cu.Email = "pretend@test.com"
	cu.Password = "test"
	cu.Sudo = true
	cu.Create()
	var apiKey string
	for _, key := range cu.Keys {
		if strings.ToLower(key.Type) == "public" {
			apiKey = key.Key
		}
	}
	t.Log("APIKEY", apiKey)

	Convey("Testing Geography", t, func() {
		//test get states
		thyme := time.Now()
		testThatHttp.Request("get", "/geography/states", "", "?key="+apiKey, GetAllStates, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &s)
		So(err, ShouldBeNil)
		So(s, ShouldHaveSameTypeAs, geography.States{})

		//test get countries
		thyme = time.Now()
		testThatHttp.Request("get", "/geography/countries", "", "?key="+apiKey, GetAllCountries, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, geography.Countries{})

		//test get countries and states
		thyme = time.Now()
		testThatHttp.Request("get", "/geography/countrystates", "", "?key="+apiKey, GetAllCountriesAndStates, nil, "")
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &c)
		So(err, ShouldBeNil)
		So(c, ShouldHaveSameTypeAs, geography.Countries{})

	})
	cu.Delete()
}
