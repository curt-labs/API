package models

import (
	"github.com/curt-labs/GoAPI/helpers/database"
	"strconv"
	"strings"
	"time"
)

var (
	CustomerPartContent_Grouped = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted, 
									ct.type,ct.allowHTML,ccb.partID
									from CustomerContent as cc
									join CustomerContentBridge as ccb on cc.id = ccb.contentID
									join ContentType as ct on cc.typeID = ct.cTypeID
									join Customer as c on cc.custID = c.cust_id
									join CustomerUser as cu on c.cust_id = cu.cust_ID
									join ApiKey as ak on cu.id = ak.user_id
									where api_key = '%s' and ccb.partID IN (%s)
									group by cc.id, ccb.partID`
)

type PartContent struct {
	PartId  int
	Content []CustomerContent
}

// Retrieves all part content for this customer
func GetAllPartContent(key string) (content []PartContent, err error) {

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

func GetGroupedPartContent(ids []string, key string) (content map[int][]CustomerContent, err error) {

	content = make(map[int][]CustomerContent, len(ids))

	for i := 0; i < len(ids); i++ {
		intId, err := strconv.Atoi(ids[i])
		if err == nil {
			content[intId] = make([]CustomerContent, 0)
		}
	}

	escaped_key := database.Db.Escape(key)

	rows, res, err := database.Db.Query(CustomerPartContent_Grouped, escaped_key, strings.Join(ids, ","))
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
		content[row.Int(partID)] = append(content[row.Int(partID)], c)
	}

	return

}
