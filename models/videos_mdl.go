package models

import (
	"../helpers/database"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type PartVideo struct {
	YouTubeVideoId, Type string
	IsPrimary            bool
	TypeIcon             *url.URL
}

type Video struct {
	YouTubeId   string
	DateAdded   time.Time
	Sort        int
	Title       string
	Description string
	Watchpage   *url.URL
	Screenshot  *url.URL
}

var (
	partVideoStmt = `select pv.video,vt.name,pv.isPrimary, vt.icon from PartVideo as pv
				join videoType vt on pv.vTypeID = vt.vTypeID
				where pv.partID = ?`

	partVideoStmt_Grouped = `select pv.partID, pv.video,vt.name,pv.isPrimary, vt.icon from PartVideo as pv
				join videoType vt on pv.vTypeID = vt.vTypeID
				where pv.partID IN (%s)`

	uniqueVideoStmt = `select distinct embed_link, dateAdded, sort, title, description, watchpage, screenshot 
				from Video
				order by sort`
)

func (p *Part) GetVideos() error {
	if p == nil {
		return errors.New("Part is nil")
	}

	qry, err := database.Db.Prepare(partVideoStmt)

	if p == nil {
		return errors.New("Part is nil")
	}
	rows, res, err := qry.Exec(p.PartId)
	if database.MysqlError(err) {
		return err
	}
	video := res.Map("video")
	typ := res.Map("name")
	prime := res.Map("isPrimary")
	icon := res.Map("icon")

	var videos []PartVideo
	for _, row := range rows {
		isPrime := false
		if row.Int(prime) == 1 {
			isPrime = true
		}
		iconPath, urlErr := url.Parse(row.Str(icon))
		if urlErr != nil {
			iconPath = nil
		}
		v := PartVideo{
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
	if len(ids) == 0 {
		return nil
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

	videos := make(map[int][]PartVideo, len(lookup.Parts))

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
		v := PartVideo{
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

func UniqueVideos() (videos []Video, err error) {
	rows, res, err := database.Db.Query(uniqueVideoStmt)
	if database.MysqlError(err) || len(rows) == 0 {
		return
	}

	youTubeId := res.Map("embed_link")
	dateAdded := res.Map("dateAdded")
	sort := res.Map("sort")
	title := res.Map("title")
	description := res.Map("description")
	watchpage := res.Map("watchpage")
	screenshot := res.Map("screenshot")

	for _, row := range rows {

		date_add, _ := time.Parse("2006-01-02 15:04:15", row.Str(dateAdded))

		page, _ := url.Parse(row.Str(watchpage))
		shot, _ := url.Parse(row.Str(screenshot))

		video := Video{
			YouTubeId:   row.Str(youTubeId),
			DateAdded:   date_add,
			Sort:        row.Int(sort),
			Title:       row.Str(title),
			Description: row.Str(description),
			Watchpage:   page,
			Screenshot:  shot,
		}
		videos = append(videos, video)
	}

	return
}
