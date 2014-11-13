package warranty

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestWarranties(t *testing.T) {
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

		ws, err := GetAllWarranties()
		So(err, ShouldBeNil)
		So(len(ws), ShouldBeGreaterThan, 0)

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
}
