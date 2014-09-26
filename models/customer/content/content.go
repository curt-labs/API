package custcontent

import (
	"database/sql"
	"errors"
	"html"
	"strings"
	"time"

	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/models/customer"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllCustomerContentStmt = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
                                ct.type,ct.allowHTML,
                                ccb.partID, ccb.catID
                                from CustomerContent as cc
                                left join CustomerContentBridge as ccb on cc.id = ccb.contentID
                                join ContentType as ct on cc.typeID = ct.cTypeID
                                join Customer as c on cc.custID = c.cust_id
                                join CustomerUser as cu on c.cust_id = cu.cust_ID
                                join ApiKey as ak on cu.id = ak.user_id
                                where api_key = ?
                                group by cc.id`
	getCustomerContentStmt = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
                             ct.type,ct.allowHTML,ccb.partID,ccb.catID
                             from CustomerContent as cc
                             join CustomerContentBridge as ccb on cc.id = ccb.contentID
                             join ContentType as ct on cc.typeID = ct.cTypeID
                             join Customer as c on cc.custID = c.cust_id
                             join CustomerUser as cu on c.cust_id = cu.cust_ID
                             join ApiKey as ak on cu.id = ak.user_id
                             where api_key = ? and cc.id = ? limit 1`
	getCustomerContentRevisionsStmt = `select ccr.old_text, ccr.new_text, ccr.date, ccr.changeType,
                                      ct1.type as newType, ct1.allowHTML as newAllowHtml,
                                      ct2.type as oldType, ct2.allowHTML as oldAllowHtml,
                                      ccr.userID as userId
                                      from CustomerContent_Revisions ccr
                                      left join ContentType ct1 on ccr.new_type = ct1.cTypeId
                                      left join ContentType ct2 on ccr.old_type = ct2.cTypeId
                                      join CustomerContent cc on ccr.contentID = cc.id
                                      join Customer as c on cc.custID = c.cust_id
                                      join CustomerUser as cu on c.cust_id = cu.cust_ID
                                      join ApiKey as ak on cu.id = ak.user_id
                                      where ak.api_key = ? and ccr.contentID = ?
                                      order by ccr.date`
	getContentTypeIdStmt      = `select cTypeID, type, allowHTML from ContentType where type = ? limit 1`
	getAllContentTypesStmt    = `select type, allowHTML from ContentType order by type`
	updateCustomerContentStmt = `update CustomerContent as cc
                                join Customer as c on cc.custID = c.cust_id
                                join CustomerUser as cu on c.cust_id = cu.cust_ID
                                join ApiKey as ak on cu.id = ak.user_id
                                set cc.text = ?, cc.modified = now(),
                                cc.userID = cu.id, cc.typeID = ?, cc.deleted = ?
                                where ak.api_key = ? and cc.id = ?`
	insertCustomerContentStmt = `insert into CustomerContent (
                                    text, custID, modified, userID, typeID, deleted
                                )
                                select ?, c.cust_id, now(), cu.id, ?, 0
                                from Customer as c
                                join CustomerUser as cu on c.cust_id = cu.cust_ID
                                join ApiKey as ak on cu.id = ak.user_id
                                where ak.api_key = ?`
	checkExistingCustomerContentBridgeStmt = `select count(id) from CustomerContentBridge
                                             where partID = ? and catID = ? and contentID = ?`
	createCustomerContentBridgeStmt = `insert into CustomerContentBridge(partID, catID, contentID) values (?,?,?)`
	deleteCustomerContentBridgeStmt = `delete from CustomerContentBridge
                                      where contentID in(
                                            select cc.id from CustomerContent as cc
                                            join Customer as c on cc.custID = c.cust_id
                                            join CustomerUser as cu on c.cust_id = cu.cust_ID
                                            join ApiKey ak on cu.id = ak.user_id
                                            where api_key = ? and contentID = ?
                                      ) and partID = ? and catID = ?`
	markCustomerContentDeletedStmt = `update CustomerContent as cc
                                     join Customer as c on cc.custID = c.cust_id
                                     join CustomerUser as cu on c.cust_id = cu.cust_ID
                                     join ApiKey as ak on cu.id = ak.user_id
                                     set cc.deleted = 1, cc.modified = now(),
                                     cc.userID = cu.id where ak.api_key = ? and cc.id = ?`
)

type CustomerContent struct {
	Id              int
	Text            string
	Added, Modified time.Time
	ContentType     ContentType
	Hidden          bool
}

type IndexedContentType struct {
	Id        int
	Type      string
	AllowHtml bool
}

type ContentType struct {
	Type      string
	AllowHtml bool
}

type CustomerContentRevision struct {
	User                           customer.CustomerUser
	OldText, NewText               string
	Date                           time.Time
	ChangeType                     string
	OldContentType, NewContentType ContentType
}

// Retrieves all content for this customer
func AllCustomerContent(key string) (content []CustomerContent, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllCustomerContentStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(key)
	if err != nil {
		return
	}

	var partID, catID *int

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
			&partID,
			&catID,
		)

		if err != nil {
			return
		}

		if partID != nil && *partID > 0 {
			c.ContentType.Type = "Part: " + c.ContentType.Type
		} else if catID != nil && *catID > 0 {
			c.ContentType.Type = "Category: " + c.ContentType.Type
		}

		content = append(content, c)
	}

	return
}

