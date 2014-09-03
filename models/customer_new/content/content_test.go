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
			err = content.Delete(partID, catID, key)
			So(err, ShouldBeNil)
		})
		Convey("Testing Save()Update", func() {
			partID, catID := easyCatPart()
			_, key := getApiKey(allCustContent)
			var content CustomerContent
			content.Text = "test text"
			content.ContentType.Type = getRandType()
			content.Id = 1
			err := content.Save(partID, catID, key)
			So(err, ShouldBeNil)
			err = content.Delete(partID, catID, key)
			So(err, ShouldBeNil)
		})
		Convey("Testing GetContentType()", func() {
			var c CustomerContent
			var t ContentType
			c.ContentType.Type = getRandType()
			ct, err := c.GetContentType()
			So(err, ShouldBeNil)
			So(ct, ShouldNotBeNil)
			So(ct, ShouldHaveSameTypeAs, t)

		})
		Convey("AllCustomerContentTypes()", func() {
			cts, err := AllCustomerContentTypes()
			So(err, ShouldBeNil)
			So(cts, ShouldNotBeNil)
		})
	})
}

func TestContentComparedToOldModel(t *testing.T) {
	Convey("ComparativeTests", t, func() {
		err := database.PrepareAll()
		So(err, ShouldBeNil)

		//Works, but dateModifed does not work in original model
		// Convey("AllContent v AllContent", func() {
		// 	_, key := getApiKey(allCustContent)
		// 	content, err := AllCustomerContent(key)
		// 	So(err, ShouldBeNil)
		// 	oldContent, err := custcontent.AllCustomerContent(key)
		// 	So(err, ShouldBeNil)
		// 	So(content, ShouldResemble, oldContent)
		// })
	})
}
