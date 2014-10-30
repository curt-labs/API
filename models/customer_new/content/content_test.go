package custcontent

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
)

const (
	inputTimeFormat = "01/02/2006"
)

var (
	allCustContent = `select cc.id, api_key
							from CustomerContent as cc
							left join CustomerContentBridge as ccb on cc.id = ccb.contentID
							join ContentType as ct on cc.typeID = ct.cTypeID
							join Customer as c on cc.custID = c.cust_id
							join CustomerUser as cu on c.cust_id = cu.cust_ID
							join ApiKey as ak on cu.id = ak.user_id
							group by cc.id`
	randomType     = `SELECT type FROM ContentType ORDER BY RAND() LIMIT 1`
	easyCatAndPart = `SELECT catID, partID FROM CustomerContentBridge ORDER BY RAND() LIMIT 1`
)

type Output struct {
	apiKey string
	id     int
}

func getApiKey(query string) (int, string) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return 0, ""
	}
	defer db.Close()

	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, ""
	}

	var outputs []Output

	res, err := stmt.Query()
	for res.Next() {
		var output Output
		res.Scan(&output.id, &output.apiKey)
		if err != nil {
			return 0, ""
		}
		outputs = append(outputs, output)
	}
	if len(outputs) == 0 {
		return 0, ""
	}

	x := rand.Intn(len(outputs))
	return outputs[x].id, outputs[x].apiKey
}

func getRandType() (t string) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ""
	}
	defer db.Close()

	stmt, err := db.Prepare(randomType)
	if err != nil {
		return ""
	}
	err = stmt.QueryRow().Scan(&t)
	return
}

func easyCatPart() (partID, catID int) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(easyCatAndPart)
	if err != nil {
		return
	}
	err = stmt.QueryRow().Scan(&partID, &catID)
	return
}

func TestContent(t *testing.T) {
	Convey("Testing Content", t, func() {

		Convey("Testing AllCustomerContent()", func() {
			_, key := getApiKey(allCustContent)
			content, err := AllCustomerContent(key)
			var allCn []CustomerContent
			So(content, ShouldHaveSameTypeAs, allCn)
			So(err, ShouldBeNil)
		})
		Convey("Testing GetCustomerContent()", func() {
			id, key := getApiKey(allCustContent)
			content, err := GetCustomerContent(id, key)
			var cn CustomerContent
			So(content, ShouldHaveSameTypeAs, cn)
			So(err, ShouldBeNil)
		})
		Convey("Testing GetCustomerContentRevisions()", func() {
			var con []CustomerContentRevision
			id, key := getApiKey(allCustContent)
			content, err := GetCustomerContentRevisions(id, key)
			So(content, ShouldHaveSameTypeAs, con)
			So(err, ShouldBeNil)
		})
		Convey("Testing Save()", func() {
			var err error
			partID, catID := easyCatPart()
			_, key := getApiKey(allCustContent)
			var content CustomerContent
			content.Text = "test text"
			content.ContentType.Type = getRandType()
			err = content.Save(partID, catID, key)
			So(err, ShouldBeNil)

			_ = content.Delete(partID, catID, key) //returns error if no bridge exists -ok

		})
		Convey("Testing Save()Update", func() {
			partID, catID := easyCatPart()
			_, key := getApiKey(allCustContent)
			var content CustomerContent
			content.Text = "test text"
			content.Id = 1
			content.ContentType.Type = getRandType()
			err := content.Save(partID, catID, key)
			So(err, ShouldBeNil)

			t.Log(content.GetContentType())
			err = content.Save(partID, catID, key)
			So(err, ShouldBeNil)
			_ = content.Delete(partID, catID, key) //returns error if no bridge exists -ok
		})
		Convey("Testing GetContentType()", func() {
			var c CustomerContent
			c.ContentType.Type = getRandType()
			err := c.GetContentType()
			So(err, ShouldBeNil)
			So(c.ContentType, ShouldNotBeNil)

		})
		Convey("AllCustomerContentTypes()", func() {
			cts, err := AllCustomerContentTypes()
			So(err, ShouldBeNil)
			So(cts, ShouldNotBeNil)
		})
	})
}

//Comparisons to old customer content model
// func TestContentComparedToOldModel(t *testing.T) {
// 	Convey("ComparativeTests", t, func() {
// 		id, key := getApiKey(allCustContent)

// 		allCon, err := AllCustomerContent(key)
// 		So(err, ShouldBeNil)
// 		var c CustomerContent
// 		if len(allCon) > 0 {
// 			c = allCon[rand.Intn(len(allCon))]
// 		}

// 		allCon2, err := custcontent.AllCustomerContent(key)
// 		var c2 custcontent.CustomerContent
// 		if len(allCon2) > 0 {
// 			c2 = allCon2[rand.Intn(len(allCon2))]
// 		}

// 		Convey("AllContent", func() {
// 			content, err := AllCustomerContent(key)
// 			So(err, ShouldBeNil)
// 			oldContent, err := custcontent.AllCustomerContent(key)
// 			So(err, ShouldBeNil)
// 			So(len(content), ShouldEqual, len(oldContent))
// 		})
// 		Convey("Content Revisions", func() {
// 			content, err := GetCustomerContentRevisions(id, key)
// 			So(err, ShouldBeNil)
// 			oldContent, err := custcontent.GetCustomerContentRevisions(id, key)
// 			So(err, ShouldBeNil)
// 			So(len(content), ShouldEqual, len(oldContent))
// 		})
// 		Convey("ContentType", func() {
// 			indexedType, err := c.GetContentType()
// 			So(err, ShouldBeNil)
// 			oldindexedType, err := c2.GetContentType()
// 			So(err, ShouldBeNil)
// 			So(indexedType.Type, ShouldEqual, oldindexedType.Type)
// 		})
// 		Convey("AllCustContentTypes", func() {
// 			types, err := AllCustomerContentTypes()
// 			So(err, ShouldBeNil)
// 			oldTypes, err := custcontent.AllCustomerContentTypes()
// 			So(err, ShouldBeNil)
// 			So(len(types), ShouldEqual, len(oldTypes))
// 		})
// 	})
// }
