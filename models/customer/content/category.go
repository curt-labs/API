package custcontent

import (
	"database/sql"

	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllCustomerCategoryContentStmt = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
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
	getCategoryContentStmt = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
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

type CategoryContent struct {
	CategoryId int
	Content    []CustomerContent
}

// Retrieves all category content for this customer
func GetAllCategoryContent(key string) (content []CategoryContent, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllCustomerCategoryContentStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(key)
	if err != nil {
		return
	}

	var catID int

	rawContent := make(map[int][]CustomerContent, 0)

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
			&catID,
		)
		if err != nil {
			return
		}

		if catID > 0 {
			rawContent[catID] = append(rawContent[catID], c)
		}
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
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getCategoryContentStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(key, catID)
	if err != nil {
		return
	}

	for rows.Next() {
		var c CustomerContent
		var categoryID int

		err = rows.Scan(
			&c.Id,
			&c.Text,
			&c.Added,
			&c.Modified,
			&c.Hidden,
			&c.ContentType.Type,
			&c.ContentType.AllowHtml,
			&categoryID,
		)
		if err != nil {
			return
		}
		content = append(content, c)
	}

	return
}
