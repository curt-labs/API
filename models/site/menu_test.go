package site

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
)

func getRandomMenuSiteContent() (msc Menu_SiteContent) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return msc
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT contentID FROM Menu_SiteContent ORDER BY RAND() LIMIT 1 WHERE contentID > 0")
	if err != nil {
		return msc
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRow().Scan(id)
	msc.Id = id
	return msc
}

func TestMenus(t *testing.T) {
	Convey("Testing Gets", t, func() {
		Convey("Testing GetPrimaryMenu()", func() {
			var m MenuWithContent
			var err error
			err = m.GetPrimaryMenu()
			if err != sql.ErrNoRows { //check for empty db
				So(m, ShouldNotBeNil)
				So(len(m.Contents), ShouldNotBeNil)
				So(m.Menu, ShouldNotBeNil)
				So(err, ShouldBeNil)
			}
		})
		Convey("Testing GetAllMenus", func() {
			var err error
			ms, err := GetAllMenus()
			So(err, ShouldBeNil)
			//if thar' be menus, get random menu
			if len(ms) > 0 {
				i := rand.Intn(len(ms))
				menu := ms[i]
				Convey("GetMenuItemsByMenuId", func() {
					menuItems, err := menu.GetMenuItemsByMenuId()
					So(err, ShouldBeNil)
					var mi MenuItems
					So(menuItems, ShouldHaveSameTypeAs, mi)
				})
				Convey("Test GetMenuWithContentByName", func() {
					var mwc MenuWithContent
					err = mwc.GetMenuWithContentByName(menu.Name)
					if err != sql.ErrNoRows {
						So(err, ShouldBeNil)
						So(mwc.DisplayName, ShouldNotBeNil)
					}
				})
			}

		})

		Convey("Testing GetFooterSitemap", func() {
			menuWithContents, err := GetFooterSitemap()
			if err != sql.ErrNoRows {
				So(len(menuWithContents), ShouldBeGreaterThan, 0)
				So(err, ShouldBeNil)
				//rand menu with contents
				i := rand.Intn(len(menuWithContents))
				mc := menuWithContents[i]
				Convey("GetMenuByContentId", func() {
					msc := getRandomMenuSiteContent()
					err = mc.GetMenuByContentId(msc.Id, false)
					t.Log(mc, msc.Id)
					if err != sql.ErrNoRows {
						So(err, ShouldBeNil)
						So(mc.Id, ShouldBeGreaterThan, 0)
					}

				})
				Convey("GetMenuIdByContentId", func() {
					menuId, err := GetMenuIdByContentId(getRandomMenuSiteContent().Id)
					So(err, ShouldBeNil)
					So(menuId, ShouldBeGreaterThanOrEqualTo, 0)
				})
			}
		})

		Convey("Testing GetMenuSitemap", func() {
			menuWithContents, err := GetMenuSitemap()
			if err != sql.ErrNoRows {
				So(len(menuWithContents), ShouldBeGreaterThan, 0)
				So(err, ShouldBeNil)
			}
		})

	})
}
