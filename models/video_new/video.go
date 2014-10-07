package video

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	// "github.com/curt-labs/GoAPI/models/site"
	"github.com/curt-labs/GoAPI/models/products"
	_ "github.com/go-sql-driver/mysql"
	// "log"
	"time"
)

type Video struct {
	ID           int
	Title        string
	VideoType    VideoType
	Description  string
	DateAdded    time.Time
	DateModified time.Time
	IsPrimary    bool
	Thumbnail    string
	Channels     Channels
	Files        CdnFiles
	Categories   []CatAssociation
	Parts        []PartAssociation
	WebsiteId    int //TODO
}
type Videos []Video

type Channel struct {
	ID           int
	Type         ChannelType
	Link         string
	EmbedCode    string
	ForiegnID    string
	DateAdded    time.Time
	DateModified time.Time
	Title        string
	Description  string
}

type Channels []Channel

type ChannelType struct {
	ID          int
	Name        string
	Description string
}

type CdnFile struct {
	ID           int
	Type         CdnFileType
	Path         string
	Bucket       string
	ObjectName   string
	FileSize     string
	DateAdded    time.Time
	DateModified time.Time
	LastUploaded string
}

type CdnFiles []CdnFile

type CdnFileType struct {
	ID          int
	MimeType    string
	Title       string
	Description string
}

type VideoType struct {
	ID   int
	Name string
	Icon string
}

type CatAssociation struct {
	CatID    int
	CatTitle string
	ParentID int
}

type PartAssociation struct {
	PartID    int
	ShortDesc string
	IsPrimary bool
}

const (
	videoFields       = ` v.ID, v.subjectTypeID, v.title, v.description, v.dateAdded, v.dateModified, v.isPrimary, v.thumbnail `
	videoTypeFields   = `  vt.name, vt.icon `
	channelFields     = ` c.ID, c.typeID, c.link, c.embedCode, c.foriegnID, c.dateAdded, c.dateModified, c.title, c.desc `
	channelTypeFields = ` ct.name, ct.description `
	cdnFileFields     = `cf.ID, cf.typeID, cf.path, cf.dateAdded, cf.lastUploaded, cf.bucket, cf.objectName, cf.fileSize `
	cdnFileTypeFields = ` cft.mimeType, cft.title, cft.description `
)

var (
	getVideo      = `SELECT ` + videoFields + `, ` + videoTypeFields + ` FROM VideoNew AS v LEFT JOIN videoType AS vt ON vt.vTypeID = v.subjectTypeID WHERE v.ID = ?`
	getAllVideos  = `SELECT ` + videoFields + `, ` + videoTypeFields + ` FROM VideoNew AS v LEFT JOIN videoType AS vt ON vt.vTypeID = v.subjectTypeID`
	getPartVideos = `SELECT ` + videoFields + `, ` + videoTypeFields + ` FROM VideoNew AS v LEFT JOIN videoType AS vt ON vt.vTypeID = v.subjectTypeID 
						LEFT JOIN VideoJoin AS vj on vj.videoID = v.ID WHERE vj.partID = ? ORDER BY v.title `
	getVideoChannels = `SELECT ` + channelFields + `, ` + channelTypeFields + ` FROM VideoNew AS v 
						LEFT JOIN VideoChannels AS vc ON vc.videoID = v.ID 
						LEFT JOIN Channel AS c ON c.ID = vc.channelID
						LEFT JOIN ChannelType AS ct ON ct.ID = c.typeID 
						WHERE c.ID = ?`
	getVideoCdns = `SELECT ` + cdnFileFields + `, ` + cdnFileTypeFields + ` FROM CdnFile AS cf 
						LEFT JOIN CdnFileType AS cft ON cft.ID = cf.typeID 
						LEFT JOIN VideoCdnFiles AS vcf ON vcf.cdnID = cf.ID 
						WHERE vcf.videoID = ? `
)

