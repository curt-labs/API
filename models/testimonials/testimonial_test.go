package testimonials

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTestimonials(t *testing.T) {
	var test Testimonial
	var err error
	Convey("Testing Create Testimonial", t, func() {
		test.Content = "Test Content"
		err = test.Create()
		So(err, ShouldBeNil)
	})
	Convey("Update", t, func() {
		test.Content = "New Content"
		test.Active = true
		test.Approved = true
		err = test.Update()
		So(err, ShouldBeNil)

	})

	Convey("Get testimonial", t, func() {
		err = test.Get()
		So(err, ShouldBeNil)
	})
	Convey("GetAll - No paging", t, func() {
		ts, err := GetAllTestimonials(0, 1, false)
		So(err, ShouldBeNil)
		So(len(ts), ShouldBeGreaterThan, 0)

	})

	Convey("GetAll - Paged", t, func() {
		ts, err := GetAllTestimonials(0, 1, false)

		So(err, ShouldBeNil)
		So(len(ts), ShouldBeGreaterThan, 0)

	})

	Convey("GetAll - randomized", t, func() {
		ts, err := GetAllTestimonials(0, 1, true)
		So(err, ShouldBeNil)
		So(len(ts), ShouldBeGreaterThan, 0)

	})
	Convey("Delete", t, func() {
		err = test.Delete()
		So(err, ShouldBeNil)

	})

}

func BenchmarkGetAllTestimonials(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetAllTestimonials(0, 1, false)
	}
}

func BenchmarkGetTestimonial(b *testing.B) {
	test := setupDummyTestimonial()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		test.Create()
		b.StartTimer()
		test.Get()
		b.StopTimer()
		test.Delete()
	}
}

func BenchmarkCreateTestimonial(b *testing.B) {
	test := setupDummyTestimonial()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		test.Create()
		b.StopTimer()
		test.Delete()
	}
}

func BenchmarkUpdateTestimonial(b *testing.B) {
	test := setupDummyTestimonial()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		test.Create()
		b.StartTimer()
		test.Content = "This is a good test."
		test.Update()
		b.StopTimer()
		test.Delete()
	}
}

func BenchmarkDeleteTestimonial(b *testing.B) {
	test := setupDummyTestimonial()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		test.Create()
		b.StartTimer()
		test.Delete()
	}
}

func setupDummyTestimonial() *Testimonial {
	return &Testimonial{
		Rating:    5,
		Title:     "Test Test",
		Content:   "This is a test.",
		Approved:  true,
		Active:    true,
		FirstName: "TESTER",
		LastName:  "TESTER",
		Location:  "Testville, Oklahoma",
	}
}
