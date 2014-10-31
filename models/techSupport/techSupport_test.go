package techSupport

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTechSupport(t *testing.T) {
	var err error
	var tc TechSupport
	tc.Issue = "This is an issue"
	//make contact
	tc.Contact.Email = "e@e.e"
	tc.Contact.LastName = "techSupport"
	tc.Contact.FirstName = "f"
	tc.Contact.Subject = "s"
	tc.Contact.Message = "m"

	Convey("Test Create TechSupport", t, func() {

		err = tc.Create()
		So(err, ShouldBeNil)
	})
	Convey("Test Get TechSupport", t, func() {

		err = tc.Get()
		So(err, ShouldBeNil)
	})
	Convey("Test Get TechSupport", t, func() {

		ts, err := tc.GetByContact()
		So(err, ShouldBeNil)
		So(len(ts), ShouldBeGreaterThanOrEqualTo, 0)
	})
	Convey("Test Get TechSupport", t, func() {

		allTs, err := GetAllTechSupport()
		So(err, ShouldBeNil)
		So(len(allTs), ShouldBeGreaterThanOrEqualTo, 0)
	})
	Convey("Test Delete TechSupport", t, func() {

		err = tc.Delete()
		So(err, ShouldBeNil)

		err = tc.Contact.Delete()
		So(err, ShouldBeNil)
	})

}
