package models

import (
	"../helpers/database"
	"net/url"
	"strconv"
	"strings"
)

type Image struct {
	Size, Sort    string
	Height, Width int
	Path          *url.URL
}

var (
	partImageStmt = `select pis.size,pi.sort,pi.height,pi.width,pi.path from PartImages as pi
				join PartImageSizes as pis on pi.sizeID = pis.sizeID
				where partID = ? order by pi.sort, pi.height`

	partImageStmt_ByGroup = `select partID, pis.size,pi.sort,pi.height,pi.width,pi.path from PartImages as pi
				join PartImageSizes as pis on pi.sizeID = pis.sizeID
				where partID IN (%s) order by pi.sort, pi.height`
)

func (p *Part) GetImages() error {

	qry, err := database.Db.Prepare(partImageStmt)
	if err != nil {
		return err
	}

	rows, res, err := qry.Exec(p.PartId)
	if database.MysqlError(err) {
		return err
	}

	size := res.Map("size")
	sort := res.Map("sort")
	height := res.Map("height")
	width := res.Map("width")
	path := res.Map("path")

	var images []Image
	for _, row := range rows {
		imgPath, urlErr := url.Parse(row.Str(path))
		if urlErr == nil {
			img := Image{
				Size:   row.Str(size),
				Sort:   row.Str(sort),
				Height: row.Int(height),
				Width:  row.Int(width),
				Path:   imgPath,
			}

			images = append(images, img)
		}
	}

	p.Images = images

	return nil
}

func GetImagesByGroup(existing map[int]Part) (parts map[int]Part, err error) {

	parts = make(map[int]Part, len(existing))
	var ids []string
	for k, _ := range existing {
		parts[k] = Part{PartId: k}
		ids = append(ids, strconv.Itoa(k))
	}

	rows, res, err := database.Db.Query(partImageStmt_ByGroup, strings.Join(ids, ","))
	if database.MysqlError(err) {
		return
	}

	partID := res.Map("partID")
	size := res.Map("size")
	sort := res.Map("sort")
	height := res.Map("height")
	width := res.Map("width")
	path := res.Map("path")

	for _, row := range rows {
		imgPath, urlErr := url.Parse(row.Str(path))
		if urlErr == nil {
			pId := row.Int(partID)
			img := Image{
				Size:   row.Str(size),
				Sort:   row.Str(sort),
				Height: row.Int(height),
				Width:  row.Int(width),
				Path:   imgPath,
			}

			tmp := parts[pId]
			tmp.Images = append(tmp.Images, img)
			parts[pId] = tmp
		}
	}

	return

}
