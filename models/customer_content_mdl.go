package models

import (
	"errors"
	"github.com/curt-labs/GoAPI/helpers/database"
	"html"
	"log"
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
	User                           CustomerUser
	OldText, NewText               string
	Date                           time.Time
	ChangeType                     string
	OldContentType, NewContentType ContentType
}

// Retrieves all content for this customer
func AllCustomerContent(key string) (content []CustomerContent, err error) {

	qry, err := database.GetStatement("AllCustomerContent")
	if database.MysqlError(err) {
		return
	}

	rows, res, err := qry.Exec(key)
	if database.MysqlError(err) {
		return
	}

	id := res.Map("id")
	text := res.Map("text")
	added := res.Map("added")
	mod := res.Map("modified")
	deleted := res.Map("deleted")
	cType := res.Map("type")
	html := res.Map("allowHTML")
	partID := res.Map("partID")
	catID := res.Map("catID")

	for _, row := range rows {
		c := CustomerContent{
			Id:       row.Int(id),
			Text:     row.Str(text),
			Added:    row.ForceTime(added, time.UTC),
			Modified: row.ForceTime(mod, time.UTC),
			Hidden:   row.ForceBool(deleted),
			ContentType: ContentType{
				AllowHtml: row.ForceBool(html),
			},
		}

		part_id := row.Int(partID)
		cat_id := row.Int(catID)
		if part_id > 0 {
			c.ContentType.Type = "Part:" + row.Str(cType)
		} else if cat_id > 0 {
			c.ContentType.Type = "Category:" + row.Str(cType)
		} else {
			c.ContentType.Type = row.Str(cType)
		}

		content = append(content, c)
	}

	return
}

func GetCustomerContent(id int, key string) (content CustomerContent, err error) {
	qry, err := database.GetStatement("CustomerContent")
	if database.MysqlError(err) {
		return
	}

	row, res, err := qry.ExecFirst(key, id)
	if database.MysqlError(err) {
		return
	}

	text := res.Map("text")
	added := res.Map("added")
	mod := res.Map("modified")
	deleted := res.Map("deleted")
	cType := res.Map("type")
	html := res.Map("allowHTML")
	partID := res.Map("partID")
	catID := res.Map("catID")

	content = CustomerContent{
		Id:       id,
		Text:     row.Str(text),
		Added:    row.ForceTime(added, time.UTC),
		Modified: row.ForceTime(mod, time.UTC),
		Hidden:   row.ForceBool(deleted),
		ContentType: ContentType{
			AllowHtml: row.ForceBool(html),
		},
	}

	part_id := row.Int(partID)
	cat_id := row.Int(catID)
	if part_id > 0 {
		content.ContentType.Type = "Part:" + row.Str(cType)
	} else if cat_id > 0 {
		content.ContentType.Type = "Category:" + row.Str(cType)
	} else {
		content.ContentType.Type = row.Str(cType)
	}

	return
}

