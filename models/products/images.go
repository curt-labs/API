package products

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
)

var (
	partImageStmt = `select pis.size,pi.sort,pi.height,pi.width,pi.path from PartImages as pi
				join PartImageSizes as pis on pi.sizeID = pis.sizeID
				where partID = ? order by pi.sort, pi.height`
	createPartImage  = `INSERT INTO PartImages (sizeID, sort, path, height, width, partID) VALUES (?,?,?,?,?,?)`
	deletePartImages = `DELETE FROM PartImages WHERE partID = ?`
)

type Image struct {
	ID     int      `json:"id,omitempty" xml:"id,omitempty"`
	Size   string   `json:"size,omitempty" xml:"size,omitempty"`
	Sort   string   `json:"sort,omitempty" xml:"sort,omitempty"`
	Height int      `json:"height,omitempty" xml:"height,omitempty"`
	Width  int      `json:"width,omitempty" xml:"width,omitempty"`
	Path   *url.URL `json:"path,omitempty" xml:"path,omitempty"`
	PartID int      `json:"partId,omitempty" xml:"partId,omitempty"`
}

func (p *Part) GetImages(dtx *apicontext.DataContext) error {

	redis_key := fmt.Sprintf("part:%d:images:%s", p.ID, dtx.BrandArray)

	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		if err = json.Unmarshal(data, &p.Images); err != nil {
			return nil
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	qry, err := db.Prepare(partImageStmt)
	if err != nil {
		return err
	}
	defer qry.Close()

	rows, err := qry.Query(p.ID)
	if err != nil {
		return err
	}

	var images []Image
	for rows.Next() {
		var img Image
		var path *string
		err = rows.Scan(
			&img.Size,
			&img.Sort,
			&img.Height,
			&img.Width,
			&path)
		if err == nil && path != nil {
			img.Path, err = url.Parse(*path)
			if err == nil {
				images = append(images, img)
			}
		}
	}
	defer rows.Close()

	p.Images = images
	go redis.Setex(redis_key, p.Images, redis.CacheTimeout)

	return nil
}

func (i *Image) Create(dtx *apicontext.DataContext) (err error) {
	go redis.Delete(fmt.Sprintf("part:%d:images:"+dtx.BrandString, i.PartID))
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(createPartImage)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(i.Size, i.Sort, i.Path, i.Height, i.Width, i.PartID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	i.ID = int(id)
	return nil
}

func (i *Image) DeleteByPart(dtx *apicontext.DataContext) (err error) {
	go redis.Delete(fmt.Sprintf("part:%d:images:"+dtx.BrandString, i.PartID))
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deletePartImages)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(i.PartID)
	if err != nil {
		return err
	}
	return nil
}
