package custcontent

import (
	"database/sql"
	"errors"
	"github.com/curt-labs/GoAPI/helpers/conversions"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/models/customer"
	_ "github.com/go-sql-driver/mysql"
	"html"

	// "strconv"
	"strings"
	"time"
)

type CustomerContent struct {
	Id          int
	Text        string
	Added       time.Time
	Modified    time.Time
	ContentType ContentType
	Hidden      bool
	Customer    *customer.Customer
	User        *customer.CustomerUser
	Revisions   CustomerContentRevisions
}

type ContentType struct {
	Id        int
	Type      string
	AllowHtml bool
}

type CustomerContentRevision struct {
	Id             int
	User           customer.CustomerUser
	Customer       customer.Customer
	OldText        string
	NewText        string
	Date           time.Time
	ChangeType     string
	ContentId      int
	OldContentType ContentType
	NewContentType ContentType
}
type CustomerContentRevisions []CustomerContentRevision

var (
	createContentType     = `insert into ContentType(type, allowHTML) values (?,?)`
	createContentRevision = `insert into Content_Revisions (userID, old_text, new_text, date, changeType, contentID, old_type, new_type) values (?,?,?,?,?,?,?,?)`
	deleteContentType     = `delete from ContentType where cTypeId = ?`
	deleteContentRevision = `delete from Content_Revisions where id = ?`

	allCustomerContent = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
							ct.type,ct.allowHTML,
							ccb.partID, ccb.catID,
							cc.userID, cc.custID
							from CustomerContent as cc
							left join CustomerContentBridge as ccb on cc.id = ccb.contentID
							join ContentType as ct on cc.typeID = ct.cTypeID
							join Customer as c on cc.custID = c.cust_id
							join CustomerUser as cu on c.cust_id = cu.cust_ID
							join ApiKey as ak on cu.id = ak.user_id
							where api_key = ?
							group by cc.id`
	customerContent = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
							ct.type,ct.allowHTML,ccb.partID,ccb.catID,
							cc.userID, cc.custID
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
									text, custID, added, modified, userID, typeID, deleted
								)
								select ?, c.cust_id, now(), now(), cu.id, ?, 0
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
	customerContentRevisions = `select ccr.id, ccr.old_text, ccr.new_text, ccr.date, ccr.changeType,
									ct1.type as newType, ct1.allowHTML as newAllowHtml,
									ct2.type as oldType, ct2.allowHTML as oldAllowHtml,
									ccr.userID as userId, ccr.custID
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

	getRevisionsByContentId = `select ccr.id, ccr.userID, ccr.custID, ccr.old_text, ccr.new_text, ccr.date, ccr.changeType, ccr.contentID, ccr.old_type, ccr.new_type
								from CustomerContent_Revisions as ccr
								where contentID = ?`
)

const (
	timeFormat     = "01/02/2006"
	timeYearFormat = "2006"
)

func (ct *ContentType) Create() error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(createContentType)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(ct.Type, ct.AllowHtml)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	ct.Id = int(id)
	return err
}
func (ccr *CustomerContentRevision) Create() error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(createContentRevision)
	if err != nil {
		return err
	}
	ccr.Date = time.Now()
	defer stmt.Close()
	res, err := stmt.Exec(ccr.User.Id, ccr.OldText, ccr.NewText, ccr.Date, ccr.ChangeType, ccr.ContentId, ccr.OldContentType.Id, ccr.NewContentType.Id)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	ccr.Id = int(id)
	return err
}

func (ct *ContentType) Delete() error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deleteContentType)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ct.Id)
	if err != nil {
		return err
	}
	return err
}

func (ccr *CustomerContentRevision) Delete() error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deleteContentRevision)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ccr.Id)
	if err != nil {
		return err
	}
	return err
}

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
	var pId, cId []byte
	var deleted *bool
	var added, modified *time.Time
	var contentType, userId *string
	var partId, catId int
	var custId *int
	var u customer.CustomerUser
	var cus customer.Customer
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
			&userId,
			&custId,
		)
		if err != nil {
			return content, err
		}
		if userId != nil {
			u.Id = *userId
			c.User = &u
		}
		if custId != nil {
			cus.Id = *custId
			c.Customer = &cus
		}
		if pId != nil {
			partId, err = conversions.ByteToInt(pId)
		}
		if cId != nil {
			catId, err = conversions.ByteToInt(cId)
		}
		if partId > 0 {
			c.ContentType.Type = "Part:" + *contentType
		} else if catId > 0 {
			c.ContentType.Type = "Category:" + *contentType
		} else {
			c.ContentType.Type = *contentType
		}
		if deleted != nil {
			c.Hidden = *deleted
		}
		if modified != nil {
			c.Modified = *modified
		}
		if added != nil {
			c.Added = *added
		}
		err = c.GetRevisions()
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

	stmt, err := db.Prepare(customerContent)
	if err != nil {
		return c, err
	}

	rows, err := stmt.Query(key, id)
	if err != nil {
		return c, err
	}

	var pId, cId []byte
	var deleted *bool
	var added, modified *time.Time
	var contentType string
	var partId, catId int
	var custId *int
	var userId *string
	var u customer.CustomerUser
	var cus customer.Customer

	for rows.Next() {
		err = rows.Scan(
			&c.Id,
			&c.Text,
			&added,
			&modified,
			&deleted,
			&contentType,
			&c.ContentType.AllowHtml,
			&pId,
			&cId,
			&userId,
			&custId,
		)
		if err != nil {
			return c, err
		}
	}

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
	if deleted != nil {
		c.Hidden = *deleted
	}
	if modified != nil {
		c.Modified = *modified
	}
	if added != nil {
		c.Added = *added
	}
	if userId != nil {
		u.Id = *userId
		c.User = &u
	}
	if custId != nil {
		cus.Id = *custId
		c.Customer = &cus
	}
	err = c.GetRevisions()
	return c, err
}