//Copied from existing GoAPI - what? it grabs a random piece of content? Great.
func GetCustomerContent(id int, key string) (content CustomerContent, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getCustomerContentStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	var partID, catID *int

	err = stmt.QueryRow(key, id).Scan(
		&content.Id,
		&content.Text,
		&content.Added,
		&content.Modified,
		&content.Hidden,
		&content.ContentType.Type,
		&content.ContentType.AllowHtml,
		&partID,
		&catID,
	)
	if err != nil {
		return
	}

	if partID != nil && *partID > 0 {
		content.ContentType.Type = "Part: " + content.ContentType.Type
	} else if catID != nil && *catID > 0 {
		content.ContentType.Type = "Category: " + content.ContentType.Type
	}

	return
}

func GetCustomerContentRevisions(id int, key string) (revs []CustomerContentRevision, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getCustomerContentRevisionsStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(key, id)
	if err != nil {
		return
	}

	users := make(map[string]customer.CustomerUser, 0)

	for rows.Next() {
		var rev CustomerContentRevision
		err = rows.Scan(
			&rev.OldText,
			&rev.NewText,
			&rev.Date,
			&rev.ChangeType,
			&rev.NewContentType.Type,
			&rev.NewContentType.AllowHtml,
			&rev.OldContentType.Type,
			&rev.OldContentType.AllowHtml,
			&rev.User.Id,
		)
		if err != nil {
			return
		}

		if _, ok := users[rev.User.Id]; !ok {
			u, err := customer.GetCustomerUserById(rev.User.Id)
			if err == nil {
				users[rev.User.Id] = u
			}
		}
		rev.User = users[rev.User.Id]
		revs = append(revs, rev)
	}

	return
}

func (content *CustomerContent) Save(partID, catID int, key string) error {

	// If the Id is 0, we're adding a new
	// content piece; so we'll invoke that
	// method and return it's error
	if content.Id == 0 {
		return content.insert(partID, catID, key)
	}

	// Validate
	if content.Text == "" {
		return errors.New("Invalid content text: Content text was empty; if attempting to remove, use deletion endpoint.")
	}

	contentType, err := content.GetContentType()
	if err != nil || contentType.Id == 0 {
		return errors.New("Failed to match up to a content type.")
	}

	// We need to escape any possible HTML
	// if this content type doesn't allow
	// HTML
	if !contentType.AllowHtml {
		content.Text = html.EscapeString(content.Text)
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(updateCustomerContentStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(content.Text, key, content.Id, contentType.Id, content.Hidden)
	if err != nil {
		return err
	}

	// We need to bind this to a part or a category
	// just in case it was deleted at some point
	// and the customer is re-enabling it
	err = content.bridge(partID, catID)

	return err
}

func (content *CustomerContent) Delete(partID, catID int, key string) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	// First we need to delete the reference
	stmt, err := db.Prepare(deleteCustomerContentBridgeStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(key, content.Id, partID, catID)
	if err != nil {
		return errors.New("Failed to delete content bridge.")
	}

	// Mark the content piece as deleted
	stmt, err = db.Prepare(markCustomerContentDeletedStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(key, content.Id)
	if err != nil {
		return errors.New("Failed to mark content as deleted.")
	}

	content.Hidden = true

	return nil
}

func (content *CustomerContent) insert(partID, catID int, key string) error {
	contentType, err := content.GetContentType()
	if err != nil {
		return err
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertCustomerContentStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// We need to escape any possible HTML
	// if this content type doesn't allow
	// HTML
	if !contentType.AllowHtml {
		content.Text = html.EscapeString(content.Text)
	}

	res, err := stmt.Exec(content.Text, contentType.Id, key)
	if err != nil {
		return err
	}

	// Get the id of the record that was just inserted
	id, _ := res.LastInsertId()
	content.Id = int(id)

	// We need to bind this to a part or a category
	err = content.bridge(partID, catID)

	return err
}

func (content CustomerContent) bridge(partID, catID int) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	// Get the query to check if
	// there's already a record with this info
	stmt, err := db.Prepare(checkExistingCustomerContentBridgeStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the check
	res, err := stmt.Exec(partID, catID, content.Id)
	if err != nil || res == nil {
		// Either the query errored
		// or we already have a reference for this
		return err
	}

	// Create a bridge between the content
	// and the part/category
	stmt, err = db.Prepare(createCustomerContentBridgeStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(partID, catID, content.Id)

	return err
}

func (content CustomerContent) GetContentType() (ct IndexedContentType, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getContentTypeIdStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	cType := content.ContentType.Type

	typeArr := strings.Split(content.ContentType.Type, ":")
	if len(typeArr) > 1 {
		cType = typeArr[1]
	}

	err = stmt.QueryRow(cType).Scan(&ct.Id, &ct.Type, &ct.AllowHtml)
	if err == sql.ErrNoRows {
		err = nil
	}

	return
}

func AllCustomerContentTypes() (types []ContentType, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllContentTypesStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	res, err := stmt.Query()
	if err != nil {
		return
	}

	for res.Next() {
		var c ContentType
		err = res.Scan(&c.Type, &c.AllowHtml)
		if err != nil {
			return
		}
		types = append(types, c)
	}

	return
}
