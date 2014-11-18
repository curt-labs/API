package lifestyle

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetLifestyles(t *testing.T) {
	var l Lifestyle

	Convey("Testing CRUD", t, func() {

		l.Name = "testName"
		l.LongDesc = "Long description"
		err := l.Create()
		So(err, ShouldBeNil)
		err = l.Get()
		So(err, ShouldBeNil)
		So(l, ShouldNotBeNil)
		So(l.Name, ShouldEqual, "testName")
		So(l.LongDesc, ShouldEqual, "Long description")

		//Update
		l.Name = "newName"
		l.Image = "image"
		l.ShortDesc = "Desc"
		err = l.Update()
		So(err, ShouldBeNil)
		err = l.Get()

		So(err, ShouldBeNil)
		So(l, ShouldNotBeNil)
		So(l.Name, ShouldEqual, "newName")
		So(l.Image, ShouldEqual, "image")
		So(l.ShortDesc, ShouldEqual, "Desc")

		//Gets
		ls, err := GetAll()
		So(err, ShouldBeNil)
		So(ls, ShouldHaveSameTypeAs, Lifestyles{})

		err = l.Get()
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
	for i := 0; i < b.N; i++ {
		GetAll()
	}
}

func BenchmarkGetLifestyle(b *testing.B) {
	ls := setupDummyLifestyle()
	ls.Create()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ls.Get()
	}
	b.StopTimer()
	ls.Delete()
}

func BenchmarkCreateLifestyle(b *testing.B) {
	ls := setupDummyLifestyle()
	for i := 0; i < b.N; i++ {
		ls.Create()
		b.StopTimer()
		ls.Delete()
		b.StartTimer()
	}
}

func BenchmarkUpdateLifestyle(b *testing.B) {
	ls := setupDummyLifestyle()
	ls.Create()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ls.ShortDesc = "TEST"
		ls.LongDesc = "THIS IS A TEST"
		ls.Update()
	}
	b.StopTimer()
	ls.Delete()
}

func BenchmarkDeleteLifestyle(b *testing.B) {
	ls := setupDummyLifestyle()
	ls.Create()
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
