package forum

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGroups(t *testing.T) {
	var g Group
	var lastGroupID int
	var err error

	Convey("Testing Gets", t, func() {
		Convey("Testing GetAll()", func() {
			groups, err := GetAllGroups()

			So(groups, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("Testing Get()", func() {
			Convey("Group with ID of 0", func() {
				err = g.Get()
				So(g.ID, ShouldEqual, 0)
				So(g.Name, ShouldEqual, "")
				So(g.Description, ShouldEqual, "")

				So(err, ShouldNotBeNil)
			})

			Convey("Group with ID of 1", func() {
				g.ID = 1
				err = g.Get()

				So(g.ID, ShouldNotEqual, 0)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Testing Add()", t, func() {
		Convey("Empty Group", func() {
			g = Group{}
			err = g.Add()

			So(g.ID, ShouldEqual, 0)
			So(g.Name, ShouldEqual, "")
			So(g.Description, ShouldEqual, "")
			So(err, ShouldNotBeNil)
		})
		Convey("Valid Group", func() {
			g = Group{}
			g.Name = "Test Group"
			g.Description = "This is a test group."
			err = g.Add()

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
			err = g.Update()

			So(g.ID, ShouldEqual, 0)
			So(g.Name, ShouldEqual, "")
			So(g.Description, ShouldEqual, "")
			So(err, ShouldNotBeNil)
		})

		Convey("Missing Name", func() {
			g = Group{ID: lastGroupID}
			g.Name = ""
			err = g.Update()

			So(g.ID, ShouldNotEqual, 0)
			So(g.Name, ShouldEqual, "")
			So(err, ShouldNotBeNil)
		})

		Convey("Last Added Group", func() {
			g = Group{ID: lastGroupID}
			g.Name = "Updated Test Group"
			g.Description = "This is a updated test group."
			err = g.Update()

			So(g.ID, ShouldNotEqual, 0)
			So(g.Name, ShouldNotEqual, "")
			So(g.Description, ShouldNotEqual, "")
			So(err, ShouldBeNil)
		})
	})

	Convey("Testing Delete()", t, func() {
		Convey("Empty Group", func() {
			g = Group{}
			err = g.Delete()

			So(err, ShouldNotBeNil)
		})

		Convey("Last Updated Group", func() {
			g = Group{ID: lastGroupID}
			err = g.Delete()

			So(err, ShouldBeNil)
		})
	})
}
