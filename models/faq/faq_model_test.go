package faq_model

import (
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetFaqs(t *testing.T) {
	var f Faq
	var err error

	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}

	f.Question = "testQuestion"
	f.Answer = "testAnswer"

	Convey("Testing Create", t, func() {
		err = f.Create(MockedDTX)
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)
		So(f.Question, ShouldEqual, "testQuestion")
		So(f.Answer, ShouldEqual, "testAnswer")
	})
	Convey("Testing Update", t, func() {
		f.Question = "testQuestion222"
		f.Answer = "testAnswer222"
		err = f.Update(MockedDTX)
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)
		So(f.Question, ShouldEqual, "testQuestion222")
		So(f.Answer, ShouldEqual, "testAnswer222")
	})
	Convey("Testing Get", t, func() {
		err = f.Get(MockedDTX)
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)
		So(f.Question, ShouldHaveSameTypeAs, "str")
		So(f.Answer, ShouldHaveSameTypeAs, "str")

		var fs Faqs
		fs, err = GetAll(MockedDTX)
		So(fs, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(len(fs), ShouldNotBeNil)

		as, err := Search(MockedDTX, f.Question, "", "1", "0")
		So(as, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(as.Pagination.Page, ShouldEqual, 1)
		So(as.Pagination.ReturnedCount, ShouldNotBeNil)
		So(as.Pagination.PerPage, ShouldNotBeNil)
		So(as.Pagination.PerPage, ShouldEqual, len(as.Objects))
	})
	Convey("Testing Delete", t, func() {
		err = f.Delete()
		So(err, ShouldBeNil)

	})

	_ = apicontextmock.DeMock(MockedDTX)
}

func BenchmarkGetAllFaq(b *testing.B) {
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetAll(MockedDTX)
	}
	b.StopTimer()
	_ = apicontextmock.DeMock(MockedDTX)
}

func BenchmarkGetFaq(b *testing.B) {
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}
	f := setupDummyFaq()
	f.Create(MockedDTX)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Get(MockedDTX)
	}
	b.StopTimer()
	f.Delete()
	_ = apicontextmock.DeMock(MockedDTX)
}

func BenchmarkCreateFaq(b *testing.B) {
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}
	f := setupDummyFaq()
	for i := 0; i < b.N; i++ {
		f.Create(MockedDTX)
		b.StopTimer()
		f.Delete()
		b.StartTimer()
	}
	_ = apicontextmock.DeMock(MockedDTX)
}

func BenchmarkUpdateFaq(b *testing.B) {
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}
	f := setupDummyFaq()
	f.Create(MockedDTX)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Question = "Testing for real?"
		f.Answer = "You betcha."
		f.Update(MockedDTX)
	}
	b.StopTimer()
	f.Delete()
	_ = apicontextmock.DeMock(MockedDTX)
}

func BenchmarkDeleteFaq(b *testing.B) {
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}
	f := setupDummyFaq()
	f.Create(MockedDTX)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Delete()
	}
	b.StopTimer()
	f.Delete()
	_ = apicontextmock.DeMock(MockedDTX)
}

func setupDummyFaq() *Faq {
	return &Faq{
		Question: "Testing 123?",
		Answer:   "Yes...this is a test.",
	}
}
