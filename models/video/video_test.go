package video

import (
	// "database/sql"
	// "github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/models/products"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestVideo_New(t *testing.T) {
	var v Video
	var ch Channel
	var cdn CdnFile
	var cdnft CdnFileType
	var vt VideoType
	var ct ChannelType
	var p products.Part
	var cat products.Category
	var err error

	//Creates
	Convey("Testing Creates", t, func() {
		ct.Name = "test Channel Type"
		err = ct.Create()
		So(err, ShouldBeNil)
	})
	Convey("Testing Creates", t, func() {
		ch.Title = "test title"
		ch.Type = ct
		err = ch.Create()
		So(err, ShouldBeNil)
	})
	Convey("Testing Creates", t, func() {
		cdnft.Title = "test cdntype"
		err = cdnft.Create()
		So(err, ShouldBeNil)
	})
	Convey("Testing Creates", t, func() {
		cdn.ObjectName = "test cdn"
		cdn.Type = cdnft
		err = cdn.Create()
		So(err, ShouldBeNil)
	})
	Convey("Testing Creates", t, func() {
		vt.Name = "test videoType"
		err = vt.Create()
		So(err, ShouldBeNil)
	})
	Convey("Testing Creates", t, func() {
		cat.Title = "test cat title"
		err = cat.Create()
		So(err, ShouldBeNil)
	})
	Convey("Testing Creates", t, func() {
		v.Title = "test vid"
		p.ID = 11000 //force part

		v.VideoType = vt
		v.Categories = append(v.Categories, cat)
		v.Channels = append(v.Channels, ch)
		v.Files = append(v.Files, cdn)
		err = v.Create()
		So(err, ShouldBeNil)
	})

	//Updates
	Convey("Testing Update", t, func() {
		ct.Name = "test Channel Type 2"
		err = ct.Update()
		So(err, ShouldBeNil)
	})
	Convey("Testing Update", t, func() {
		ch.Title = "test title 2"
		err = ch.Update()
		So(err, ShouldBeNil)
	})
	Convey("Testing Update", t, func() {
		cdnft.Title = "test cdntype 2"
		err = cdnft.Update()
		So(err, ShouldBeNil)
	})
	Convey("Testing Update", t, func() {
		vt.Name = "test videoType 2"
		err = vt.Update()
		So(err, ShouldBeNil)
	})
	Convey("Testing Update", t, func() {
		cdn.ObjectName = "test cdn 2"
		err = cdn.Update()
		So(err, ShouldBeNil)
	})
	Convey("Testing Update", t, func() {
		v.Title = "test vid 2"

		p.ID = 110001 //force part
		v.Parts = append(v.Parts, p)

		err = v.Update()
		So(err, ShouldBeNil)
	})

	//Get
	Convey("Testing Get", t, func() {
		err = v.Get()
		So(err, ShouldBeNil)
	})
	Convey("Testing Get Details", t, func() {
		err = v.GetVideoDetails()
		So(err, ShouldBeNil)
	})
	Convey("Testing GetAllVideos", t, func() {
		vs, err := GetAllVideos()
		So(err, ShouldBeNil)
		So(len(vs), ShouldBeGreaterThan, 0)
	})
	Convey("Testing GetAllVideos", t, func() {
		vs, err := GetPartVideos(p)
		So(err, ShouldBeNil)
		So(len(vs), ShouldBeGreaterThan, 0)
	})
	Convey("Testing GetChannels", t, func() {
		chs, err := v.GetChannels()
		So(err, ShouldBeNil)
		So(len(chs), ShouldBeGreaterThan, 0)
	})
	Convey("Testing GetCdnFiles", t, func() {
		cdns, err := v.GetCdnFiles()
		So(err, ShouldBeNil)
		So(len(cdns), ShouldBeGreaterThan, 0)
	})
	//Gets
	Convey("Testing Get Channel", t, func() {
		err = ch.Get()
		So(err, ShouldBeNil)
	})
	Convey("Testing Get CdnFile", t, func() {
		err = cdn.Get()
		So(err, ShouldBeNil)
	})
	Convey("Testing Get CdnFileType", t, func() {
		err = cdnft.Get()
		So(err, ShouldBeNil)
	})
	Convey("Testing Get VideoType", t, func() {
		err = vt.Get()
		So(err, ShouldBeNil)
	})
	Convey("Testing Get ChannelType", t, func() {
		err = ct.Get()
		So(err, ShouldBeNil)
	})
	//Get Alls
	Convey("Testing Get All Channel", t, func() {
		chs, err := GetAllChannels()
		So(err, ShouldBeNil)
		So(len(chs), ShouldBeGreaterThan, 0)
	})
	Convey("Testing Get All CdnFile", t, func() {
		cdns, err := GetAllCdnFiles()
		So(err, ShouldBeNil)
		So(len(cdns), ShouldBeGreaterThan, 0)
	})
	Convey("Testing Get All CdnFileType", t, func() {
		cdnfts, err := GetAllCdnFileTypes()
		So(err, ShouldBeNil)
		So(len(cdnfts), ShouldBeGreaterThan, 0)
	})
	Convey("Testing Get All VideoType", t, func() {
		vts, err := GetAllVideoTypes()
		So(err, ShouldBeNil)
		So(len(vts), ShouldBeGreaterThan, 0)
	})
	Convey("Testing Get All ChannelType", t, func() {
		cts, err := GetAllChannelTypes()
		So(err, ShouldBeNil)
		So(len(cts), ShouldBeGreaterThan, 0)
	})

	//Deletes
	Convey("Testing Delete", t, func() {
		err = ct.Delete()
		So(err, ShouldBeNil)
	})
	Convey("Testing Delete", t, func() {
		err = ch.Delete()
		So(err, ShouldBeNil)
	})
	Convey("Testing Delete", t, func() {
		err = cdnft.Delete()
		So(err, ShouldBeNil)
	})
	Convey("Testing Delete", t, func() {
		err = vt.Delete()
		So(err, ShouldBeNil)
	})
	Convey("Testing Delete", t, func() {
		err = cdn.Delete()
		So(err, ShouldBeNil)
	})
	Convey("Testing Delete", t, func() {
		err = cat.Delete()
		So(err, ShouldBeNil)
	})
	Convey("Testing Delete", t, func() {
		err = v.Delete()
		So(err, ShouldBeNil)
	})

}

