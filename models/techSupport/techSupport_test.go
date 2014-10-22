package techSupport

import (
	// "database/sql"
	// "github.com/curt-labs/GoAPI/models/contact"
	. "github.com/smartystreets/goconvey/convey"
	// "math/rand"
	"testing"
)

func TestTechSupport(t *testing.T) {
	Convey("Test Create TechSupport", t, func() {
		var err error
		var tc TechSupport
		tc.Issue = "This is an issue"

		//make contact
		tc.Contact.Email = "e@e.e"
		tc.Contact.LastName = "l"
		tc.Contact.FirstName = "f"
		tc.Contact.Type = "t"
		tc.Contact.Subject = "s"
		tc.Contact.Message = "m"
		err = tc.Contact.Add()
		So(err, ShouldBeNil)
		err = tc.Contact.Get()
		So(err, ShouldBeNil)
		t.Log(tc.Contact)

		err = tc.Create()
		So(err, ShouldBeNil)

		err = tc.Get()
		t.Log(tc)
		So(err, ShouldBeNil)

		ts, err := GetAllTechSupportByContact(tc.Contact)
		So(err, ShouldBeNil)
		So(len(ts), ShouldBeGreaterThanOrEqualTo, 1)

		allTs, err := GetAllTechSupport()
		So(err, ShouldBeNil)
		So(len(allTs), ShouldBeGreaterThanOrEqualTo, 1)

		tc.Contact.Delete()
		err = tc.Delete()
		So(err, ShouldBeNil)

	})
}
