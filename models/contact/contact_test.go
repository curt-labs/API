package contact

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestContacts(t *testing.T) {
	var err error
	var c Contact

	Convey("Testing Add/Update/Delete", t, func() {
		c = Contact{
			FirstName: "TEST",
			LastName:  "TEST",
			Email:     "test@test.com",
			Type:      "TEST",
			Subject:   "TEST",
			Message:   "Testing this awesome code!",
		}
		Convey("Add Missing Values", func() {
			Convey("Missing First Name", func() {
				c.FirstName = ""
				err = c.Add()
				So(c.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				c.FirstName = "TEST"
			})
			Convey("Missing Last Name", func() {
				c.LastName = ""
				err = c.Add()
				So(c.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				c.LastName = "TEST"
			})
			Convey("Bad Email", func() {
				c.Email = "INVALID"
				err = c.Add()
				So(c.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				c.Email = "test@test.com"
			})
			Convey("Missing Type", func() {
				c.Type = ""
				err = c.Add()
				So(c.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				c.Type = "TEST"
			})
			Convey("Missing Subject", func() {
				c.Subject = ""
				err = c.Add()
				So(c.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				c.Subject = "TEST"
			})
			Convey("Missing Message", func() {
				c.Message = ""
				err = c.Add()
				So(c.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				c.Message = "Testing this awesome code!"
			})
			Convey("Empty Contact", func() {
				con := Contact{}
				err = con.Add()
				So(c.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("Update Missing Values", func() {
			Convey("Missing First Name", func() {
				c.FirstName = ""
				err = c.Update()
				So(c.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				c.FirstName = "TEST"
			})
			Convey("Missing Last Name", func() {
				c.LastName = ""
				err = c.Update()
				So(c.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				c.LastName = "TEST"
			})
			Convey("Bad Email", func() {
				c.Email = "INVALID"
				err = c.Update()
				So(c.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				c.Email = "test@test.com"
			})
			Convey("Missing Type", func() {
				c.Type = ""
				err = c.Update()
				So(c.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				c.Type = "TEST"
			})
			Convey("Missing Subject", func() {
				c.Subject = ""
				err = c.Update()
				So(c.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				c.Subject = "TEST"
			})
			Convey("Missing Message", func() {
				c.Message = ""
				err = c.Update()
				So(c.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				c.Message = "Testing this awesome code!"
			})
			Convey("Empty Contact", func() {
				con := Contact{}
				err = con.Update()
				So(c.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("Add Valid Contact", func() {
			err = c.Add()
			So(c.ID, ShouldNotEqual, 0)
			So(err, ShouldBeNil)

			Convey("Update Valid Contact", func() {
				c.FirstName = "FRED"
				c.LastName = "FLINTSTONE"
				c.Message = "There was this one time at band camp...we played the drums."
				err = c.Update()
				So(c.FirstName, ShouldEqual, "FRED")
				So(err, ShouldBeNil)

				Convey("Delete Valid Contact", func() {
					err = c.Delete()
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("Delete Empty Contact", func() {
			con := Contact{}
			err = con.Delete()
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Testing Gets", t, func() {
		Convey("Testing GetAll()", func() {
			contacts, err := GetAllContacts(1, 1)
			So(len(contacts), ShouldBeGreaterThanOrEqualTo, 0)
			So(err, ShouldBeNil)
		})

		Convey("Testing Get()", func() {
			Convey("Contact with ID of 0", func() {
				c = Contact{}
				err = c.Get()

				So(c.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
			})

			Convey("Contact with non-zero ID", func() {
				c = Contact{ID: 1}
				err = c.Get()

				So(c.ID, ShouldNotEqual, 0)
				So(err, ShouldBeNil)
			})
		})
	})
}

func BenchmarkGetAllContacts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetAllContacts(1, 1)
	}
}

func BenchmarkGetContact(b *testing.B) {
	c := setupDummyContact()
	c.Add()
	for i := 0; i < b.N; i++ {
		c.Get()
	}
	c.Delete()
}

func BenchmarkAddContact(b *testing.B) {
	c := setupDummyContact()
	for i := 0; i < b.N; i++ {
		c.Add()
		b.StopTimer()
		c.Delete()
		b.StartTimer()
	}
}

func BenchmarkUpdateContact(b *testing.B) {
	c := setupDummyContact()
	c.Add()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.FirstName = "TESTER"
		c.LastName = "TESTER"
		c.Update()
	}
	b.StopTimer()
	c.Delete()
}

func BenchmarkDeleteContact(b *testing.B) {
	c := setupDummyContact()
	c.Add()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Delete()
	}
}

func setupDummyContact() *Contact {
	return &Contact{
		FirstName: "TEST",
		LastName:  "TEST",
		Email:     "test@test.com",
		Type:      "TEST",
		Subject:   "TEST",
		Message:   "Testing this awesome code!",
	}
}
