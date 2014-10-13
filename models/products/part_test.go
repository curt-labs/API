package products

import (
	// "database/sql"
	// "github.com/curt-labs/GoAPI/helpers/database"
	. "github.com/smartystreets/goconvey/convey"
	// "log"
	"math/rand"
	"testing"
	"time"
)

// func checkPartNumber(num int) (err error) {
// 	db, err := sql.Open("mysql", database.ConnectionString())
// 	if err != nil {
// 		return err
// 	}
// 	defer db.Close()
// 	stmt, err := db.Prepare("SELECT partID FROM Part WHERE partID = ?")
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()
// 	var z int
// 	err = stmt.QueryRow().Scan(&z)
// 	return err
// }

func randNum() (y int) {
	rand.Seed(int64(time.Second))
	y = rand.Intn(1000) + 9000000
	return y
}

func TestParts(t *testing.T) {
	Convey("Testing Create Part", t, func() {
		var p Part
		var err error

		// create
		p.ID = randNum()
		p.Status = 800
		p.ShortDesc = "Test"
		err = p.Create()
		So(err, ShouldBeNil)
		So(p, ShouldNotBeNil)
		t.Log(p.ID)

		// get
		err = p.Get(generateAPIkey())
		So(err, ShouldBeNil)

		//update
		p.Status = 900
		p.ShortDesc = "Test2"
		err = p.Update()
		So(err, ShouldBeNil)

		// delete
		err = p.Delete()
		So(err, ShouldBeNil)

	})
	// Convey("Testing Create - Long Version", t, func() {
	// 	var p Part
	// 	var err error
	// 	var a Attribute
	// 	var i Installation
	// 	var c Content
	// 	var pr Price
	// 	var r Review
	// 	var im Image
	// 	var related int
	// 	var cat Category
	// 	var v Video
	// 	var pa Package

	// 	p.ID = randNum()
	// 	p.Status = 800
	// 	p.ShortDesc = "Test"

	// })
}

func TestLatest(t *testing.T) {
	Convey("Testing GetLatest", t, func() {
		parts, err := Latest(generateAPIkey(), 10)
		So(err, ShouldBeNil)
		So(parts, ShouldHaveSameTypeAs, []Part{})
	})
}
