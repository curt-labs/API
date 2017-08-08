package video

import (
	"net/url"
	"time"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

//Pulls from "Video" table, as opposed to "video_new", used for legacy support.
type Video_Old struct {
	YouTubeId   string
	DateAdded   time.Time
	Sort        int
	Title       string
	Description string
	Watchpage   *url.URL
	Screenshot  *url.URL
	BrandID     int
}

var (
	uniqueVideoStmt = `select distinct v.embed_link, v.dateAdded, v.sort, v.title, v.description, v.watchpage, v.screenshot, v.brandID
		from Video as v
		join ApiKeyToBrand as akb on akb.brandID = v.brandID
		join ApiKey as ak on ak.id = akb.keyID
        && ak.api_key = ? && (v.brandID = ? or 0 = ?)
        order by sort`
)

// Gets a list of all of the old videos - used for legacy support.
func UniqueVideos(dtx *apicontext.DataContext) (videos []Video_Old, err error) {

	err = database.Init()
	if err != nil {
		return
	}

	stmt, err := database.DB.Prepare(uniqueVideoStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
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
			&v.BrandID,
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
