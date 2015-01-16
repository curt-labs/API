package lifestyle

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetLifestyles(t *testing.T) {
	var l Lifestyle
	var err error
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}

	Convey("Testing CRUD", t, func() {

		l.Name = "testName"
		l.LongDesc = "Long description"
		err := l.Create(MockedDTX)
		So(err, ShouldBeNil)
		err = l.Get(MockedDTX)
		So(err, ShouldBeNil)
		So(l, ShouldNotBeNil)
		So(l.Name, ShouldEqual, "testName")
		So(l.LongDesc, ShouldEqual, "Long description")

		//Update
		l.Name = "newName"
		l.Image = "image"
		l.ShortDesc = "Desc"
		err = l.Update(MockedDTX)
		So(err, ShouldBeNil)
		err = l.Get(MockedDTX)

		So(err, ShouldBeNil)
		So(l, ShouldNotBeNil)
		So(l.Name, ShouldEqual, "newName")
		So(l.Image, ShouldEqual, "image")
		So(l.ShortDesc, ShouldEqual, "Desc")

		//Gets
		ls, err := GetAll(MockedDTX)
		So(err, ShouldBeNil)
		So(ls, ShouldHaveSameTypeAs, Lifestyles{})

		err = l.Get(MockedDTX)
		So(err, ShouldBeNil)
		So(l, ShouldNotBeNil)
		So(l.Name, ShouldHaveSameTypeAs, "str")
		So(l.ShortDesc, ShouldHaveSameTypeAs, "str")

		//Delete
		err = l.Delete()
		So(err, ShouldBeNil)

	})
}

func BenchmarkGetAllLifestyles(b *testing.B) {
	var err error
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	for i := 0; i < b.N; i++ {
		GetAll(MockedDTX)
	}
}

func BenchmarkGetLifestyle(b *testing.B) {
	var err error
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	ls := setupDummyLifestyle()
	ls.Create(MockedDTX)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ls.Get(MockedDTX)
	}
	b.StopTimer()
	ls.Delete()
}

func BenchmarkCreateLifestyle(b *testing.B) {
	ls := setupDummyLifestyle()
	var err error
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	for i := 0; i < b.N; i++ {
		ls.Create(MockedDTX)
		b.StopTimer()
		ls.Delete()
		b.StartTimer()
	}
}

func BenchmarkUpdateLifestyle(b *testing.B) {
	var err error
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	ls := setupDummyLifestyle()
	ls.Create(MockedDTX)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ls.ShortDesc = "TEST"
		ls.LongDesc = "THIS IS A TEST"
		ls.Update(MockedDTX)
	}
	b.StopTimer()
	ls.Delete()
}

func BenchmarkDeleteLifestyle(b *testing.B) {
	var err error
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	ls := setupDummyLifestyle()
	ls.Create(MockedDTX)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ls.Delete()
	}
	b.StopTimer()
	ls.Delete()
}

func setupDummyLifestyle() *Lifestyle {
	return &Lifestyle{
		Name:      "TESTER",
		ShortDesc: "TESTER",
		LongDesc:  "TESTER TESTER",
	}
}
