package contact

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestReceivers(t *testing.T) {
	var err error
	var cr ContactReceiver

	Convey("Testing Gets", t, func() {
		Convey("Testing GetAll()", func() {
			receivers, err := GetAllContactReceivers()
			So(len(receivers), ShouldBeGreaterThanOrEqualTo, 0)
			So(err, ShouldBeNil)
		})

		Convey("Testing Get()", func() {
			Convey("ContactReceiver with ID of 0", func() {
				cr = ContactReceiver{}
				err = cr.Get()

				So(cr.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
			})

			Convey("ContactReceiver with non-zero ID", func() {
				cr = ContactReceiver{ID: 18}
				err = cr.Get()

				So(cr.ID, ShouldNotEqual, 0)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Testing Add/Update/Delete", t, func() {
		ct := ContactType{Name: "TestType"}
		ct.Add()
		cr = ContactReceiver{
			FirstName: "TEST",
			LastName:  "TEST",
			Email:     "test@test.com",
		}
		cr.ContactTypes = append(cr.ContactTypes, ct)

		Convey("Add Missing Values", func() {
			Convey("Missing Email", func() {
				cr.Email = "INVALID"
				err = cr.Add()
				So(cr.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				cr.Email = "test@test.com"
			})
		})
		Convey("Update Missing Values", func() {
			Convey("Missing Email", func() {
				cr.Email = "INVALID"
				err = cr.Update()
				So(cr.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				cr.Email = "test@test.com"
			})
		})

		Convey("Add Valid ContactReceiver", func() {
			err = cr.Add()
			So(cr.ID, ShouldNotEqual, 0)
			So(err, ShouldBeNil)

			Convey("Update Valid ContactReceiver", func() {
				cr.FirstName = "Fred"
				cr.LastName = "Flintstone"
				err = cr.Update()
				So(cr.FirstName, ShouldEqual, "Fred")
				So(cr.LastName, ShouldEqual, "Flintstone")
				So(err, ShouldBeNil)

				Convey("Test Get now that it's valid and has a new type", func() {
					err = cr.Get()
					So(err, ShouldBeNil)

				})

				Convey("Delete Valid ContactReceiver", func() {
					err = cr.Delete()
					So(err, ShouldBeNil)
				})
			})

		})

		Convey("Delete Empty ContactReceiver", func() {
			rec := ContactReceiver{}
			err = rec.Delete()
			So(err, ShouldNotBeNil)
			err = ct.Delete()
			So(err, ShouldBeNil)

		})
	})
}
