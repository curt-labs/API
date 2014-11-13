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
	var cc CustomerContent
	var ct ContentType
	var err error
	_, key := getApiKey(allCustContent)
	partID, catID := easyCatPart()

	Convey("Testing Create", t, func() {
		cc.Text = "test text"
		err = cc.Save(11000, 1, key)

		ct.Type = "test type"
		err = ct.Create()
		So(err, ShouldBeNil)

	})

	Convey("Testing AllCustomerContent()", t, func() {
		_, key := getApiKey(allCustContent)
		content, err := AllCustomerContent(key)
		var allCn []CustomerContent
		So(content, ShouldHaveSameTypeAs, allCn)
		So(err, ShouldBeNil)
	})
	Convey("Testing GetCustomerContent()", t, func() {
		id, key := getApiKey(allCustContent)
		content, err := GetCustomerContent(id, key)
		var cn CustomerContent
		So(content, ShouldHaveSameTypeAs, cn)
		So(err, ShouldBeNil)
	})
	Convey("Testing GetCustomerContentRevisions()", t, func() {
		var con []CustomerContentRevision
		id, key := getApiKey(allCustContent)
		content, err := GetCustomerContentRevisions(id, key)
		So(content, ShouldHaveSameTypeAs, con)
		So(err, ShouldBeNil)
	})
	Convey("Testing Save()", t, func() {
		var err error
		partID, catID := easyCatPart()
		_, key := getApiKey(allCustContent)
		cc.Text = "test text"
		cc.ContentType = ct
		err = cc.Save(partID, catID, key)
		So(err, ShouldBeNil)

	})
	Convey("Testing Save()Update", t, func() {
		_, key := getApiKey(allCustContent)
		cc.Text = "test text"
		cc.Id = 1
		cc.ContentType = ct
		err := cc.Save(partID, catID, key)
		So(err, ShouldBeNil)

	})
	Convey("Testing GetContentType()", t, func() {
		cc.ContentType = ct
		err := cc.GetContentType()
		So(err, ShouldBeNil)
		So(cc.ContentType, ShouldNotBeNil)

	})
	Convey("AllCustomerContentTypes()", t, func() {
		cts, err := AllCustomerContentTypes()
		So(err, ShouldBeNil)
		So(cts, ShouldNotBeNil)
	})
	Convey("Test Delete", t, func() {
		err = ct.Delete()
		So(err, ShouldBeNil)
		_ = cc.Delete(partID, catID, key) //returns error if no bridge exists -ok
	})

}
