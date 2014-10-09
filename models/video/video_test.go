package video

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/models/products"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

//Data Modelling Helper Functions
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

//Test Videos
func TestVideo_New(t *testing.T) {
	Convey("Testing Videos", t, func() {
		var v Video
		var p products.Part
		var ch Channel
		var f CdnFile
		var c Category
		var err error
		//Create
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

		//Update
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

		//Get
		err = v.Get()
		So(err, ShouldBeNil)
		t.Log("ID", v.ID)
		//Get Chans

		chans, err := v.GetChannels()
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(chans), ShouldBeGreaterThan, 0)
		}

		//Get CDNs
		cdns, err := v.GetCdnFiles()
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(cdns), ShouldBeGreaterThan, 0)
		}

		//Delete
		err = v.Delete()
		So(err, ShouldBeNil)

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

//Test Video-related items
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

		//Updates
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
		//Gets
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
		//Deletes
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

	Convey("Testing GetAll", t, func() {
		//Channels
		chs, err := GetAllChannels()
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(chs), ShouldBeGreaterThan, 0)
		}
		//CdnFiles
		cdns, err := GetAllCdnFiles()
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(cdns), ShouldBeGreaterThan, 0)
		}
		//VideoTypes
		vts, err := GetAllVideoTypes()
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(vts), ShouldBeGreaterThan, 0)
		}
		//ChannelTypes
		cts, err := GetAllChannelTypes()
		t.Log(cts)
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(cts), ShouldBeGreaterThan, 0)
		}

		//CdnFileTypes
		cds, err := GetAllCdnFileTypes()
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(cds), ShouldBeGreaterThan, 0)
		}
	})
}