func GetCustomerContentRevisions(id int, key string) (revs []CustomerContentRevision, err error) {
	qry, err := database.GetStatement("CustomerContentRevisions")
	if database.MysqlError(err) {
		return
	}

	rows, res, err := qry.Exec(key, id)
	if database.MysqlError(err) {
		return
	}

	txtOld := res.Map("old_text")
	txtNew := res.Map("new_text")
	date := res.Map("date")
	change := res.Map("changeType")
	newType := res.Map("newType")
	oldType := res.Map("oldType")
	newHTML := res.Map("newAllowHtml")
	oldHTML := res.Map("oldAllowHtml")
	userId := res.Map("userId")

	users := make(map[string]CustomerUser, 0)

	for _, row := range rows {
		ccr := CustomerContentRevision{
			OldText:    row.Str(txtOld),
			NewText:    row.Str(txtNew),
			Date:       row.ForceTime(date, time.UTC),
			ChangeType: row.Str(change),
			OldContentType: ContentType{
				Type:      row.Str(oldType),
				AllowHtml: row.ForceBool(oldHTML),
			},
			NewContentType: ContentType{
				Type:      row.Str(newType),
				AllowHtml: row.ForceBool(newHTML),
			},
		}

		if _, ok := users[row.Str(userId)]; !ok {
			u, err := GetCustomerUserById(row.Str(userId))
			if err == nil {
				users[row.Str(userId)] = u
			}
		}
		ccr.User = users[row.Str(userId)]

		revs = append(revs, ccr)
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

	qry, err := database.GetStatement("UpdateCustomerContent")
	if database.MysqlError(err) {
		return err
	}

	// We need to escape any possible HTML
	// if this content type doesn't allow
	// HTML
	if !contentType.AllowHtml {
		content.Text = html.EscapeString(content.Text)
	}

	hidden := 0
	if content.Hidden {
		hidden = 1
	}

	qry.Bind(content.Text, key, content.Id, contentType.Id, hidden)
	_, _, err = qry.Exec()
	if database.MysqlError(err) {
		return err
	}

	// We need to bind this to a part or a category
	// just in case it was deleted at some point
	// and the customer is re-enabling it
	err = content.bridge(partID, catID)

	return err
}

func (content *CustomerContent) Delete(partID, catID int, key string) error {

	// First we need to delete the reference
	qry, err := database.GetStatement("DeleteCustomerContentBridge")
	if database.MysqlError(err) {
		return err
	}

	_, _, err = qry.Exec(key, content.Id, partID, catID)
	if database.MysqlError(err) {
		return errors.New("Failed to delete content bridge.")
	}

	// Mark the content piece as deleted
	qry, err = database.GetStatement("MarkCustomerContentDeleted")
	if database.MysqlError(err) {
		return err
	}

	_, _, err = qry.Exec(key, content.Id)
	if database.MysqlError(err) {
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

	qry, err := database.GetStatement("InsertCustomerContent")
	if database.MysqlError(err) {
		return err
	}

	// We need to escape any possible HTML
	// if this content type doesn't allow
	// HTML
	if !contentType.AllowHtml {
		content.Text = html.EscapeString(content.Text)
	}

	log.Println(contentType)

	qry.Bind(content.Text, contentType.Id, key)
	_, res, err := qry.Exec()
	if database.MysqlError(err) {
		return err
	}

	// Get the id of the record that was just inserted
	content.Id = int(res.InsertId())

	// We need to bind this to a part or a category
	err = content.bridge(partID, catID)

	return err
}

func (content CustomerContent) bridge(partID, catID int) error {

	// Get the query to check if
	// there's already a record with this info
	qry, err := database.GetStatement("CheckExistingCustomerContentBridge")
	if database.MysqlError(err) {
		return err
	}

	// Execute the check
	row, _, err := qry.ExecFirst(partID, catID, content.Id)
	if err != nil || row.Int(0) > 0 {
		// Either the query errored
		// or we already have a reference for this
		return err
	}

	// Create a bridge between the content
	// and the part/category
	qry, err = database.GetStatement("CreateCustomerContentBridge")
	if database.MysqlError(err) {
		return err
	}

	qry.Bind(partID, catID, content.Id)
	_, _, err = qry.Exec()
	if database.MysqlError(err) {
		return err
	}

	return err

}

func (content CustomerContent) GetContentType() (ct IndexedContentType, err error) {

	qry, err := database.GetStatement("GetContentTypeId")
	if database.MysqlError(err) {
		return
	}

	cType := content.ContentType.Type

	typeArr := strings.Split(content.ContentType.Type, ":")
	if len(typeArr) > 1 {
		cType = typeArr[1]
	}

	row, _, err := qry.ExecFirst(cType)
	if database.MysqlError(err) || row == nil {
		return
	}

	ct = IndexedContentType{
		Id:        row.Int(0),
		Type:      row.Str(1),
		AllowHtml: row.ForceBool(2),
	}

	return
}

func AllCustomerContentTypes() (types []ContentType, err error) {
	qry, err := database.GetStatement("GetAllContentTypes")
	if database.MysqlError(err) {
		return
	}

	rows, res, err := qry.Exec()
	if database.MysqlError(err) || rows == nil {
		return
	}

	typ := res.Map("type")
	html := res.Map("allowHTML")

	for _, row := range rows {
		ct := ContentType{
			Type:      row.Str(typ),
			AllowHtml: row.ForceBool(html),
		}
		types = append(types, ct)
	}
	return
}
