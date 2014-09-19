package salesrep

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSalesReps(t *testing.T) {
	var err error
	var rep SalesRep
	var lastSalesRepID int

	//create a test rep
	testRep := SalesRep{
		Name: "Test SalesRep",
		Code: "9999",
	}
	testRep.Add()

	//run our tests
	Convey("Testing Gets", t, func() {
		Convey("Testing GetAll()", func() {
			reps, err := GetAllSalesReps()

			So(len(reps), ShouldBeGreaterThanOrEqualTo, 0)
			So(err, ShouldBeNil)
		})

		Convey("Testing Get()", func() {
			Convey("SalesRep with ID of 0", func() {
				err = rep.Get()

				So(rep.ID, ShouldEqual, 0)
				So(rep.Name, ShouldEqual, "")
				So(rep.Code, ShouldEqual, "")
			})

			Convey("SalesRep with non-zero ID", func() {
				rep = SalesRep{ID: testRep.ID}
				err = rep.Get()

				So(rep.ID, ShouldNotEqual, 0)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Testing Add()", t, func() {
		Convey("Empty SalesRep", func() {
			rep = SalesRep{}
			err = rep.Add()

			So(rep.ID, ShouldEqual, 0)
			So(rep.Name, ShouldEqual, "")
			So(rep.Code, ShouldEqual, "")
			So(err, ShouldNotBeNil)
		})
		Convey("Missing name", func() {
			rep = SalesRep{
				Code: "9999",
			}

			err = rep.Add()

			So(rep.ID, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})
		Convey("Valid SalesRep", func() {
			rep = SalesRep{
				Name: "Test SalesRep",
				Code: "9999",
			}

			err = rep.Add()

			So(rep.ID, ShouldBeGreaterThan, 0)
			So(rep.Name, ShouldNotEqual, "")
			So(rep.Code, ShouldNotEqual, "")
			So(err, ShouldBeNil)

			lastSalesRepID = rep.ID
		})
	})

	Convey("Testing Update()", t, func() {
		Convey("Empty SalesRep", func() {
			rep = SalesRep{}
			err = rep.Update()

			So(rep.ID, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})

		Convey("Missing Name", func() {
			rep = SalesRep{
				Code: "9999",
			}
			err = rep.Update()

			So(rep.ID, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})

		Convey("Last Added SalesRep", func() {
			rep = SalesRep{ID: lastSalesRepID}
			rep.Name = "Test SalesRep"
			rep.Code = "9999"

			err = rep.Update()

			So(rep.ID, ShouldNotEqual, 0)
			So(rep.Name, ShouldNotEqual, "")
			So(rep.Code, ShouldNotEqual, "")
			So(err, ShouldBeNil)
		})
	})

	Convey("Testing Delete()", t, func() {
		Convey("Empty SalesRep", func() {
			rep = SalesRep{}
			err = rep.Delete()

			So(err, ShouldNotBeNil)
		})

		Convey("Last Updated Group", func() {
			rep = SalesRep{ID: lastSalesRepID}
			err = rep.Delete()

			So(err, ShouldBeNil)
		})
	})

	//delete our test record
	testRep.Delete()
}
