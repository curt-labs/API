package warranty

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	"github.com/curt-labs/GoAPI/models/contact"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestWarranties(t *testing.T) {
	var err error
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	Convey("Testing CRUD", t, func() {
		var w Warranty
		var err error
		w.PartNumber = 123456
		date := time.Now()
		w.Date = &date
		//make contact
		w.Contact.Email = "e@e.e"
		w.Contact.LastName = "l"
		w.Contact.FirstName = "f"
		w.Contact.Type = "t"
		w.Contact.Subject = "s"
		w.Contact.Message = "m"

		err = w.Create()
		So(err, ShouldBeNil)

		err = w.Get()
		So(err, ShouldBeNil)

		wts, err := w.GetByContact()
		So(err, ShouldBeNil)
		So(len(wts), ShouldBeGreaterThan, 0)

		ws, err := GetAllWarranties(MockedDTX)
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(ws), ShouldBeGreaterThan, 0)
		}

		Convey("DELETE", func() {
			err = w.Delete()
			So(err, ShouldBeNil)
			Convey("Cleanup", func() {
				//cleanup contact
				if w.Contact.ID > 0 {
					err = w.Contact.Delete()
					So(err, ShouldBeNil)

				}
			})
		})

	})
	_ = apicontextmock.DeMock(MockedDTX)
}

func BenchmarkGetAllWarranties(b *testing.B) {
	var err error
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	for i := 0; i < b.N; i++ {
		GetAllWarranties(MockedDTX)
	}
	_ = apicontextmock.DeMock(MockedDTX)
}

func BenchmarkGetWarranty(b *testing.B) {
	w := setupDummyWarranty()
	w.Create()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Get()
	}
	b.StopTimer()
	w.Delete()
}

func BenchmarkGetWarrantyByContact(b *testing.B) {
	w := setupDummyWarranty()
	w.Create()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.GetByContact()
	}
	b.StopTimer()
	w.Delete()
}

func BenchmarkCreateWarranty(b *testing.B) {
	w := setupDummyWarranty()
	for i := 0; i < b.N; i++ {
		w.Create()
		b.StopTimer()
		w.Delete()
		b.StartTimer()
	}
}

func BenchmarkDeleteWarranty(b *testing.B) {
	w := setupDummyWarranty()
	w.Create()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Delete()
	}
}

func setupDummyWarranty() *Warranty {
	date := time.Now()
	return &Warranty{
		Date:       &date,
		PartNumber: 999999,
		Contact: contact.Contact{
			Email:     "test@test.com",
			FirstName: "TESTER",
			LastName:  "TESTER",
			Type:      "TESTER",
			Subject:   "TESTER",
			Message:   "This is a test.",
		},
	}
}
