package custcontent

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/conversions"
	"github.com/curt-labs/GoAPI/helpers/database"
	// "time"
)

type CategoryContent struct {
	CategoryId int
	Content    []CustomerContent
}

var (
	allCustomerCategoryContent = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
									ct.type,ct.allowHTML,ccb.catID
									from CustomerContent as cc
									join CustomerContentBridge as ccb on cc.id = ccb.contentID
									join ContentType as ct on cc.typeID = ct.cTypeID
									join Customer as c on cc.custID = c.cust_id
									join CustomerUser as cu on c.cust_id = cu.cust_ID
									join ApiKey as ak on cu.id = ak.user_id
									where api_key = ? and ccb.catID > 0
									group by ccb.catID, cc.id
									order by ccb.catID`
	customerCategoryContent = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
									ct.type,ct.allowHTML,ccb.catID
									from CustomerContent as cc
									join CustomerContentBridge as ccb on cc.id = ccb.contentID
									join ContentType as ct on cc.typeID = ct.cTypeID
									join Customer as c on cc.custID = c.cust_id
									join CustomerUser as cu on c.cust_id = cu.cust_ID
									join ApiKey as ak on cu.id = ak.user_id
									where api_key = ? and ccb.catID = ?
									group by cc.id`
)

//TODO test me!

// Retrieves all category content for this customer
func GetAllCategoryContent(key string) (content []CategoryContent, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return content, err
	}
	defer db.Close()

	stmt, err := db.Prepare(allCustomerCategoryContent)
	if err != nil {
		return content, err
	}
	res, err := stmt.Query(key)
	var deleted bool
	var catID int
	var added, modified []byte
	rawContent := make(map[int][]CustomerContent, 0)
	for res.Next() {
		var c CustomerContent
		err = res.Scan(
			&c.Id,
			&c.Text,
			&added,
			&modified,
			&deleted,
			&c.ContentType.Type,
			&c.ContentType.AllowHtml,
			&catID,
		)

		if catID > 0 {
			rawContent[catID] = append(rawContent[catID], c)
		}
		c.Added, err = conversions.ByteToTime(added, timeYearFormat)
		c.Modified, err = conversions.ByteToTime(modified, timeYearFormat)
	}

	for k, _ := range rawContent {
		catCon := CategoryContent{
			CategoryId: k,
			Content:    rawContent[k],
		}
		content = append(content, catCon)
	}
	return
}

// Retrieves specific category content for this customer
func GetCategoryContent(catID int, key string) (content []CustomerContent, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return content, err
	}
	defer db.Close()

	stmt, err := db.Prepare(customerCategoryContent)
	if err != nil {
		return content, err
	}
	res, err := stmt.Query(key, catID)
	var deleted bool
	var added, modified []byte

	for res.Next() {
		var c CustomerContent
		err = res.Scan(
			&c.Id,
			&c.Text,
			&added,
			&modified,
			&deleted,
			&c.ContentType.Type,
			&c.ContentType.AllowHtml,
			&catID,
		)
		c.Added, err = conversions.ByteToTime(added, timeYearFormat)
		c.Modified, err = conversions.ByteToTime(modified, timeYearFormat)
		content = append(content, c)
	}
	return
}