func (v *Video) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getVideo)
	if err != nil {
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(v.ID)
	ch := make(chan Video)
	go populateVideo(row, ch)
	*v = <-ch

	return err
}

func GetAllVideos() (vs Videos, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return vs, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllVideos)
	if err != nil {
		return vs, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return vs, err
	}

	ch := make(chan Videos)
	go populateVideos(rows, ch)
	vs = <-ch
	if len(vs) == 0 {
		err = sql.ErrNoRows
	}

	return vs, err
}

func GetPartVideos(p products.Part) (vs Videos, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return vs, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getPartVideos)
	if err != nil {
		return vs, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(p.ID)
	if err != nil {
		return vs, err
	}

	ch := make(chan Videos)
	go populateVideos(rows, ch)
	vs = <-ch
	if len(vs) == 0 {
		err = sql.ErrNoRows
	}

	return vs, err
}

func (v *Video) GetChannels() (chs Channels, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return chs, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getVideoChannels)
	if err != nil {
		return chs, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(v.ID)
	if err != nil {
		return chs, err
	}

	ch := make(chan Channels)
	go populateChannels(rows, ch)
	chs = <-ch
	if len(chs) == 0 {
		err = sql.ErrNoRows
	}

	return chs, err
}

func (v *Video) GetCdnFiles() (cdns CdnFiles, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cdns, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getVideoCdns)
	if err != nil {
		return cdns, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(v.ID)
	if err != nil {
		return cdns, err
	}

	ch := make(chan CdnFiles)
	go populateCdns(rows, ch)
	cdns = <-ch
	if len(cdns) == 0 {
		err = sql.ErrNoRows
	}

	return cdns, err
}

//Populates video and video type fields
func populateVideo(row *sql.Row, ch chan Video) {
	var v Video
	err := row.Scan(
		&v.ID,
		&v.VideoType.ID,
		&v.Title,
		&v.Description,
		&v.DateAdded,
		&v.DateModified,
		&v.IsPrimary,
		&v.Thumbnail,
		&v.VideoType.Name,
		&v.VideoType.Icon,
	)
	if err != nil {
		ch <- Video{}
		return
	}
	ch <- v
	return
}

//Populates video and video type fields
func populateVideos(rows *sql.Rows, ch chan Videos) {
	var v Video
	var vs Videos
	for rows.Next() {
		err := rows.Scan(
			&v.ID,
			&v.VideoType.ID,
			&v.Title,
			&v.Description,
			&v.DateAdded,
			&v.DateModified,
			&v.IsPrimary,
			&v.Thumbnail,
			&v.VideoType.Name,
			&v.VideoType.Icon,
		)
		if err != nil {
			ch <- Videos{}
			return
		}
		vs = append(vs, v)
	}
	ch <- vs
	return
}

//Populates channels and channel type fields
func populateChannels(rows *sql.Rows, ch chan Channels) {
	var c Channel
	var cs Channels
	for rows.Next() {
		err := rows.Scan(
			&c.ID,
			&c.Type.ID,
			&c.Link,
			&c.EmbedCode,
			&c.ForiegnID,
			&c.DateAdded,
			&c.DateModified,
			&c.Title,
			&c.Description,
			&c.Type.Name,
			&c.Type.Description,
		)
		if err != nil {
			ch <- Channels{}
			return
		}
		cs = append(cs, c)
	}
	ch <- cs
	return
}

//Populates channels and channel type fields
func populateCdns(rows *sql.Rows, ch chan CdnFiles) (err error) {
	var c CdnFile
	var cs CdnFiles
	for rows.Next() {
		err := rows.Scan(
			&c.ID,
			&c.Type.ID,
			&c.Path,
			&c.DateAdded,
			&c.LastUploaded,
			&c.Bucket,
			&c.ObjectName,
			&c.FileSize,
			&c.Type.MimeType,
			&c.Type.Title,
			&c.Type.Description,
		)
		if err != nil {
			ch <- CdnFiles{}
			return err
		}
		cs = append(cs, c)
	}
	ch <- cs
	return nil
}
