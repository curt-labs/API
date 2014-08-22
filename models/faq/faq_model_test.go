package faq_model

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetFaqs(t *testing.T) {
	Convey("Testing Gets", t, func() {
		Convey("Testing Get()", func() {
			var f Faq
			f.ID = 1
			f.Get()
			So(f, ShouldNotBeNil)
			So(f.Question, ShouldEqual, "Can you use weight distribution on trailers with surge brakes?")
			So(f.Answer, ShouldEqual, "Yes. The coupler still has solid contact to the ball with a free range of motion.")
		})
		Convey("Testing GetAll()", func() {
			var fs Faqs
			fs, err := GetAll()
			So(fs, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(len(fs), ShouldNotBeNil)
		})
		Convey("Testing GetQuestions()", func() {
			qs, err := GetQuestions("1", "3")
			So(qs, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(qs.Pagination.Page, ShouldEqual, 1)
			So(qs.Pagination.ReturnedCount, ShouldNotBeNil)
			So(qs.Pagination.PerPage, ShouldEqual, len(qs.Objects))
		})
		Convey("Testing GetAnswers()", func() {
			as, err := GetAnswers("1", "0")
			So(as, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(as.Pagination.Page, ShouldEqual, 1)
			So(as.Pagination.ReturnedCount, ShouldNotBeNil)
			So(as.Pagination.PerPage, ShouldNotBeNil)
			So(as.Pagination.PerPage, ShouldEqual, len(as.Objects))
		})
		Convey("Testing Search()", func() {
			as, err := Search("hitch", "", "1", "0")
			So(as, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(as.Pagination.Page, ShouldEqual, 1)
			So(as.Pagination.ReturnedCount, ShouldNotBeNil)
			So(as.Pagination.PerPage, ShouldNotBeNil)
			So(as.Pagination.PerPage, ShouldEqual, len(as.Objects))
		})

	})
	Convey("Testing CUD", t, func() {
		Convey("Testing Create()", func() {
			var f Faq
			f.Question = "testQuestion"
			f.Answer = "testAnswer"
			f.Create()
			So(f, ShouldNotBeNil)
			So(f.Question, ShouldEqual, "testQuestion")
			So(f.Answer, ShouldEqual, "testAnswer")
		})
		Convey("Testing Update()", func() {
			var f Faq
			f.ID = 15
			f.Question = "testQuestion222"
			f.Answer = "testAnswer222"
			f.Update()
			So(f, ShouldNotBeNil)
			So(f.Question, ShouldEqual, "testQuestion222")
			So(f.Answer, ShouldEqual, "testAnswer222")
		})
		Convey("Testing Delete()", func() {
			var f Faq
			f.ID = 15
			f.Delete()
			f.Get()
			So(f.Question, ShouldBeBlank)
			So(f.Answer, ShouldBeBlank)
		})

	})

}
