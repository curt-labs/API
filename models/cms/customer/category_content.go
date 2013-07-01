package customer_cms

import (
	"../../../helpers/database"
	"time"
)

type CategoryContent struct {
	CategoryId int
	Content    []CustomerContent
}

// Retrieves all category content for this customer
func AllCategoryContent(key string) (content []CategoryContent, err error) {
	qry, err := database.GetStatement("AllCustomerCategoryContent")
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
	catID := res.Map("catID")

	rawContent := make(map[int][]CustomerContent, 0)

	for _, row := range rows {
		c := CustomerContent{
			Id:       row.Int(id),
			Text:     row.Str(text),
			Added:    row.ForceTime(added, time.UTC),
			Modified: row.ForceTime(mod, time.UTC),
			Hidden:   row.ForceBool(deleted),
			ContentType: ContentType{
				Type:      "Category:" + row.Str(cType),
				AllowHtml: row.ForceBool(html),
			},
		}

		cat_id := row.Int(catID)
		if cat_id > 0 {
			rawContent[cat_id] = append(rawContent[cat_id], c)
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

	qry, err := database.GetStatement("CustomerCategoryContent")
	if database.MysqlError(err) {
		return
	}

	rows, res, err := qry.Exec(key, catID)
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

	for _, row := range rows {
		c := CustomerContent{
			Id:       row.Int(id),
			Text:     row.Str(text),
			Added:    row.ForceTime(added, time.UTC),
			Modified: row.ForceTime(mod, time.UTC),
			Hidden:   row.ForceBool(deleted),
			ContentType: ContentType{
				Type:      "Category:" + row.Str(cType),
				AllowHtml: row.ForceBool(html),
			},
		}
		content = append(content, c)
	}

	return
}
