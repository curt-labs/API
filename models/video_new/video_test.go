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

func TestVideo_New(t *testing.T) {
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
		Convey("Testing Get", func() {
			err = v.Get()
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
		Convey("Testing Delete", func() {
			err = v.Delete()
			So(err, ShouldBeNil)

		})
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

func TestVideoCatsnStuff_New(t *testing.T) {
	Convey("Testing Creates", t, func() {
		var err error
		var ch Channel
		var f CdnFile
		var v VideoType
		var cht ChannelType
		var ft CdnFileType

		ch.Description = "test"
		f.ObjectName = "test"
		v.Name = "test"
		cht.Name = "test"
		ft.Title = "test"

		err = ch.Create()
		So(err, ShouldBeNil)
		err = f.Create()
		So(err, ShouldBeNil)
		err = v.Create()
		So(err, ShouldBeNil)
		err = cht.Create()
		So(err, ShouldBeNil)
		err = ft.Create()
		So(err, ShouldBeNil)

		Convey("Testing Updates", func() {
			ch.Description = "new test"
			f.ObjectName = "new test"
			v.Name = "new test"
			cht.Name = "new test"
			ft.Title = "new test"

			err = ch.Update()
			So(err, ShouldBeNil)
			err = f.Update()
			So(err, ShouldBeNil)
			err = v.Update()
			So(err, ShouldBeNil)
			err = cht.Update()
			So(err, ShouldBeNil)
			err = ft.Update()
			So(err, ShouldBeNil)
		})
		Convey("Testing Gets", func() {
			err = ch.Get()
			So(err, ShouldBeNil)
			err = f.Get()
			So(err, ShouldBeNil)
			err = v.Get()
			So(err, ShouldBeNil)
			err = cht.Get()
			So(err, ShouldBeNil)
			err = ft.Get()
			So(err, ShouldBeNil)
		})
		Convey("Testing Deletes", func() {
			err = ch.Delete()
			So(err, ShouldBeNil)
			err = f.Delete()
			So(err, ShouldBeNil)
			err = v.Delete()
			So(err, ShouldBeNil)
			err = cht.Delete()
			So(err, ShouldBeNil)
			err = ft.Delete()
			So(err, ShouldBeNil)
		})

	})
}
