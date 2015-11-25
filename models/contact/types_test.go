package contact

import (
	"database/sql"
	"github.com/curt-labs/API/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTypes(t *testing.T) {
	var err error
	var ct ContactType
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}

	Convey("Testing Add/Update/Delete", t, func() {
		ct = ContactType{
			Name: "TEST",
		}
		Convey("Add Missing Values", func() {
			Convey("Missing Name", func() {
				ct.Name = ""
				err = ct.Add()
				So(ct.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				ct.Name = "TEST"
			})
		})
		Convey("Update Missing Values", func() {
			Convey("Missing Name", func() {
				ct.Name = ""
				err = ct.Update()
				So(ct.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				ct.Name = "TEST"
			})
		})

		Convey("Add Valid ContactType", func() {
			err = ct.Add()
			So(ct.ID, ShouldNotEqual, 0)
			So(err, ShouldBeNil)

			Convey("Update Valid ContactType", func() {
				ct.Name = "TESTER"
				err = ct.Update()
				So(ct.Name, ShouldEqual, "TESTER")
				So(err, ShouldBeNil)

				Convey("Test getReceiversByType", func() {
					var cr ContactReceiver
					cr.LastName = "testLName"
					cr.Email = "testEmail@test.com"
					err = cr.Add()
					So(err, ShouldBeNil)

					err = cr.Get()

					Convey("Test Join", func() {
						err = cr.CreateTypeJoin(ct)
						crs, err := ct.GetReceivers()
						if err != sql.ErrNoRows {
							So(err, ShouldBeNil)
							So(len(crs), ShouldBeGreaterThanOrEqualTo, 1)
						}
						//cleanup
						cr.DeleteTypeJoin(ct)
						So(err, ShouldBeNil)
						cr.Delete()
						So(err, ShouldBeNil)

					})

				})

				Convey("Delete Valid ContactType", func() {
					err = ct.Delete()
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("Delete Empty ContactType", func() {
			ctype := ContactType{}
			err = ctype.Delete()
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Testing Gets", t, func() {
		Convey("Testing GetAll()", func() {
			types, err := GetAllContactTypes(MockedDTX)
			So(len(types), ShouldBeGreaterThanOrEqualTo, 0)
			if err != sql.ErrNoRows {
				So(err, ShouldBeNil)
			}
		})

		Convey("Testing Get()", func() {
			Convey("ContactType with ID of 0", func() {
				ct = ContactType{}
				err = ct.Get()

				So(ct.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
			})

			Convey("ContactType with non-zero ID", func() {
				ct2 := ContactType{ID: 1}
				ct2.Name = "TESTER"
				ct2.ShowOnWebsite = true
				ct2.Add()
				err = ct2.Get()

				So(ct2.ID, ShouldNotEqual, 0)
				So(err, ShouldBeNil)
			})
		})
	})
	_ = apicontextmock.DeMock(MockedDTX)
}

func BenchmarkGetAllContactTypes(b *testing.B) {
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}

	for i := 0; i < b.N; i++ {
		GetAllContactTypes(MockedDTX)
	}
	_ = apicontextmock.DeMock(MockedDTX)
}

func BenchmarkGetContactType(b *testing.B) {
	ct := setupDummyContactType()
	ct.Add()
	for i := 0; i < b.N; i++ {
		ct.Get()
	}
	ct.Delete()
}

func BenchmarkAddContactType(b *testing.B) {
	ct := setupDummyContactType()
	for i := 0; i < b.N; i++ {
		ct.Add()
		b.StopTimer()
		ct.Delete()
		b.StartTimer()
	}
}

func BenchmarkUpdateContactType(b *testing.B) {
	ct := setupDummyContactType()
	ct.Add()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ct.Name = "TESTING"
		ct.Update()
	}
	b.StopTimer()
	ct.Delete()
}

func BenchmarkDeleteContactType(b *testing.B) {
	ct := setupDummyContactType()
	ct.Add()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ct.Delete()
	}
}

func setupDummyContactType() *ContactType {
	return &ContactType{
		Name:          "TESTER",
		ShowOnWebsite: false,
	}
}
