package part

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
)

var (
	partImageStmt = `select pis.size,pi.sort,pi.height,pi.width,pi.path from PartImages as pi
				join PartImageSizes as pis on pi.sizeID = pis.sizeID
				where partID = ? order by pi.sort, pi.height`
)

type Image struct {
	Size, Sort    string
	Height, Width int
	Path          *url.URL
}

func (p *Part) GetImages() error {
	redis_key := fmt.Sprintf("part:%d:images", p.PartId)

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

	rows, err := qry.Query(p.PartId)
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

	p.Images = images

	go redis.Setex(redis_key, p.Images, redis.CacheTimeout)

	return nil
}
