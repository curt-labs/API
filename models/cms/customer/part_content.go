package customer_cms

import (
	"../../../helpers/database"
	"time"
)

type PartContent struct {
	PartId  int
	Content []CustomerContent
}

// Retrieves all part content for this customer
func AllPartContent(key string) (content []PartContent, err error) {

	qry, err := database.GetStatement("AllCustomerPartContent")
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

	rawContent := make(map[int][]CustomerContent, 0)

	for _, row := range rows {
		c := CustomerContent{
			Id:       row.Int(id),
			Text:     row.Str(text),
			Added:    row.ForceTime(added, time.UTC),
			Modified: row.ForceTime(mod, time.UTC),
			Hidden:   row.ForceBool(deleted),
			ContentType: ContentType{
				Type:      "Part:" + row.Str(cType),
				AllowHtml: row.ForceBool(html),
			},
		}

		part_id := row.Int(partID)
		if part_id > 0 {
			rawContent[part_id] = append(rawContent[part_id], c)
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

	qry, err := database.GetStatement("CustomerPartContent")
	if database.MysqlError(err) {
		return
	}

	rows, res, err := qry.Exec(key, partID)
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
				Type:      "Part:" + row.Str(cType),
				AllowHtml: row.ForceBool(html),
			},
		}
		content = append(content, c)
	}

	return
}
