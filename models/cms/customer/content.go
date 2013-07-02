package customer_cms

import (
	"../../../helpers/database"
	"errors"
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
	User                           string
	OldText, NewText               string
	Date                           time.Time
	ChangeType                     string
	OldContentType, NewContentType ContentType
}

// Retrieves all content for this customer
func AllContent(key string) (content []CustomerContent, err error) {

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
