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
