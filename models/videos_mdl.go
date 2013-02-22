package models

import (
	"../helpers/database"
	"net/url"
)

type Video struct {
	YouTubeVideoId, Type string
	IsPrimary            bool
	TypeIcon             *url.URL
}

var (
	partVideoStmt = `select pv.video,vt.name,pv.isPrimary, vt.icon from PartVideo as pv
				join videoType vt on pv.vTypeID = vt.vTypeID
				where pv.partID = %d`
)

func (p *Part) GetVideos() error {
	db := database.Db

	rows, res, err := db.Query(partVideoStmt, p.PartId)
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