// by customer ID
func (cc *CustomerContent) GetRevisions() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getRevisionsByContentId)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Query(cc.Id)
	var r CustomerContentRevision
	var oldT, newT *string
	var oct, nct *int

	for res.Next() {
		err = res.Scan(
			&r.Id,
			&r.User.Id,
			&r.Customer.Id,
			&oldT,
			&newT,
			&r.Date,
			&r.ChangeType,
			&r.ContentId,
			&oct,
			&nct,
		)
		if err != nil {
			return err
		}
		if oldT != nil {
			r.OldText = *oldT
		}
		if newT != nil {
			r.NewText = *newT
		}
		if oct != nil {
			r.OldContentType.Id = *oct
		}
		if nct != nil {
			r.NewContentType.Id = *nct
		}
		cc.Revisions = append(cc.Revisions, r)
	}
	return err
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

	var ot, nt, octt *string
	var octa *bool

	for res.Next() {
		var ccr CustomerContentRevision
		err = res.Scan(
			&ccr.Id,
			&ot,
			&nt,
			&ccr.Date,
			&ccr.ChangeType,
			&ccr.NewContentType.Type,
			&ccr.NewContentType.AllowHtml,
			&octt,
			&octa,
			&ccr.User.Id,
			&ccr.Customer.Id,
		)
		if err != nil {
			return revs, err
		}
		if ot != nil {
			ccr.OldText = *ot
		}
		if nt != nil {
			ccr.NewText = *nt
		}
		if octt != nil {
			ccr.OldContentType.Type = *octt
		}
		if octa != nil {
			ccr.OldContentType.AllowHtml = *octa
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

	if content.ContentType.Id == 0 && content.ContentType.Type == "" {
		return errors.New("content type must be provided")
	} else {
		if err := content.GetContentType(); err != nil {
			return errors.New("faield to retrieve content type")
		}
	}

	// Validate
	if content.Text == "" {
		return errors.New("Invalid content text: Content text was empty; if attempting to remove, use deletion endpoint.")
	}

	// If the Id is 0, we're adding a new
	// content piece; so we'll invoke that
	// method and return it's error
	if content.Id == 0 {
		return content.insert(partID, catID, key)
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

	_, err = stmt.Exec(content.Text, content.ContentType.Id, hidden, key, content.Id) //TODO this right?
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
		return errors.New("Content Bridge already deleted.")
	}
	_, err = stmt.Exec(key, content.Id, partID, catID)
	if err != nil {
		tx.Rollback()
		return errors.New("Failed to delete content bridge.")
	}

	stmt, err = tx.Prepare(markCustomerContentDeleted)
	if err != nil {
		return errors.New("Content already deleted.")
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

	err = content.GetContentType()
	if err != nil {
		return errors.New("Error getting content type.")
	}

	if !content.ContentType.AllowHtml {
		content.Text = html.EscapeString(content.Text)
	}

	res, err := stmt.Exec(content.Text, content.ContentType.Id, key)

	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	id, err := res.LastInsertId()
	content.Id = int(id)

	err = content.bridge(partID, catID)
	return err
}

func (content *CustomerContent) bridge(partID, catID int) error {
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
	if count != 0 {
		return err
	}

	tx, err := db.Begin()

	stmt, err = tx.Prepare(createCustomerContentBridge)
	_, err = stmt.Exec(partID, catID, content.Id)

	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

//gets content by name
func (content *CustomerContent) GetContentType() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getContentTypeId)
	if err != nil {
		return
	}
	cType := content.ContentType.Type

	typeArr := strings.Split(content.ContentType.Type, ":")
	if len(typeArr) > 1 {
		cType = typeArr[1]
	}

	err = stmt.QueryRow(cType).Scan(&content.ContentType.Id, &content.ContentType.Type, &content.ContentType.AllowHtml)
	if err == sql.ErrNoRows {
		err = nil
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
	defer stmt.Close()

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
