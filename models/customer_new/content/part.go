package custcontent

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/conversions"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
	// "time"
)

var (
	customerPartContent_Grouped = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
									ct.type,ct.allowHTML,ccb.partID
									from CustomerContent as cc
									join CustomerContentBridge as ccb on cc.id = ccb.contentID
									join ContentType as ct on cc.typeID = ct.cTypeID
									join Customer as c on cc.custID = c.cust_id
									join CustomerUser as cu on c.cust_id = cu.cust_ID
									join ApiKey as ak on cu.id = ak.user_id
									where api_key = ? and ccb.partID IN (?)
									group by cc.id, ccb.partID`

	allCustomerPartContent = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
								ct.type,ct.allowHTML,ccb.partID
								from CustomerContent as cc
								join CustomerContentBridge as ccb on cc.id = ccb.contentID
								join ContentType as ct on cc.typeID = ct.cTypeID
								join Customer as c on cc.custID = c.cust_id
								join CustomerUser as cu on c.cust_id = cu.cust_ID
								join ApiKey as ak on cu.id = ak.user_id
								where api_key = ? and ccb.partID > 0
								group by ccb.partID, cc.id
								order by ccb.partID`

	customerPartContent = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
								ct.type,ct.allowHTML,ccb.partID
								from CustomerContent as cc
								join CustomerContentBridge as ccb on cc.id = ccb.contentID
								join ContentType as ct on cc.typeID = ct.cTypeID
								join Customer as c on cc.custID = c.cust_id
								join CustomerUser as cu on c.cust_id = cu.cust_ID
								join ApiKey as ak on cu.id = ak.user_id
								where api_key = ? and ccb.partID = ?
								group by cc.id`
)

type PartContent struct {
	PartId  int
	Content []CustomerContent
}

// Retrieves all part content for this customer
func GetAllPartContent(key string) (content []PartContent, err error) {

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return content, err
	}
	defer db.Close()

	stmt, err := db.Prepare(allCustomerPartContent)
	if err != nil {
		return content, err
	}

	res, err := stmt.Query(key)

	rawContent := make(map[int][]CustomerContent, 0)
	var partId int
	var added, modified, deleted []byte
	for res.Next() {
		var cc CustomerContent
		err = res.Scan(
			&cc.Id,
			&cc.Text,
			&added,
			&modified,
			&deleted, //Not Used
			&cc.ContentType.Type,
			&cc.ContentType.AllowHtml,
			&partId,
		)
		if err != nil {
			return content, err
		}
		cc.Added, _ = conversions.ByteToTime(added, timeYearFormat)
		cc.Modified, _ = conversions.ByteToTime(modified, timeYearFormat)
		part_id := partId
		if part_id > 0 {
			rawContent[part_id] = append(rawContent[part_id], cc)
		}

	}
	for k, _ := range rawContent {
		pCon := PartContent{
			PartId:  k,
			Content: rawContent[k],
		}
		content = append(content, pCon)
	}
	return

}

// Retrieves specific part content for this customer
func GetPartContent(partID int, key string) (content []CustomerContent, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return content, err
	}
	defer db.Close()

	stmt, err := db.Prepare(customerPartContent)
	if err != nil {
		return content, err
	}

	res, err := stmt.Query(key, partID)
	var partId int
	var added, modified, deleted []byte
	for res.Next() {
		var cc CustomerContent
		err = res.Scan(
			&cc.Id,
			&cc.Text,
			&added,
			&modified,
			&deleted, //Not Used
			&cc.ContentType.Type,
			&cc.ContentType.AllowHtml,
			&partId,
		)
		if err != nil {
			return content, err
		}
		cc.Added, _ = conversions.ByteToTime(added, timeYearFormat)
		cc.Modified, _ = conversions.ByteToTime(modified, timeYearFormat)
		content = append(content, cc)
	}
	return
}

func GetGroupedPartContent(ids []string, key string) (content map[int][]CustomerContent, err error) {
	content = make(map[int][]CustomerContent, len(ids))

	for i := 0; i < len(ids); i++ {
		intId, err := strconv.Atoi(ids[i])
		if err == nil {
			content[intId] = make([]CustomerContent, 0)
		}
	}
	escaped_key := database.Db.Escape(key)

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return content, err
	}
	defer db.Close()

	stmt, err := db.Prepare(customerPartContent_Grouped)
	if err != nil {
		return content, err
	}
	var partId int
	var added, modified, deleted []byte

	res, err := stmt.Query(escaped_key, strings.Join(ids, ","))
	for res.Next() {
		var cc CustomerContent
		err = res.Scan(
			&cc.Id,
			&cc.Text,
			&added,
			&modified,
			&deleted, //Not Used
			&cc.ContentType.Type,
			&cc.ContentType.AllowHtml,
			&partId,
		)
		if err != nil {
			return content, err
		}
		cc.Added, _ = conversions.ByteToTime(added, timeYearFormat)
		cc.Modified, _ = conversions.ByteToTime(modified, timeYearFormat)
		content[partId] = append(content[partId], cc)
	}
	return
}
