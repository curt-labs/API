package webProperty_model

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/api"
	"github.com/curt-labs/GoAPI/helpers/database"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func getRandomWebProperty() (wp WebProperty) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return wp
	}
	defer db.Close()

	stmt, err := db.Prepare("select id, name, cust_ID, badgeID, url, isEnabled,sellerID, typeID, isFinalApproved, isEnabledDate,isDenied, requestedDate, addedDate from WebProperties")
	if err != nil {
		return wp
	}

	var wps []WebProperty

	res, err := stmt.Query()
	for res.Next() {
		var w WebProperty
		err = res.Scan(
			&w.ID,
			&w.Name,
			&w.CustID,
			&w.BadgeID,
			&w.Url,
			&w.IsEnabled,
			&w.SellerID,
			&w.WebPropertyType.ID,
			&w.IsFinalApproved,
			&w.IsEnabledDate,
			&w.IsDenied,
			&w.RequestedDate,
			&w.AddedDate)
		if err != nil {
			return wp
		}
		wps = append(wps, w)
	}
	x := rand.Intn(len(wps))
	wp = wps[x]
	return
}

func getRandomNote() (wp WebPropertyNote) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return wp
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT * FROM WebPropNotes")
	if err != nil {
		return wp
	}

	var wps []WebPropertyNote

	res, err := stmt.Query()
	for res.Next() {
		var w WebPropertyNote
		err = res.Scan(&w.ID, &w.WebPropID, &w.Text, &w.DateAdded)
		if err != nil {
			return wp
		}
		wps = append(wps, w)
	}
	x := rand.Intn(len(wps))
	wp = wps[x]
	return
}

func getRandomType() (wp WebPropertyType) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return wp
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT * FROM WebPropertyTypes")
	if err != nil {
		return wp
	}

	var wps []WebPropertyType

	res, err := stmt.Query()
	for res.Next() {
		var w WebPropertyType
		err = res.Scan(&w.ID, &w.TypeID, &w.Type)
		if err != nil {
			return wp
		}
		wps = append(wps, w)
	}
	x := rand.Intn(len(wps))
	wp = wps[x]
	return
}

func getRandomRequirement() (wp WebPropertyRequirement) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return wp
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT * FROM WebPropRequirements")
	if err != nil {
		return wp
	}

	var wps []WebPropertyRequirement

	res, err := stmt.Query()
	for res.Next() {
		var w WebPropertyRequirement
		err = res.Scan(&w.ID, &w.ReqType, &w.Requirement)
		if err != nil {
			return wp
		}
		wps = append(wps, w)
	}
	x := rand.Intn(len(wps))
	wp = wps[x]
	return
}