// //Data Modelling Helper Functions
// func getPartId() (id int) {
// 	db, err := sql.Open("mysql", database.ConnectionString())
// 	if err != nil {
// 		return 0
// 	}
// 	defer db.Close()
// 	stmt, err := db.Prepare("SELECT partID FROM VehiclePart ORDER BY RAND() LIMIT 1")
// 	if err != nil {
// 		return 0
// 	}
// 	defer stmt.Close()
// 	err = stmt.QueryRow().Scan(&id)
// 	return id
// }
// func getCatId() (id int) {
// 	db, err := sql.Open("mysql", database.ConnectionString())
// 	if err != nil {
// 		return 0
// 	}
// 	defer db.Close()
// 	stmt, err := db.Prepare("SELECT catID FROM Categories ORDER BY RAND() LIMIT 1")
// 	if err != nil {
// 		return 0
// 	}
// 	defer stmt.Close()
// 	err = stmt.QueryRow().Scan(&id)
// 	return id
// }
// func getChannelId() (id int) {
// 	db, err := sql.Open("mysql", database.ConnectionString())
// 	if err != nil {
// 		return 0
// 	}
// 	defer db.Close()
// 	stmt, err := db.Prepare("SELECT ID FROM Channel ORDER BY RAND() LIMIT 1")
// 	if err != nil {
// 		return 0
// 	}
// 	defer stmt.Close()
// 	err = stmt.QueryRow().Scan(&id)
// 	return id
// }
// func getFileId() (id int) {
// 	db, err := sql.Open("mysql", database.ConnectionString())
// 	if err != nil {
// 		return 0
// 	}
// 	defer db.Close()
// 	stmt, err := db.Prepare("SELECT ID FROM CdnFile ORDER BY RAND() LIMIT 1")
// 	if err != nil {
// 		return 0
// 	}
// 	defer stmt.Close()
// 	err = stmt.QueryRow().Scan(&id)
// 	return id
// }

// //Test Videos
// func TestVideo_New(t *testing.T) {
// 	Convey("Testing Videos", t, func() {
// 		var v Video
// 		var p products.Part
// 		var ch Channel
// 		var f CdnFile
// 		var c Category
// 		var err error
// 		//Create
// 		p.ID = 11000
// 		ch.ID = 1
// 		f.ID = 1
// 		c.ID = 1

// 		v.Title = "Test Video"
// 		v.Parts = append(v.Parts, p)
// 		v.Channels = append(v.Channels, ch)
// 		v.Files = append(v.Files, f)
// 		v.Categories = append(v.Categories, c)
// 		err = v.Create()
// 		So(err, ShouldBeNil)

// 		//Update
// 		p.ID = 110001
// 		ch.ID = 11
// 		f.ID = 11
// 		// c.ID = getChannelId()
// 		c.ID = 11
// 		ch.Type.ID = 1
// 		v.Title = "Test Video Redux"
// 		v.Parts = append(v.Parts, p)
// 		v.Channels = append(v.Channels, ch)
// 		v.Files = append(v.Files, f)

// 		v.Categories = append(v.Categories, c)

// 		err = v.Update()
// 		So(err, ShouldBeNil)

// 		//Get
// 		err = v.Get()
// 		So(err, ShouldBeNil)

