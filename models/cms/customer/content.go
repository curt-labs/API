package customer_cms

import (
	"../../../helpers/database"
	"time"
)

type CustomerContent struct {
	Id              int
	Text            string
	Added, Modified time.Time
	ContentType     ContentType
	Hidden          bool
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