func TestGetWebProperties(t *testing.T) {
	var wp WebProperty
	var rs WebPropertyRequirement
	Convey("Testing Gets", t, func() {
		Convey("Testing GetAll()", func() {
			var ws WebProperties
			ws, err := GetAll()
			So(ws, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(len(ws), ShouldNotBeNil)
			if len(ws) > 0 {
				wp = ws[0]
				Convey("Testing Get()", func() {
					var w WebProperty
					w = ws[api_helpers.RandGenerator(len(ws)-1)]
					err = w.Get()
					So(w, ShouldNotBeNil)
					So(err, ShouldBeNil)
					So(w.Name, ShouldNotEqual, "")
					So(w.CustID, ShouldNotEqual, 0)
					So(w.WebPropertyType.Type, ShouldNotEqual, "")
					var ns WebPropertyNotes
					var rs WebPropertyRequirements
					So(w.WebPropertyNotes, ShouldHaveSameTypeAs, ns)
					So(w.WebPropertyRequirements, ShouldHaveSameTypeAs, rs)
				})
				Convey("Testing Get(); focus on dates", func() {
					var w WebProperty
					w = ws[api_helpers.RandGenerator(len(ws)-1)]
					err = w.Get()
					So(w, ShouldNotBeNil)
					var t time.Time
					So(w.IsEnabledDate, ShouldHaveSameTypeAs, t)
					So(w.RequestedDate, ShouldHaveSameTypeAs, t)

					So(w.AddedDate, ShouldHaveSameTypeAs, t)
					So(err, ShouldBeNil)
				})
			}

			Convey("Testing GetAllWebPropertyRequirements()", func() {
				qs, err := GetAllWebPropertyRequirements()
				if len(qs) > 0 {
					rs = qs[0]
				}
				So(qs, ShouldNotBeNil)
				So(err, ShouldBeNil)
			})
		})

		Convey("Testing GetAllWebPropertyTypes()", func() {
			qs, err := GetAllWebPropertyTypes()
			So(qs, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("Testing GetNote()", func() {
			testNote := getRandomNote()
			var n WebPropertyNote
			n.ID = testNote.ID
			err := n.Get()
			So(n.Text, ShouldEqual, testNote.Text)
			So(err, ShouldBeNil)
		})
		Convey("Testing GetWebPropertyType()", func() {
			testType := getRandomType()
			var n WebPropertyType
			n.ID = testType.ID
			err := n.Get()
			So(n.Type, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(n.Type, ShouldEqual, testType.Type)
		})
		Convey("Testing GetWebPropertyRequirement()", func() {
			testReq := getRandomRequirement()
			var n WebPropertyRequirement
			n.RequirementID = testReq.ID
			err := n.Get()
			So(n.Requirement, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(n.Requirement, ShouldEqual, testReq.Requirement)
		})
		// Convey("Testing GetWebPropertyRequirementCheck()", func() {
		// 	n := getRandomRequirement()
		// 	err := n.GetJoin()
		// 	So(n.WebPropID, ShouldNotBeNil)
		// 	So(err, ShouldBeNil)
		// })
		Convey("Testing GetAllWebPropertyNotes()", func() {
			qs, err := GetAllWebPropertyNotes()
			So(qs, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("Testing Search()", func() {
			testProp := getRandomWebProperty()
			as, err := Search(testProp.Name, "", "", "", "", "", "", "", "", "", "", "", "1", "0")
			So(as, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(as.Pagination.Page, ShouldEqual, 1)
			So(as.Pagination.ReturnedCount, ShouldNotBeNil)
			So(as.Pagination.PerPage, ShouldNotBeNil)
			So(as.Pagination.PerPage, ShouldEqual, len(as.Objects))
		})
	})
	Convey("Testing CUD", t, func() {
		Convey("Testing Create(), Update(), Delete()", func() {
			var f WebProperty
			var err error
			f.Name = "testTitle"
			f.CustID = 12345
			f.BadgeID = strconv.Itoa(rand.Int())
			f.IsEnabledDate, err = time.Parse(timeFormat, "2004-03-03 9:15:00")
			f.Url = "www.test.com"
			f.IsEnabled = true
			f.SellerID = "test"
			f.WebPropertyType.ID = 2
			f.IsFinalApproved = true
			f.IsDenied = false
			f.RequestedDate, err = time.Parse(timeFormat, "2004-03-03 9:15:00")

			f.Create()
			So(f.ID, ShouldNotBeNil)
			f.Get()
			So(f, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(f.Name, ShouldEqual, "testTitle")
			So(f.CustID, ShouldEqual, 12345)
			So(f.BadgeID, ShouldNotBeNil)
			So(f.Url, ShouldEqual, "www.test.com")
			So(f.SellerID, ShouldEqual, "test")
			So(f.WebPropertyType.ID, ShouldEqual, 2)
			So(f.IsFinalApproved, ShouldBeTrue)
			So(f.IsDenied, ShouldBeFalse)
			t, err := time.Parse(timeFormat, "2004-03-03 09:15:00")
			So(f.IsEnabledDate, ShouldResemble, t)
			u, err := time.Parse(timeFormat, "2004-03-03 09:15:00")
			So(f.RequestedDate, ShouldResemble, u)
			f.Name = "testTitle2"
			f.CustID = 123452
			f.BadgeID = strconv.Itoa(rand.Int())
			f.IsEnabledDate, err = time.Parse(timeFormat, "2004-03-03 9:15:22")
			f.Url = "www.test.com2"
			f.IsEnabled = false
			f.SellerID = "test2"
			f.WebPropertyType.ID = 22
			f.IsFinalApproved = false
			f.IsDenied = true
			f.RequestedDate, err = time.Parse(timeFormat, "2004-03-03 9:15:22")
			f.Update()
			So(f.ID, ShouldNotBeNil)
			f.Get()
			So(f, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(f.Name, ShouldEqual, "testTitle2")
			So(f.CustID, ShouldEqual, 123452)
			So(f.BadgeID, ShouldNotBeNil)
			So(f.Url, ShouldEqual, "www.test.com2")
			So(f.SellerID, ShouldEqual, "test2")
			So(f.WebPropertyType.ID, ShouldEqual, 22)
			So(f.IsFinalApproved, ShouldBeFalse)
			So(f.IsDenied, ShouldBeTrue)
			t, err = time.Parse(timeFormat, "2004-03-03 09:15:22")
			So(f.IsEnabledDate, ShouldHaveSameTypeAs, t)
			u, err = time.Parse(timeFormat, "2004-03-03 09:15:22")
			So(f.RequestedDate, ShouldHaveSameTypeAs, u)
			err = f.Get()
			err = f.Delete()
			So(err, ShouldBeNil)
		})
		Convey("Testing Create(), Update(), Delete() WebPropNotes", func() {
			var n WebPropertyNote
			f := getRandomWebProperty()
			n.WebPropID = f.ID
			n.Text = "test note"
			c := make(chan int)
			go func() {
				n.Create()
				c <- 1
			}()
			<-c
			f.Get()
			So(f.Name, ShouldHaveSameTypeAs, "")
			var ns WebPropertyNotes
			So(f.WebPropertyNotes, ShouldHaveSameTypeAs, ns)

			n.Text = "Funk"
			err := n.Update()
			So(err, ShouldBeNil)

			err = n.Delete()
			So(err, ShouldBeNil)

		})
		Convey("Testing Create(), Update(), Delete() WebProperyRequirementsCheck", func() {

			var err error

			rs.WebPropID = wp.ID
			rs.RequirementID = rs.ID
			rs.Compliance = true

			c := make(chan int)
			go func() {
				rs.CreateJoin()
				c <- 1
			}()
			<-c
			err = wp.Get()
			var tmpReq WebPropertyRequirements
			So(wp.WebPropertyRequirements, ShouldHaveSameTypeAs, tmpReq)
			So(err, ShouldBeNil)

			rs.Compliance = true
			err = rs.UpdateJoin()
			So(err, ShouldBeNil)

			err = rs.DeleteJoin()
			So(err, ShouldBeNil)
		})

		Convey("Testing Create(), Update(), Delete() Requirement", func() {
			var n WebPropertyRequirement
			n.ReqType = "Approved"
			n.Requirement = "TEST Requirement"
			err := n.Create()
			So(err, ShouldBeNil)
			err = n.Get()
			n.Requirement = "booger"
			err = n.Update()
			So(err, ShouldBeNil)
			So(n.Requirement, ShouldEqual, "booger")
			err = n.Get()
			err = n.Delete()
			So(err, ShouldBeNil)
		})

		Convey("Testing Create(), Update(), Delete() Type", func() {
			var n WebPropertyType
			n.Type = "TEST Type"
			err := n.Create()
			So(err, ShouldBeNil)
			err = n.Get()
			n.Type = "booger"
			err = n.Update()
			So(err, ShouldBeNil)
			So(n.Type, ShouldEqual, "booger")
			err = n.Get()
			err = n.Delete()
			So(err, ShouldBeNil)

		})
	})
}
