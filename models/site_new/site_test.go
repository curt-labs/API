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
				err = menu.Get()
				if err != sql.ErrNoRows { //check for empty db
					So(menu, ShouldNotBeNil)
					So(menu.DisplayName, ShouldNotBeNil)

					So(err, ShouldBeNil)
				}
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

}
