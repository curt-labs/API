package testimonials

import (
	"database/sql"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTestimonials(t *testing.T) {
	Convey("Testing CRUD", t, func() {
		var test Testimonial
		var err error
		test.Content = "Test Content"
		err = test.Create()
		So(err, ShouldBeNil)

		test.Content = "New Content"
		err = test.Update()
		So(err, ShouldBeNil)

		err = test.Get()
		So(err, ShouldBeNil)

		err = test.Delete()
		So(err, ShouldBeNil)

	})
	Convey("Testing GetAll", t, func() {
		var ts Testimonials
		var err error
		Convey("GetAll - No paging", func() {
			ts, err = GetAllTestimonials(0, 0, false)
			if err != sql.ErrNoRows {
				So(err, ShouldBeNil)
				So(len(ts), ShouldBeGreaterThan, 0)
			}
		})

		Convey("GetAll - Paged", func() {
			ts, err = GetAllTestimonials(1, 20, false)
			if err != sql.ErrNoRows {
				So(err, ShouldBeNil)
				So(len(ts), ShouldBeGreaterThan, 0)
			}
		})

		Convey("GetAll - randomized", func() {
			ts, err = GetAllTestimonials(1, 2, true)
			if err != sql.ErrNoRows {
				So(err, ShouldBeNil)
				So(len(ts), ShouldBeGreaterThan, 0)
			}
		})
	})
}
