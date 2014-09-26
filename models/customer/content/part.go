package custcontent

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/curt-labs/GoAPI/helpers/api"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllCustomerPartContentStmt = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
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
	getPartContentStmt = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
                         ct.type,ct.allowHTML,ccb.partID
                         from CustomerContent as cc
                         join CustomerContentBridge as ccb on cc.id = ccb.contentID
                         join ContentType as ct on cc.typeID = ct.cTypeID
                         join Customer as c on cc.custID = c.cust_id
                         join CustomerUser as cu on c.cust_id = cu.cust_ID
                         join ApiKey as ak on cu.id = ak.user_id
                         where api_key = ? and ccb.partID = ?
                         group by cc.id`
	getGroupedCustomerPartContentStmt = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
                                        ct.type,ct.allowHTML,ccb.partID
                                        from CustomerContent as cc
                                        join CustomerContentBridge as ccb on cc.id = ccb.contentID
                                        join ContentType as ct on cc.typeID = ct.cTypeID
                                        join Customer as c on cc.custID = c.cust_id
                                        join CustomerUser as cu on c.cust_id = cu.cust_ID
                                        join ApiKey as ak on cu.id = ak.user_id
                                        where api_key = ? and ccb.partID IN (?)
                                        group by cc.id, ccb.partID`
)

type PartContent struct {
	PartId  int
	Content []CustomerContent
}

// Retrieves all part content for this customer
func GetAllPartContent(key string) (content []PartContent, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllCustomerPartContentStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(key)
	if err != nil {
		return
	}

	rawContent := make(map[int][]CustomerContent, 0)

	for rows.Next() {
		var c CustomerContent
		var partID int
		err = rows.Scan(
			&c.Id,
			&c.Text,
			&c.Added,
			&c.Modified,
			&c.Hidden,
			&c.ContentType.Type,
			&c.ContentType.AllowHtml,
			&partID,
		)
		if err != nil {
			return
		}

		if partID > 0 {
			rawContent[partID] = append(rawContent[partID], c)
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
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getPartContentStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(key, partID)
	if err != nil {
		return
	}

	for rows.Next() {
		var c CustomerContent
		err = rows.Scan(
			&c.Id,
			&c.Text,
			&c.Added,
			&c.Modified,
			&c.Hidden,
			&c.ContentType.Type,
			&c.ContentType.AllowHtml,
		)
		if err != nil {
			return
		}
		content = append(content, c)
	}

	return
}

func GetGroupedPartContent(ids []string, key string) (content map[int][]CustomerContent, err error) {
	content = make(map[int][]CustomerContent, len(ids))
	escaped_key := api_helpers.Escape(key)

	for i := 0; i < len(ids); i++ {
		intId, err := strconv.Atoi(ids[i])
		if err == nil {
			content[intId] = make([]CustomerContent, 0)
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getGroupedCustomerPartContentStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(escaped_key, strings.Join(ids, ","))
	if err != nil {
		return
	}

	for rows.Next() {
		var c CustomerContent
		var partID int
		err = rows.Scan(
			&c.Id,
			&c.Text,
			&c.Added,
			&c.Modified,
			&c.Hidden,
			&c.ContentType.Type,
			&c.ContentType.AllowHtml,
			&partID,
		)
		if err != nil {
			return
		}
		c.ContentType.Type = "Part:" + c.ContentType.Type
		content[partID] = append(content[partID], c)
	}

	return
}
