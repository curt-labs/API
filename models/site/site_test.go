package site

import (
	"database/sql"
	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
)

func TestSite_New(t *testing.T) {
	var c Content
	var r ContentRevision
	var m Menu
	var err error

	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}

	Convey("Testing Create Menus", t, func() {
		m.Name = "name"
		m.ShowOnSitemap = true
		m.WebsiteId = MockedDTX.WebsiteID
		err = m.Create()
		So(err, ShouldBeNil)
		//Update menu
		m.Name = "name2"
		m.ShowOnSitemap = false
		err = m.Update()
		So(err, ShouldBeNil)
		//Create content
		c.Type = "type"
		c.Title = "title"
		c.MetaTitle = "mTitle"
		c.MetaDescription = "mDesc"
		c.WebsiteId = MockedDTX.WebsiteID
		err = c.Create()
		So(err, ShouldBeNil)
		//update content
		c.Type = "type2"
		c.Title = "title2"
		c.MetaTitle = "mTitle2"
		c.MetaDescription = "mDesc2"
		c.Slug = "testSlug"
		err = c.Update()
		So(err, ShouldBeNil)

		//create revisions
		//Revisions
		r.Text = "text"
		r.Active = true
		r.ContentId = c.Id
		err = r.Create()
		So(err, ShouldBeNil)
		//update revisions
		r.Text = "text2"
		r.Active = false
		err = r.Update()
		So(err, ShouldBeNil)

		//get menus
		err = m.Get(MockedDTX)
		So(m.Name, ShouldEqual, "name2")
		So(err, ShouldBeNil)

		ms, err := GetAllMenus(MockedDTX)
		So(err, ShouldBeNil)
		So(len(ms), ShouldBeGreaterThan, 0)

		err = m.GetByName(MockedDTX)
		So(err, ShouldBeNil)

		err = m.GetContents()
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			// So(len(m.Contents), ShouldBeGreaterThan, 0)
		}
		//get revisions
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
		//get content

		err = c.Get(MockedDTX)
		So(c.Title, ShouldEqual, "title2")
		So(err, ShouldBeNil)

		//get revisions/content
		//get revisions
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
		//get contents
		cs, err := GetAllContents(MockedDTX)
		So(err, ShouldBeNil)
		So(len(cs), ShouldBeGreaterThanOrEqualTo, 0)
		//get contents by slug
		t.Log(c.Slug)
		err = c.GetBySlug(MockedDTX)
		So(err, ShouldBeNil)

		//menu-content join
		err = m.JoinToContent(c)
		So(err, ShouldBeNil)

		//check actual join
		err = m.GetContents()
		So(err, ShouldBeNil)

		So(len(m.Contents), ShouldBeGreaterThan, 0)
		err = m.DeleteMenuContentJoin(c)
		So(err, ShouldBeNil)

		//delete content
		err = c.Delete()
		So(err, ShouldBeNil)
		//delete revision
		err = r.Delete()
		So(err, ShouldBeNil)

		//delete menu
		err = m.Delete()
		So(err, ShouldBeNil)
	})
	_ = apicontextmock.DeMock(MockedDTX)
}

func TestWebsite(t *testing.T) {
	var w Website
	var err error

	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
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
		err = w.GetDetails(MockedDTX)
		So(err, ShouldBeNil)
	})
	Convey("Testing Delete Website", t, func() {
		err = w.Delete()
		So(err, ShouldBeNil)
	})

}
