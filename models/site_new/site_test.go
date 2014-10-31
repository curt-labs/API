package site_new

import (
	"database/sql"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
)

func TestSite_New(t *testing.T) {
	var c Content
	var r ContentRevision
	var m Menu
	var err error

	Convey("Testing Create Menus", t, func() {
		m.Name = "name"
		m.ShowOnSitemap = true
		err := m.Create()
		So(err, ShouldBeNil)
	})
	Convey("Testing Update Menus", t, func() {

		m.Name = "name2"
		m.ShowOnSitemap = false
		err = m.Update()
		So(err, ShouldBeNil)
	})

	Convey("Testing Create Content", t, func() {
		c.Type = "type"
		c.Title = "title"
		c.MetaTitle = "mTitle"
		c.MetaDescription = "mDesc"
		err := c.Create()
		So(err, ShouldBeNil)
	})
	Convey("Testing Update Content", t, func() {
		c.Type = "type2"
		c.Title = "title2"
		c.MetaTitle = "mTitle2"
		c.MetaDescription = "mDesc2"
		c.Slug = "testSlug"
		err = c.Update()
		So(err, ShouldBeNil)

	})
	Convey("Testing Create Revision", t, func() {

		//Revisions
		r.Text = "text"
		r.Active = true
		r.ContentId = c.Id
		err = r.Create()
		So(err, ShouldBeNil)
	})
	Convey("Testing Update Revision", t, func() {
		r.Text = "text2"
		r.Active = false
		err = r.Update()
		So(err, ShouldBeNil)

	})
	Convey("Testing Get Menus", t, func() {
		err = m.Get()
		So(m.Name, ShouldEqual, "name2")
		So(err, ShouldBeNil)

		ms, err := GetAllMenus()
		So(err, ShouldBeNil)
		So(len(ms), ShouldBeGreaterThan, 0)

		err = m.GetByName()
		So(err, ShouldBeNil)

		err = m.GetContents()
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			// So(len(m.Contents), ShouldBeGreaterThan, 0)
		}
	})
	Convey("Testing Get Revisions", t, func() {
		err = r.Get()
		So(r.Text, ShouldEqual, "text2")
		So(err, ShouldBeNil)

		cr, err := GetAllContentRevisions()
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(cr), ShouldBeGreaterThan, 0)
			i := rand.Intn(len(cr))
			r := cr[i] //random revision
			err = r.Get()
			So(err, ShouldBeNil)
			So(r, ShouldNotBeNil)
		}
	})
	Convey("Testing Get Content", t, func() {

		err = c.Get()
		So(c.Title, ShouldEqual, "title2")
		So(err, ShouldBeNil)

		Convey("Testing Get Revisions", func() {
			err = c.GetContentRevisions()
			if err != sql.ErrNoRows {
				So(err, ShouldBeNil)
				So(len(c.ContentRevisions), ShouldBeGreaterThan, 0)
			}
			err = c.GetLatestRevision()
			if err != sql.ErrNoRows {
				So(err, ShouldBeNil)
				So(c.ContentRevisions, ShouldNotBeNil)
			}
		})
		Convey("Testing GetAllContents", func() {
			cs, err := GetAllContents()
			So(err, ShouldBeNil)
			So(len(cs), ShouldBeGreaterThanOrEqualTo, 0)
		})
		Convey("Testing GetContentBySlug", func() {
			err = c.GetBySlug()
			So(err, ShouldBeNil)
		})
	})

	Convey("Testing Menu-Content Joins", t, func() {
		err = m.JoinToContent(c)
		So(err, ShouldBeNil)

		//check actual join
		err = m.GetContents()
		So(err, ShouldBeNil)

		So(len(m.Contents), ShouldBeGreaterThan, 0)
		err = m.DeleteMenuContentJoin(c)
		So(err, ShouldBeNil)

	})

	Convey("Testing Delete Content", t, func() {
		err = c.Delete()
		So(err, ShouldBeNil)
	})

	Convey("Testing Delete Revision", t, func() {
		err = r.Delete()
		So(err, ShouldBeNil)

	})
	Convey("Testing Delete Menu", t, func() {
		err = m.Delete()
		So(err, ShouldBeNil)
	})
}

func TestWebsite(t *testing.T) {
	var w Website
	var err error
	Convey("Testing Create Website", t, func() {
		w.Url = "www.example.com"
		err = w.Create()
		So(err, ShouldBeNil)
	})
	Convey("Testing Update Website", t, func() {

		w.Description = "Desc"
		err = w.Update()
		So(err, ShouldBeNil)
	})
	Convey("Testing Get Website", t, func() {
		err = w.Get()
		So(err, ShouldBeNil)
	})
	Convey("Testing GetAll Websites", t, func() {
		ws, err := GetAllWebsites()
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(ws), ShouldBeGreaterThan, 0)
		}
	})
	Convey("Testing GetSiteDetails", t, func() {
		err = w.GetDetails()
		So(err, ShouldBeNil)
	})
	Convey("Testing Delete Website", t, func() {
		err = w.Delete()
		So(err, ShouldBeNil)
	})

}