// 		//Get Details
// 		err = v.GetVideoDetails()
// 		So(err, ShouldBeNil)
// 		t.Log(v)
// 		//Get Chans
// 		chans, err := v.GetChannels()
// 		if err != sql.ErrNoRows {
// 			So(err, ShouldBeNil)
// 			t.Log(v.ID)
// 			So(len(chans), ShouldBeGreaterThanOrEqualTo, 1)
// 		}

// 		//Get CDNs
// 		cdns, err := v.GetCdnFiles()
// 		if err != sql.ErrNoRows {
// 			So(err, ShouldBeNil)
// 			So(cdns, ShouldHaveSameTypeAs, CdnFiles{})
// 			So(len(cdns), ShouldBeGreaterThan, 0)
// 		}

// 		//Delete
// 		err = v.Delete()
// 		So(err, ShouldBeNil)

// 	})

// 	Convey("Testing Get-all", t, func() {
// 		//All Videos
// 		vs, err := GetAllVideos()
// 		if err != sql.ErrNoRows {
// 			So(err, ShouldBeNil)
// 			So(len(vs), ShouldBeGreaterThan, 0)
// 		}
// 		//Part Vids
// 		var p products.Part
// 		p.ID = getPartId()
// 		partVids, err := GetPartVideos(p)
// 		if err != sql.ErrNoRows {
// 			So(err, ShouldBeNil)
// 			So(len(partVids), ShouldBeGreaterThan, 0)
// 		}

// 	})
// }

// //Test Video-related items
// func TestVideoCatsnStuff_New(t *testing.T) {
// 	Convey("Testing Creates", t, func() {
// 		var err error
// 		var ch Channel
// 		var f CdnFile
// 		var v VideoType
// 		var cht ChannelType
// 		var ft CdnFileType

// 		ch.Description = "test"
// 		f.ObjectName = "test"
// 		v.Name = "test"
// 		cht.Name = "test"
// 		ft.Title = "test"

// 		err = ch.Create()
// 		So(err, ShouldBeNil)
// 		err = f.Create()
// 		So(err, ShouldBeNil)
// 		err = v.Create()
// 		So(err, ShouldBeNil)
// 		err = cht.Create()
// 		So(err, ShouldBeNil)
// 		err = ft.Create()
// 		So(err, ShouldBeNil)

// 		//Updates
// 		ch.Description = "new test"
// 		f.ObjectName = "new test"
// 		v.Name = "new test"
// 		cht.Name = "new test"
// 		ft.Title = "new test"

// 		err = ch.Update()
// 		So(err, ShouldBeNil)
// 		err = f.Update()
// 		So(err, ShouldBeNil)
// 		err = v.Update()
// 		So(err, ShouldBeNil)
// 		err = cht.Update()
// 		So(err, ShouldBeNil)
// 		err = ft.Update()
// 		So(err, ShouldBeNil)
// 		//Gets
// 		err = ch.Get()
// 		So(err, ShouldBeNil)
// 		err = f.Get()
// 		So(err, ShouldBeNil)
// 		err = v.Get()
// 		So(err, ShouldBeNil)
// 		err = cht.Get()
// 		So(err, ShouldBeNil)
// 		err = ft.Get()
// 		So(err, ShouldBeNil)

// 		Convey("Testing GetAll", func() {
// 			//Channels
// 			chs, err := GetAllChannels()
// 			if err != sql.ErrNoRows {
// 				So(err, ShouldBeNil)
// 				So(len(chs), ShouldBeGreaterThan, 0)
// 			}
// 			//CdnFiles
// 			cdns, err := GetAllCdnFiles()
// 			if err != sql.ErrNoRows {
// 				So(err, ShouldBeNil)
// 				So(len(cdns), ShouldBeGreaterThan, 0)
// 			}
// 			//VideoTypes
// 			vts, err := GetAllVideoTypes()
// 			if err != sql.ErrNoRows {
// 				So(err, ShouldBeNil)
// 				So(len(vts), ShouldBeGreaterThan, 0)
// 			}
// 			//ChannelTypes
// 			cts, err := GetAllChannelTypes()
// 			if err != sql.ErrNoRows {
// 				So(err, ShouldBeNil)
// 				So(len(cts), ShouldBeGreaterThan, 0)
// 			}

// 			//CdnFileTypes
// 			cds, err := GetAllCdnFileTypes()
// 			if err != sql.ErrNoRows {
// 				So(err, ShouldBeNil)
// 				So(len(cds), ShouldBeGreaterThan, 0)
// 			}
// 		})
// 		//Deletes
// 		err = ch.Delete()
// 		So(err, ShouldBeNil)
// 		err = f.Delete()
// 		So(err, ShouldBeNil)
// 		err = v.Delete()
// 		So(err, ShouldBeNil)
// 		err = cht.Delete()
// 		So(err, ShouldBeNil)
// 		err = ft.Delete()
// 		So(err, ShouldBeNil)

// 	})

// }
