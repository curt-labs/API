package faq_model

import (
	// "database/sql"
	// "github.com/curt-labs/GoAPI/helpers/database"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"strconv"
	"testing"
)

func TestGetFaqs(t *testing.T) {
	Convey("Testing Gets", t, func() {
		Convey("Testing GetAll()", func() {
			var fs Faqs
			fs, err := GetAll()
			So(fs, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(len(fs), ShouldNotBeNil)
		})
		Convey("Gets a faqs and a random faq, page, resultsperpage", func() {
			fs, err := GetAll()
			So(err, ShouldBeNil)
			if len(fs) > 0 {
				x := rand.Intn(len(fs))
				f := fs[x]
				totalResults := len(fs)
				result := strconv.Itoa(totalResults / 3)
				page := "3"

				Convey("Testing Get()", func() {
					f.Get()
					So(f, ShouldNotBeNil)
					So(f.Question, ShouldHaveSameTypeAs, "str")
					So(f.Answer, ShouldHaveSameTypeAs, "str")
					Convey("Testing Bad Get()", func() {
						getFaq = "Bad Query Stmt"
						err = f.Get()
						So(err, ShouldNotBeNil)
					})
				})

				Convey("Testing GetQuestions()", func() {

					qs, err := GetQuestions(page, result)
					So(qs, ShouldNotBeNil)
					So(err, ShouldBeNil)
					So(qs.Pagination.Page, ShouldEqual, 3)
					So(qs.Pagination.ReturnedCount, ShouldNotBeNil)
					So(qs.Pagination.PerPage, ShouldEqual, len(qs.Objects))
					Convey("Testing Bad Stmt()", func() {
						getQuestions = "Bad Query Stmt"
						qs, err = GetQuestions(page, result)
						So(err, ShouldNotBeNil)
					})
				})
				Convey("Testing GetAnswers()", func() {
					as, err := GetAnswers(page, result)
					So(as, ShouldNotBeNil)
					So(err, ShouldBeNil)
					So(as.Pagination.Page, ShouldEqual, 3)
					So(as.Pagination.ReturnedCount, ShouldNotBeNil)
					So(as.Pagination.PerPage, ShouldNotBeNil)
					So(as.Pagination.PerPage, ShouldEqual, len(as.Objects))
					Convey("Testing Bad Stmt()", func() {
						getAnswers = "Bad Query Stmt"
						as, err = GetAnswers(page, result)
						So(err, ShouldNotBeNil)
					})
				})
				Convey("Testing Search()", func() {
					as, err := Search(f.Question, "", "1", "0")
					So(as, ShouldNotBeNil)
					So(err, ShouldBeNil)
					So(as.Pagination.Page, ShouldEqual, 1)
					So(as.Pagination.ReturnedCount, ShouldNotBeNil)
					So(as.Pagination.PerPage, ShouldNotBeNil)
					So(as.Pagination.PerPage, ShouldEqual, len(as.Objects))
					Convey("Testing Bad Stmt()", func() {
						search = "Bad Query Stmt"
						as, err = Search(f.Question, "", "1", "0")
						So(err, ShouldNotBeNil)
					})
				})
			}
		})

	})
	Convey("Testing C_UD", t, func() {
		Convey("Testing Create()", func() {
			var f Faq
			f.Question = "testQuestion"
			f.Answer = "testAnswer"
			f.Create()
			So(f, ShouldNotBeNil)
			So(f.Question, ShouldEqual, "testQuestion")
			So(f.Answer, ShouldEqual, "testAnswer")

			Convey("Testing Update()", func() {
				f.Question = "testQuestion222"
				f.Answer = "testAnswer222"
				f.Update()
				So(f, ShouldNotBeNil)
				So(f.Question, ShouldEqual, "testQuestion222")
				So(f.Answer, ShouldEqual, "testAnswer222")

				Convey("Testing Delete()", func() {
					f.Get()
					err := f.Delete()
					So(err, ShouldBeNil)

				})
			})
		})

	})

}
