package salesrep

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSalesReps(t *testing.T) {
	var err error
	var rep SalesRep

	Convey("Testing Add()", t, func() {
		rep.Name = "Name"
		err = rep.Add()
		So(err, ShouldBeNil)
	})

	Convey("Testing Update()", t, func() {
		rep.Name = "testname"
		err = rep.Update()
		So(err, ShouldBeNil)
	})

	Convey("Testing Gets()", t, func() {
		err = rep.Get()
		So(rep.ID, ShouldNotEqual, 0)
		So(rep.Name, ShouldEqual, "testname")

		reps, err := GetAllSalesReps()
		So(len(reps), ShouldBeGreaterThan, 0)
		So(err, ShouldBeNil)
	})

	Convey("Testing Delete()", t, func() {
		Convey("Empty SalesRep", func() {
			err = rep.Delete()
			So(err, ShouldBeNil)
		})
	})
}

func BenchmarkGetAllSalesReps(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetAllSalesReps()
	}
}

func BenchmarkGetSalesRep(b *testing.B) {
	s := SalesRep{Name: "TESTER", Code: "9999"}
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s.Add()
		b.StartTimer()
		s.Get()
		b.StopTimer()
		s.Delete()
	}
}

func BenchmarkAddSalesRep(b *testing.B) {
	s := SalesRep{Name: "TESTER", Code: "9999"}
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		s.Add()
		b.StopTimer()
		s.Delete()
	}
}

func BenchmarkUpdateSalesRep(b *testing.B) {
	s := SalesRep{Name: "TESTER", Code: "9999"}
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s.Add()
		b.StartTimer()
		s.Code = "99999"
		s.Update()
		b.StopTimer()
		s.Delete()
	}
}

func BenchmarkDeleteSalesRep(b *testing.B) {
	s := SalesRep{Name: "TESTER", Code: "9999"}
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s.Add()
		b.StartTimer()
		s.Delete()
	}
}
