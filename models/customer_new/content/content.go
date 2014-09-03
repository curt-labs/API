package custcontent

import (
	"database/sql"
	"errors"
	"github.com/curt-labs/GoAPI/helpers/conversions"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/models/customer"
	_ "github.com/go-sql-driver/mysql"
	"html"
	// "log"
	// "strconv"
	"strings"
	"time"
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

var (
	allCustomerContent = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
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
	customerContent = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
							ct.type,ct.allowHTML,ccb.partID,ccb.catID
							from CustomerContent as cc
							join CustomerContentBridge as ccb on cc.id = ccb.contentID
							join ContentType as ct on cc.typeID = ct.cTypeID
							join Customer as c on cc.custID = c.cust_id
							join CustomerUser as cu on c.cust_id = cu.cust_ID
							join ApiKey as ak on cu.id = ak.user_id
							where api_key = ? and cc.id = ?
							limit 1`
	updateCustomerContent = `update CustomerContent as cc
								join Customer as c on cc.custID = c.cust_id
								join CustomerUser as cu on c.cust_id = cu.cust_ID
								join ApiKey as ak on cu.id = ak.user_id
								set cc.text = ?, cc.modified = now(),
								cc.userID = cu.id, cc.typeID = ?, cc.deleted = ?
								where ak.api_key = ? and cc.id = ?`
	insertCustomerContent = `insert into CustomerContent (
									text, custID, modified, userID, typeID, deleted
								)
								select ?, c.cust_id, now(), cu.id, ?, 0
								from Customer as c
								join CustomerUser as cu on c.cust_id = cu.cust_ID
								join ApiKey as ak on cu.id = ak.user_id
								where ak.api_key = ?`
	checkExistingCustomerContentBridge = `select count(id) from CustomerContentBridge
												where partID = ? and catID = ? and contentID = ?`

	createCustomerContentBridge = `insert into CustomerContentBridge
										(partID, catID, contentID)
										values (?,?,?)`
	getContentTypeId         = `select cTypeID, type, allowHTML from ContentType where type = ? limit 1`
	getAllContentTypes       = `select type, allowHTML from ContentType order by type`
	customerContentRevisions = `select ccr.old_text, ccr.new_text, ccr.date, ccr.changeType,
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

	deleteCustomerContentBridge = `delete from CustomerContentBridge
									where contentID in(
										select cc.id from CustomerContent as cc
										join Customer as c on cc.custID = c.cust_id
										join CustomerUser as cu on c.cust_id = cu.cust_ID
										join ApiKey ak on cu.id = ak.user_id
										where api_key = ? and contentID = ?
									) and partID = ? and catID = ?`

	markCustomerContentDeleted = `update CustomerContent as cc
									join Customer as c on cc.custID = c.cust_id
									join CustomerUser as cu on c.cust_id = cu.cust_ID
									join ApiKey as ak on cu.id = ak.user_id
									set cc.deleted = 1, cc.modified = now(),
									cc.userID = cu.id where ak.api_key = ?
									and cc.id = ?`
)

const (
	timeFormat     = "01/02/2006"
	timeYearFormat = "2006"
)

// Retrieves all content for this customer
func AllCustomerContent(key string) (content []CustomerContent, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return content, err
	}
	defer db.Close()

	stmt, err := db.Prepare(allCustomerContent)
	if err != nil {
		return content, err
	}
	var deleted, added, pId, cId, modified []byte
	var contentType string
	var partId, catId int
	res, err := stmt.Query(key)
	for res.Next() {
		var c CustomerContent
		err = res.Scan(
			&c.Id,
			&c.Text,
			&added,
			&modified,
			&deleted,
			&contentType,
			&c.ContentType.AllowHtml,
			&pId,
			&cId,
		)
		if pId != nil {
			partId, err = conversions.ByteToInt(pId)
		}
		if cId != nil {
			catId, err = conversions.ByteToInt(cId)
		}
		if partId > 0 {
			c.ContentType.Type = "Part:" + contentType
		} else if catId > 0 {
			c.ContentType.Type = "Category:" + contentType
		} else {
			c.ContentType.Type = contentType
		}
		if modified != nil {
			m, err := conversions.ByteToString(modified)
			c.Modified, err = time.Parse("2006", m)
			if err != nil {
				return content, err
			}
		}
		if added != nil {
			a, err := conversions.ByteToString(added)
			c.Added, err = time.Parse(timeYearFormat, a)
			if err != nil {
				return content, err
			}
		}

		if err != nil {
			return content, err
		}

		content = append(content, c)
	}
	return
}

