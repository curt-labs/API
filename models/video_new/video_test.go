package video

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/models/products"
	. "github.com/smartystreets/goconvey/convey"
	// "math/rand"
	"testing"
)

func getPartId() (id int) {
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
	err = stmt.QueryRow().Scan(&id)
	return id
}
func getCatId() (id int) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return 0
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT catID FROM Categories ORDER BY RAND() LIMIT 1")
	if err != nil {
		return 0
	}
	defer stmt.Close()
	err = stmt.QueryRow().Scan(&id)
	return id
}
func getChannelId() (id int) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return 0
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT ID FROM Channel ORDER BY RAND() LIMIT 1")
	if err != nil {
		return 0
	}
	defer stmt.Close()
	err = stmt.QueryRow().Scan(&id)
	return id
}
func getFileId() (id int) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return 0
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT ID FROM CdnFile ORDER BY RAND() LIMIT 1")
	if err != nil {
		return 0
	}
	defer stmt.Close()
	err = stmt.QueryRow().Scan(&id)
	return id
}

func TestSite_New(t *testing.T) {
	Convey("Testing Create", t, func() {
		var v Video
		var p products.Part
		var ch Channel
		var f CdnFile
		var c Category
		var err error
		p.ID = getPartId()
		ch.ID = getChannelId()
		f.ID = getFileId()
		c.ID = getChannelId()
		v.Title = "Test Video"
		v.Parts = append(v.Parts, p)
		v.Channels = append(v.Channels, ch)
		v.Files = append(v.Files, f)
		v.Categories = append(v.Categories, c)
		err = v.Create()
		So(err, ShouldBeNil)
		t.Log(v)
		Convey("Testing Update", func() {
			p.ID = getPartId()
			ch.ID = getChannelId()
			f.ID = getFileId()
			c.ID = getChannelId()
			v.Title = "Test Video"
			v.Parts = append(v.Parts, p)
			v.Channels = append(v.Channels, ch)
			v.Files = append(v.Files, f)
			v.Categories = append(v.Categories, c)
			err = v.Update()
			So(err, ShouldBeNil)

		})
		Convey("Testing Delete", func() {
			err = v.Delete()
			So(err, ShouldBeNil)

		})
	})

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
