package video

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/models/products"
	. "github.com/smartystreets/goconvey/convey"
	// "math/rand"
	"testing"
)

func getPartId() (pID int) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return 0
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT partID FROM VehiclePart ORDER BY RAND() LIMIT 1")
	if err != nil {
		return 0
	}
	defer stmt.Close()
	err = stmt.QueryRow().Scan(&pID)
	return pID
}

func TestSite_New(t *testing.T) {
	Convey("Testing Gets", t, func() {
		var v Video
		v.ID = 19
		err := v.Get()
		So(err, ShouldBeNil)

		chs, err := v.GetChannels()
		t.Log(err)
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(chs), ShouldBeGreaterThan, 0)
		}

		cdns, err := v.GetCdnFiles()
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(cdns), ShouldBeGreaterThan, 0)
		}
	})
	Convey("Testing Get-all", t, func() {
		//All Videos
		vs, err := GetAllVideos()
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(vs), ShouldBeGreaterThan, 0)
		}

		//Part Vids
		var p products.Part
		p.ID = getPartId()
		partVids, err := GetPartVideos(p)
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(partVids), ShouldBeGreaterThan, 0)
		}

	})
}
