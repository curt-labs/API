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
		ts, err := GetAllTestimonials()
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(ts), ShouldBeGreaterThan, 0)
		}
	})

}
