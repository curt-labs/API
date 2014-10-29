package webProperty_model

import (
	"database/sql"
	// "github.com/curt-labs/GoAPI/helpers/api"
	// "github.com/curt-labs/GoAPI/helpers/database"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestWebPropertiesBetter(t *testing.T) {
	Convey("Testing Create", t, func() {
		var w WebProperty
		var wr WebPropertyRequirement
		var wn WebPropertyNote
		var wt WebPropertyType
		var err error

		//New WebProperty
		w.Name = "test prop"
		w.Url = "www.hotdavid.com"

		//make up badge
		seed := int64(time.Now().Second() + time.Now().Minute() + time.Now().Hour() + time.Now().Year())
		rand.Seed(seed)
		w.BadgeID = strconv.Itoa(rand.Int()) //random badge

		//Test Requirement
		wr.ReqType = "Req Type"
		err = wr.Create()
		So(err, ShouldBeNil)

		wr.Compliance = true
		err = wr.Update()
		So(err, ShouldBeNil)

		err = wr.Get()
		So(err, ShouldBeNil)

		//Test Note
		wn.Text = "Note text"
		err = wn.Create()
		So(err, ShouldBeNil)

		wn.Text = "New Text"
		err = wn.Update()
		So(err, ShouldBeNil)

		err = wn.Get()
		So(err, ShouldBeNil)

		//Test Type
		wt.Type = "A type"
		err = wt.Create()
		So(err, ShouldBeNil)

		wt.Type = "B type"
		err = wt.Update()
		So(err, ShouldBeNil)

		err = wt.Get()
		So(err, ShouldBeNil)

		w.WebPropertyRequirements = append(w.WebPropertyRequirements, wr)
		w.WebPropertyNotes = append(w.WebPropertyNotes, wn)
		w.WebPropertyType = wt

		//Create Web Property
		err = w.Create()
		So(err, ShouldBeNil)
		So(w, ShouldNotBeNil)

		//Search
		obj, err := Search(w.Name, "", "", "", "", "", "", "", "", "", "", "", "1", "1")
		So(err, ShouldBeNil)
		So(len(obj.Objects), ShouldEqual, 0)

		//Update Property
		w.Name = "New Name"
		err = w.Update()
		So(err, ShouldBeNil)
		//Get Property
		err = w.Get()
		So(err, ShouldBeNil)

		//Deletes
		err = w.Delete()
		So(err, ShouldBeNil)
		err = wn.Delete()
		So(err, ShouldBeNil)
		err = wt.Delete()
		So(err, ShouldBeNil)

		err = wr.Delete()
		So(err, ShouldBeNil)

		Convey("Testing GetAll", func() {
			ws, err := GetAll()
			if err != sql.ErrNoRows {
				So(err, ShouldBeNil)
				So(len(ws), ShouldBeGreaterThan, 0)
			}
			ns, err := GetAllWebPropertyNotes()
			if err != sql.ErrNoRows {
				So(err, ShouldBeNil)
				So(len(ns), ShouldBeGreaterThan, 0)
			}
			rs, err := GetAllWebPropertyRequirements()
			if err != sql.ErrNoRows {
				So(err, ShouldBeNil)
				So(len(rs), ShouldBeGreaterThan, 0)
			}
			ts, err := GetAllWebPropertyTypes()
			if err != sql.ErrNoRows {
				So(err, ShouldBeNil)
				So(len(ts), ShouldBeGreaterThan, 0)
			}

		})

	})

}
