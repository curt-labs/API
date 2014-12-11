package testimonials

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var (
	MockedDTX = &apicontext.DataContext{BrandID: 1, WebsiteID: 1, APIKey: "NOT_GENERATED_YET"}
)

func TestTestimonials(t *testing.T) {
	err := MockedDTX.Mock()
	if err != nil {
		return
	}
	var test Testimonial
	Convey("Testing Create Testimonial", t, func() {
		test.Content = "Test Content"
		err = test.Create(MockedDTX)
		So(err, ShouldBeNil)
	})
	Convey("Update", t, func() {
		test.Content = "New Content"
		test.Active = true
		test.Approved = true
		err = test.Update(MockedDTX)
		So(err, ShouldBeNil)

	})

	Convey("Get testimonial", t, func() {
		err = test.Get(MockedDTX)
		So(err, ShouldBeNil)
	})
	Convey("GetAll - No paging", t, func() {
		ts, err := GetAllTestimonials(0, 1, false, MockedDTX)
		So(err, ShouldBeNil)
		So(len(ts), ShouldBeGreaterThan, 0)

	})

	Convey("GetAll - Paged", t, func() {
		ts, err := GetAllTestimonials(0, 1, false, MockedDTX)

		So(err, ShouldBeNil)
		So(len(ts), ShouldBeGreaterThan, 0)

	})

	Convey("GetAll - randomized", t, func() {
		ts, err := GetAllTestimonials(0, 1, true, MockedDTX)
		So(err, ShouldBeNil)
		So(len(ts), ShouldBeGreaterThan, 0)

	})
	Convey("Delete", t, func() {
		err = test.Delete()
		So(err, ShouldBeNil)

	})

	err = MockedDTX.DeMock()
	if err != nil {
		return
	}
}

func BenchmarkGetAllTestimonials(b *testing.B) {
	err := MockedDTX.Mock()
	if err != nil {
		return
	}
	for i := 0; i < b.N; i++ {
		GetAllTestimonials(0, 1, false, MockedDTX)
	}
	err = MockedDTX.DeMock()
	if err != nil {
		return
	}
}

func BenchmarkGetTestimonial(b *testing.B) {
	err := MockedDTX.Mock()
	if err != nil {
		return
	}
	test := setupDummyTestimonial()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		test.Create(MockedDTX)
		b.StartTimer()
		test.Get(MockedDTX)
		b.StopTimer()
		test.Delete()
	}
	err = MockedDTX.DeMock()
	if err != nil {
		return
	}
}

func BenchmarkCreateTestimonial(b *testing.B) {
	err := MockedDTX.Mock()
	if err != nil {
		return
	}
	test := setupDummyTestimonial()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		test.Create(MockedDTX)
		b.StopTimer()
		test.Delete()
	}
	err = MockedDTX.DeMock()
	if err != nil {
		return
	}
}

func BenchmarkUpdateTestimonial(b *testing.B) {
	err := MockedDTX.Mock()
	if err != nil {
		return
	}
	test := setupDummyTestimonial()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		test.Create(MockedDTX)
		b.StartTimer()
		test.Content = "This is a good test."
		test.Update(MockedDTX)
		b.StopTimer()
		test.Delete()
	}
	err = MockedDTX.DeMock()
	if err != nil {
		return
	}
}

func BenchmarkDeleteTestimonial(b *testing.B) {
	err := MockedDTX.Mock()
	if err != nil {
		return
	}
	test := setupDummyTestimonial()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		test.Create(MockedDTX)
		b.StartTimer()
		test.Delete()
	}
	err = MockedDTX.DeMock()
	if err != nil {
		return
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
		BrandID:   1,
	}
}
