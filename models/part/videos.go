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

type PartVideo struct {
	YouTubeVideoId, Type string
	IsPrimary            bool
	TypeIcon             *url.URL
}

var (
	partVideoStmt = `select pv.video,vt.name,pv.isPrimary, vt.icon from PartVideo as pv
				join videoType vt on pv.vTypeID = vt.vTypeID
				where pv.partID = ?`
)

func (p *Part) GetVideos() error {
	redis_key := fmt.Sprintf("part:%d:videos", p.PartId)

	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		if err = json.Unmarshal(data, &p.Videos); err != nil {
			return nil
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	qry, err := db.Prepare(partVideoStmt)
	if err != nil {
		return err
	}
	defer qry.Close()

	rows, err := qry.Query(p.PartId)
	if err != nil {
		return err
	}

	var videos []PartVideo
	for rows.Next() {
		var v PartVideo
		var icon *string
		err = rows.Scan(
			&v.YouTubeVideoId,
			&v.Type,
			&v.IsPrimary,
			&icon)
		if err != nil {
			continue
		}

		v.TypeIcon, _ = url.Parse(*icon)
		videos = append(videos, v)
	}

	p.Videos = videos

	go redis.Setex(redis_key, p.Videos, redis.CacheTimeout)

	return nil
}
