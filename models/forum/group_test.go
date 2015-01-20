package forum

import (
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGroups(t *testing.T) {
	var g Group
	var lastGroupID int
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}

	//setup an existing groupy
	groupy := Group{}
	groupy.Name = "Test Group"
	groupy.Description = "This is a test group."
	groupy.Add(MockedDTX)

	//run our tests
	Convey("Testing Gets", t, func() {
		Convey("Testing GetAll()", func() {
			groups, err := GetAllGroups(MockedDTX)

			So(groups, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("Testing Get()", func() {
			Convey("Group with ID of 0", func() {
				err = g.Get(MockedDTX)
				So(g.ID, ShouldEqual, 0)
				So(g.Name, ShouldEqual, "")
				So(g.Description, ShouldEqual, "")

				So(err, ShouldNotBeNil)
			})

			Convey("Group with non-zero ID", func() {
				g = Group{ID: groupy.ID}
				err = g.Get(MockedDTX)

				So(g.ID, ShouldNotEqual, 0)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Testing Add()", t, func() {
		Convey("Empty Group", func() {
			g = Group{}
			err = g.Add(MockedDTX)

			So(g.ID, ShouldEqual, 0)
			So(g.Name, ShouldEqual, "")
			So(g.Description, ShouldEqual, "")
			So(err, ShouldNotBeNil)
		})
		Convey("Valid Group", func() {
			g = Group{}
			g.Name = "Test Group"
			g.Description = "This is a test group."
			err = g.Add(MockedDTX)

			So(g.ID, ShouldBeGreaterThan, 0)
			So(g.Name, ShouldNotEqual, "")
			So(g.Description, ShouldNotEqual, "")
			So(err, ShouldBeNil)

			lastGroupID = g.ID
		})
	})

	Convey("Testing Update()", t, func() {
		Convey("Empty Group", func() {
			g = Group{}
			err = g.Update(MockedDTX)

			So(g.ID, ShouldEqual, 0)
			So(g.Name, ShouldEqual, "")
			So(g.Description, ShouldEqual, "")
			So(err, ShouldNotBeNil)
		})

		Convey("Missing Name", func() {
			g = Group{ID: lastGroupID}
			g.Name = ""
			err = g.Update(MockedDTX)

			So(g.ID, ShouldNotEqual, 0)
			So(g.Name, ShouldEqual, "")
			So(err, ShouldNotBeNil)
		})

		Convey("Last Added Group", func() {
			g = Group{ID: lastGroupID}
			g.Name = "Updated Test Group"
			g.Description = "This is a updated test group."
			err = g.Update(MockedDTX)

			So(g.ID, ShouldNotEqual, 0)
			So(g.Name, ShouldNotEqual, "")
			So(g.Description, ShouldNotEqual, "")
			So(err, ShouldBeNil)
		})
	})

	Convey("Testing Delete()", t, func() {
		Convey("Empty Group", func() {
			g = Group{}
			err = g.Delete(MockedDTX)

			So(err, ShouldNotBeNil)
		})

		Convey("Last Updated Group", func() {
			g = Group{ID: lastGroupID}
			err = g.Delete(MockedDTX)

			So(err, ShouldBeNil)
		})
	})

	//destroy the test group, we're done with it
	groupy.Delete(MockedDTX)

	_ = apicontextmock.DeMock(MockedDTX)
}
