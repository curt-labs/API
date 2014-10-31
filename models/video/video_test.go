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
