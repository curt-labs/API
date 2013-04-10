package models

import (
	"../helpers/database"
	"net/url"
	"strconv"
	"strings"
)

type Video struct {
	YouTubeVideoId, Type string
	IsPrimary            bool
	TypeIcon             *url.URL
}

var (
	partVideoStmt = `select pv.video,vt.name,pv.isPrimary, vt.icon from PartVideo as pv
				join videoType vt on pv.vTypeID = vt.vTypeID
				where pv.partID = ?`

	partVideoStmt_Grouped = `select pv.partID, pv.video,vt.name,pv.isPrimary, vt.icon from PartVideo as pv
				join videoType vt on pv.vTypeID = vt.vTypeID
				where pv.partID IN (%s)`
)

func (p *Part) GetVideos() error {
	qry, err := database.Db.Prepare(partVideoStmt)

	rows, res, err := qry.Exec(p.PartId)
	if database.MysqlError(err) {
		return err
	}
	video := res.Map("video")
	typ := res.Map("name")
	prime := res.Map("isPrimary")
	icon := res.Map("icon")

	var videos []Video
	for _, row := range rows {
		isPrime := false
		if row.Int(prime) == 1 {
			isPrime = true
		}
		iconPath, urlErr := url.Parse(row.Str(icon))
		if urlErr != nil {
			iconPath = nil
		}
		v := Video{
			YouTubeVideoId: row.Str(video),
			Type:           row.Str(typ),
			IsPrimary:      isPrime,
			TypeIcon:       iconPath,
		}
		videos = append(videos, v)
	}

	p.Videos = videos

	return nil
}

func (lookup *Lookup) GetVideos() error {

	var ids []string
	for _, p := range lookup.Parts {
		ids = append(ids, strconv.Itoa(p.PartId))
	}

	rows, res, err := database.Db.Query(partVideoStmt_Grouped, strings.Join(ids, ","))
	if database.MysqlError(err) || len(rows) == 0 {
		return err
	}

	partID := res.Map("partID")
	video := res.Map("video")
	typ := res.Map("name")
	prime := res.Map("isPrimary")
	icon := res.Map("icon")

	videos := make(map[int][]Video, len(lookup.Parts))

	for _, row := range rows {
		pId := row.Int(partID)
		isPrime := false
		if row.Int(prime) == 1 {
			isPrime = true
		}
		iconPath, urlErr := url.Parse(row.Str(icon))
		if urlErr != nil {
			iconPath = nil
		}
		v := Video{
			YouTubeVideoId: row.Str(video),
			Type:           row.Str(typ),
			IsPrimary:      isPrime,
			TypeIcon:       iconPath,
		}

		videos[pId] = append(videos[pId], v)
	}

	for _, p := range lookup.Parts {
		p.Videos = videos[p.PartId]
	}

	return nil
}
