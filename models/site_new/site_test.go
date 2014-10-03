package site_new

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
)

func getRandomMenuWithContents() (m Menu, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return m, err
	}
	defer db.Close()
	stmt, err := db.Prepare("select DISTINCT(msc.menuID) from Menu_SiteContent as msc  WHERE contentID > 0  ORDER BY RAND() LIMIT 1")
	if err != nil {
		return m, err
	}
	defer stmt.Close()
	err = stmt.QueryRow().Scan(&m.Id)
	if err != nil {
		return m, err
	}
	return m, err
}

func TestSite_New(t *testing.T) {
	Convey("Testing Gets", t, func() {
		Convey("Testing GetAllMenus", func() {

			ms, err := GetAllMenus()
			So(err, ShouldBeNil)
			//if thar' be menus, get random menu
			if len(ms) > 0 {
				i := rand.Intn(len(ms))
				menu := ms[i]

				//get by id
				err = menu.Get()
				if err != sql.ErrNoRows { //check for empty db
					So(menu, ShouldNotBeNil)
					So(menu.DisplayName, ShouldNotBeNil)

					So(err, ShouldBeNil)
				}
				//get by name
				err = menu.GetByName()
				So(err, ShouldBeNil)

				//get menu's contents
				err = menu.GetContents()
				t.Log(err)
				if err != sql.ErrNoRows {
					So(err, ShouldBeNil)
					// So(len(menu.Contents), ShouldBeGreaterThan, 0)
				}

			}
		})
		Convey("Test with random menu", func() {
			m, err := getRandomMenuWithContents()
			if err != sql.ErrNoRows {
				err = m.GetContents()
				So(err, ShouldBeNil)
				So(len(m.Contents), ShouldBeGreaterThan, 0)
			}

		})
		Convey("Testing GetAll Contents", func() {
			cs, err := GetAllContents()
			if err != sql.ErrNoRows {
				So(err, ShouldBeNil)
				So(len(cs), ShouldBeGreaterThan, 0)
				i := rand.Intn(len(cs))
				c := cs[i] //random content object
				Convey("Testing GetContent", func() {
					err = c.Get()
					So(err, ShouldBeNil)
					So(c, ShouldNotBeNil)
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
				})
				Convey("Testing GetContentBySlug", func() {
					err = c.GetbySlug()
					So(err, ShouldBeNil)
				})
			}
		})
		Convey("Testing ContentRevisions", func() {
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
	})
	Convey("Testing Content CRUD", t, func() {
		var c Content
		c.Type = "type"
		c.Title = "title"
		c.MetaTitle = "mTitle"
		c.MetaDescription = "mDesc"
		err := c.Create()
		So(err, ShouldBeNil)
		c.Get()
		So(c.Title, ShouldEqual, "title")
		c.Type = "type2"
		c.Title = "title2"
		c.MetaTitle = "mTitle2"
		c.MetaDescription = "mDesc2"
		err = c.Update()
		So(err, ShouldBeNil)
		c.Get()
		So(c.Title, ShouldEqual, "title2")
		err = c.Delete()
		So(err, ShouldBeNil)
	})

	Convey("Testing Revision CRUD", t, func() {
		var r ContentRevision
		//get rand content, for its id, tis a FK relation
		cs, err := GetAllContents()
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			i := rand.Intn(len(cs))
			c := cs[i] //random content object

			r.Text = "text"
			r.Active = true
			r.ContentId = c.Id

			err = r.Create()
			So(err, ShouldBeNil)
			r.Get()
			So(r.Text, ShouldEqual, "text")
			r.Text = "text2"
			r.Active = false
			err = r.Update()
			So(err, ShouldBeNil)
			r.Get()
			So(r.Text, ShouldEqual, "text2")
			err = r.Delete()
			So(err, ShouldBeNil)
		}
	})
	Convey("Testing Menu CRUD", t, func() {
		var m Menu
		m.Name = "name"
		m.ShowOnSitemap = true
		err := m.Create()
		So(err, ShouldBeNil)
		m.Get()
		So(m.Name, ShouldEqual, "name")
		m.Name = "name2"
		m.ShowOnSitemap = false
		err = m.Update()
		So(err, ShouldBeNil)
		m.Get()
		So(m.Name, ShouldEqual, "name2")
		err = m.Delete()
		So(err, ShouldBeNil)
	})
	Convey("Testing Menu-Content Joins", t, func() {
		var m Menu
		m.Name = "name"
		err := m.Create()
		So(err, ShouldBeNil)
		var c Content
		c.Title = "title"
		err = c.Create()
		So(err, ShouldBeNil)
		err = m.JoinToContent(c)
		So(err, ShouldBeNil)
		err = m.DeleteMenuContentJoin(c)
		So(err, ShouldBeNil)
	})

}
