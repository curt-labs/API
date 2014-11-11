package video

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
	"time"
)

//Pulls from "video" table, as opposed to "video_new"
type Video_Old struct {
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

func UniqueVideos() (videos []Video_Old, err error) {

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(uniqueVideoStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	for rows.Next() {
		var v Video_Old
		var watch, screen *string
		err = rows.Scan(
			&v.YouTubeId,
			&v.DateAdded,
			&v.Sort,
			&v.Title,
			&v.Description,
			&watch,
			&screen,
		)
		if err != nil {
			return
		}
		v.Watchpage, err = url.Parse(*watch)
		v.Screenshot, err = url.Parse(*screen)
		videos = append(videos, v)
	}
	defer rows.Close()

	return
}
