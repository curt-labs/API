package video

import (
	"github.com/curt-labs/GoAPI/helpers/database"
	"net/url"
	"time"
)

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
	uniqueVideoStmt = `select distinct embed_link, dateAdded, sort, title, description, watchpage, screenshot
				from Video
				order by sort`
)

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