func GetCustomerContent(id int, key string) (c CustomerContent, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return c, err
	}
	defer db.Close()

	stmt, err := db.Prepare(allCustomerContent)
	if err != nil {
		return c, err
	}
	var deleted, added, modified, pId, cId []byte
	var contentType string
	var partId, catId int
	err = stmt.QueryRow(key).Scan(
		&c.Id,
		&c.Text,
		&added,
		&modified,
		&deleted,
		&contentType,
		&c.ContentType.AllowHtml,
		&pId,
		&cId,
	)
	if pId != nil {
		partId, err = conversions.ByteToInt(pId)
	}
	if cId != nil {
		catId, err = conversions.ByteToInt(cId)
	}
	if partId > 0 {
		c.ContentType.Type = "Part:" + contentType
	} else if catId > 0 {
		c.ContentType.Type = "Category:" + contentType
	} else {
		c.ContentType.Type = contentType
	}
	if modified != nil {
		m, err := conversions.ByteToString(modified)
		c.Modified, err = time.Parse(timeYearFormat, m)
		if err != nil {
			return c, err
		}
	}
	if added != nil {
		a, err := conversions.ByteToString(added)
		c.Added, err = time.Parse(timeYearFormat, a)
		if err != nil {
			return c, err
		}
	}
	return c, err
}

func GetCustomerContentRevisions(id int, key string) (revs []CustomerContentRevision, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return revs, err
	}
	defer db.Close()

	stmt, err := db.Prepare(customerContentRevisions)
	if err != nil {
		return revs, err
	}
	res, err := stmt.Query(key, id)

	users := make(map[string]customer.CustomerUser, 0)

	for res.Next() {
		var ccr CustomerContentRevision
		err = res.Scan(
			&ccr.OldText,
			&ccr.NewText,
			&ccr.Date,
			&ccr.ChangeType,
			&ccr.NewContentType.Type,
			&ccr.NewContentType.AllowHtml,
			&ccr.OldContentType.Type,
			&ccr.OldContentType.AllowHtml,
			&ccr.User.Id,
		)
		if err != nil {
			return revs, err
		}

		if _, ok := users[ccr.User.Id]; !ok {
			u, err := customer.GetCustomerUserById(ccr.User.Id)
			if err == nil {
				users[ccr.User.Id] = u
			}
		}
		ccr.User = users[ccr.User.Id]
		revs = append(revs, ccr)
	}
	return
}

func (content *CustomerContent) Save(partID, catID int, key string) error { //TODO - I would determine create/update in the controller
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

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()

	stmt, err := tx.Prepare(updateCustomerContent)
	if err != nil {
		return err
	}
	hidden := 0
	if content.Hidden {
		hidden = 1
	}
	_, err = stmt.Exec(content.Text, key, content.Id, contentType.Id, hidden) //TODO this right?
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

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
	tx, err := db.Begin()

	stmt, err := tx.Prepare(deleteCustomerContentBridge)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(key, content.Id, partID, catID)
	if err != nil {
		tx.Rollback()
		return errors.New("Failed to delete content bridge.")
	}
	tx.Commit()

	stmt, err = tx.Prepare(markCustomerContentDeleted)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(key, content.Id)
	if err != nil {
		tx.Rollback()
		return errors.New("Failed to mark content as deleted.")
	}
	tx.Commit()
	content.Hidden = true

	return nil
}

func (content *CustomerContent) insert(partID, catID int, key string) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()

	stmt, err := tx.Prepare(insertCustomerContent)
	if err != nil {
		return err
	}

	contentType, err := content.GetContentType()
	if err != nil {
		return errors.New("Error getting content type.")
	}

	if !contentType.AllowHtml {
		content.Text = html.EscapeString(content.Text)
	}

	res, err := stmt.Exec(content.Text, contentType.Id, key)
	if err != nil {
		tx.Rollback()
		return errors.New("Error executing statement.")
	}

	id, err := res.LastInsertId()
	content.Id = int(id)

	err = content.bridge(partID, catID)
	return err
}

func (content CustomerContent) bridge(partID, catID int) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(checkExistingCustomerContentBridge)
	if err != nil {
		return err
	}

	var count int
	err = stmt.QueryRow(partID, catID, content.Id).Scan(&count)

	tx, err := db.Begin()

	stmt, err = tx.Prepare(createCustomerContentBridge)
	_, err = stmt.Exec(partID, catID, content.Id)

	if err != nil {
		tx.Rollback()
		return err
	}
	return err
}

func (content CustomerContent) GetContentType() (ct IndexedContentType, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ct, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getContentTypeId)
	if err != nil {
		return ct, err
	}
	cType := content.ContentType.Type

	typeArr := strings.Split(content.ContentType.Type, ":")
	if len(typeArr) > 1 {
		cType = typeArr[1]
	}

	err = stmt.QueryRow(cType).Scan(&ct.Id, &ct.Type, &ct.AllowHtml)
	if err != nil {
		return ct, err
	}
	return
}

func AllCustomerContentTypes() (types []ContentType, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return types, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllContentTypes)
	if err != nil {
		return types, err
	}
	res, err := stmt.Query()
	for res.Next() {
		var ct ContentType
		err = res.Scan(&ct.Type, &ct.AllowHtml)
		if err != nil {
			return types, err
		}
		types = append(types, ct)
	}
	return
}
