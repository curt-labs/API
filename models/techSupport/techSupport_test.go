package techSupport

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTechSupport(t *testing.T) {
	Convey("Test Create TechSupport", t, func() {
		var err error
		var tc TechSupport
		tc.Issue = "This is an issue"

		//make contact
		tc.Contact.Email = "e@e.e"
		tc.Contact.LastName = "techSupport"
		tc.Contact.FirstName = "f"
		tc.Contact.Type = "t"
		tc.Contact.Subject = "s"
		tc.Contact.Message = "m"
		err = tc.Contact.Add()
		So(err, ShouldBeNil)
		err = tc.Contact.Get()
		So(err, ShouldBeNil)

		err = tc.Create()
		So(err, ShouldBeNil)

		err = tc.Get()
		So(err, ShouldBeNil)

		ts, err := tc.GetByContact()
		So(err, ShouldBeNil)
		So(len(ts), ShouldBeGreaterThanOrEqualTo, 1)

		allTs, err := GetAllTechSupport()
		So(err, ShouldBeNil)
		So(len(allTs), ShouldBeGreaterThanOrEqualTo, 1)

		err = tc.Delete()
		So(err, ShouldBeNil)

		err = tc.Contact.Delete()
		So(err, ShouldBeNil)

	})
}
