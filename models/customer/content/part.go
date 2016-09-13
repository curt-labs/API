package custcontent

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/curt-labs/API/helpers/api"
	"github.com/curt-labs/API/helpers/database"
	_ "github.com/go-sql-driver/mysql"
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
	defer stmt.Close()
	res, err := stmt.Query(key)

	rawContent := make(map[int][]CustomerContent, 0)
	var partId int
	var deleted *bool
	var added, modified *time.Time
	var ctype string
	for res.Next() {
		var cc CustomerContent
		err = res.Scan(
			&cc.Id,
			&cc.Text,
			&added,
			&modified,
			&deleted, //Not Used
			&ctype,
			&cc.ContentType.AllowHtml,
			&partId,
		)
		if err != nil {
			return content, err
		}
		cc.ContentType.Type = "Part:" + ctype
		if added != nil {
			cc.Added = *added
		}
		if modified != nil {
			cc.Modified = *modified
		}
		if deleted != nil {
			cc.Hidden = *deleted
		}
		part_id := partId
		if part_id > 0 {
			rawContent[part_id] = append(rawContent[part_id], cc)
		}

	}
	defer res.Close()

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
	content = make([]CustomerContent, 0) // initializer

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return content, err
	}
	defer db.Close()

	stmt, err := db.Prepare(customerPartContent)
	if err != nil {
		return content, err
	}
	defer stmt.Close()

	res, err := stmt.Query(key, partID)
	var partId int
	var deleted *bool
	var added, modified *time.Time
	var ctype string
	for res.Next() {
		var cc CustomerContent
		err = res.Scan(
			&cc.Id,
			&cc.Text,
			&added,
			&modified,
			&deleted, //Not Used
			&ctype,
			&cc.ContentType.AllowHtml,
			&partId,
		)
		if err != nil {
			return content, err
		}
		cc.ContentType.Type = "Part:" + ctype
		if added != nil {
			cc.Added = *added
		}
		if modified != nil {
			cc.Modified = *modified
		}
		if deleted != nil {
			cc.Hidden = *deleted
		}
		content = append(content, cc)
	}
	defer res.Close()
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
	escaped_key := api_helpers.Escape(key)

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return content, err
	}
	defer db.Close()

	stmt, err := db.Prepare(customerPartContent_Grouped)
	if err != nil {
		return content, err
	}
	defer stmt.Close()
	var partId int
	var deleted *bool
	var added, modified *time.Time
	var ctype string

	res, err := stmt.Query(escaped_key, strings.Join(ids, ","))
	for res.Next() {
		var cc CustomerContent
		err = res.Scan(
			&cc.Id,
			&cc.Text,
			&added,
			&modified,
			&deleted, //Not Used
			&ctype,
			&cc.ContentType.AllowHtml,
			&partId,
		)
		if err != nil {
			return content, err
		}
		cc.ContentType.Type = "Part:" + ctype
		if added != nil {
			cc.Added = *added
		}
		if modified != nil {
			cc.Modified = *modified
		}
		if deleted != nil {
			cc.Hidden = *deleted
		}
		content[partId] = append(content[partId], cc)
	}
	defer res.Close()
	return
}
