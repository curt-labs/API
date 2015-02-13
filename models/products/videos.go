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

type PartVideo struct {
	ID             int
	PartID         int
	YouTubeVideoId string
	Type           string
	IsPrimary      bool
	TypeIcon       *url.URL
	VideoType      VideoType
}

type VideoType struct {
	ID   int
	Name string
	Icon string
}

var (
	partVideoStmt = `select pv.video,vt.name,pv.isPrimary, vt.icon from PartVideo as pv
				join videoType vt on pv.vTypeID = vt.vTypeID
				where pv.partID = ?`
	createPartVideo  = `INSERT INTO PartVideo (partID, video, vTypeID, isPrimary) VALUES (?,?,?,?)`
	deletePartVideos = `DELETE FROM PartVideo WHERE partID = ?`
)

func (p *Part) GetVideos() error {
	redis_key := fmt.Sprintf("part:%d:%d:videos", p.BrandID, p.ID)

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

	rows, err := qry.Query(p.ID)
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
	defer rows.Close()

	go redis.Setex(redis_key, p.Videos, redis.CacheTimeout)

	return nil
}

func (p *PartVideo) CreatePartVideo(dtx *apicontext.DataContext) (err error) {
	go redis.Delete(fmt.Sprintf("part:%s:%d:videos", dtx.BrandString, p.PartID))
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(createPartVideo)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(p.PartID, p.YouTubeVideoId, p.VideoType.ID, p.IsPrimary)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	p.ID = int(id)
	return nil
}

func (p *PartVideo) DeleteByPart(dtx *apicontext.DataContext) (err error) {
	go redis.Delete(fmt.Sprintf("part:%s:%d:videos", dtx.BrandString, p.PartID))
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deletePartVideos)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.PartID)
	if err != nil {
		return err
	}
	return nil
}
