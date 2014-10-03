package site

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
)

func getRandomLandingPage() (lp LandingPage) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return lp
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT id FROM LandingPage AS l WHERE startDate <= NOW() && endDate >= NOW() ORDER BY RAND() limit 1")
	if err != nil {
		return lp
	}
	defer stmt.Close()
	err = stmt.QueryRow().Scan(&lp.Id)
	if err != nil {
		return lp
	}
	return lp
}

func TestContents(t *testing.T) {
	Convey("Testing Gets", t, func() {
		Convey("Testing GetContentPageByName", func() {
			var cp ContentPage
			var err error
			ms, err := GetAllMenus()
			So(err, ShouldBeNil)
			//if thar' be menus, get random menu
			if len(ms) > 0 {
				i := rand.Intn(len(ms))
				menu := ms[i]
				err = cp.GetContentPageByName(menu.Id, true)
				if err != sql.ErrNoRows { //check for empty db
					So(cp, ShouldNotBeNil)
					So(cp.SiteContent, ShouldNotBeNil)
					So(cp.Revision, ShouldNotBeNil)
					So(cp.MenuWithContent, ShouldNotBeNil)
					So(err, ShouldBeNil)
				}
			}
		})
	})
	Convey("GetPrimaryContentPage", t, func() {
		var cp ContentPage
		err := cp.GetPrimaryContentPage()
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(cp, ShouldNotBeNil)
		}
	})
	Convey("GetSitemapCP", t, func() {
		cps, err := GetSitemapCP()
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(cps), ShouldBeGreaterThan, 0)
		}
	})
	Convey("GetLandingPage", t, func() {
		//get random LP id
		lp := getRandomLandingPage()
		err := lp.Get()
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(lp, ShouldNotBeNil)
		}
	})
}
