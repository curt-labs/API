package video

import (
	"database/sql"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
	"github.com/curt-labs/GoAPI/models/products"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

type Video struct {
	ID           int                 `json:"id,omitempty" xml:"id,omitempty"`
	Title        string              `json:"title, omitempty" xml:"title,omitempty"`
	VideoType    VideoType           `json:"videoType,omitempty" xml:"v,omitempty"`
	Description  string              `json:"description,omitempty" xml:"description,omitempty"`
	DateAdded    time.Time           `json:"dateAdded,omitempty" xml:"dateAdded,omitempty"`
	DateModified time.Time           `json:"dateModified,omitempty" xml:"dateModified,omitempty"`
	IsPrimary    bool                `json:"isPrimary,omitempty" xml:"v,omitempty"`
	Thumbnail    string              `json:"thumbnail,omitempty" xml:"thumbnail,omitempty"`
	Channels     Channels            `json:"channels,omitempty" xml:"channels,omitempty"`
	Files        CdnFiles            `json:"files,omitempty" xml:"files,omitempty"`
	Categories   []products.Category `json:"categories,omitempty" xml:"categories,omitempty"`
	Parts        []products.Part     `json:"parts,omitempty" xml:"parts,omitempty"`
	WebsiteId    int                 `json:"websiteId,omitempty" xml:"websiteId,omitempty"` //TODO
}
type Videos []Video

type Channel struct {
	ID           int         `json:"id,omitempty" xml:"id,omitempty"`
	Type         ChannelType `json:"type,omitempty" xml:"type,omitempty"`
	Link         string      `json:"link,omitempty" xml:"link,omitempty"`
	EmbedCode    string      `json:"embedCode,omitempty" xml:"embedCode,omitempty"`
	ForiegnID    string      `json:"foreignId,omitempty" xml:"foreignId,omitempty"`
	DateAdded    time.Time   `json:"dateAdded,omitempty" xml:"dateAdded,omitempty"`
	DateModified time.Time   `json:"dateModified,omitempty" xml:"dateModified,omitempty"`
	Title        string      `json:"title,omitempty" xml:"title,omitempty"`
	Description  string      `json:"description,omitempty" xml:"description,omitempty"`
}

type Channels []Channel

type ChannelType struct {
	ID          int    `json:"id,omitempty" xml:"id,omitempty"`
	Name        string `json:"name,omitempty" xml:"name,omitempty"`
	Description string `json:"description,omitempty" xml:"description,omitempty"`
}

type CdnFile struct {
	ID           int         `json:"id,omitempty" xml:"id,omitempty"`
	Type         CdnFileType `json:"type,omitempty" xml:"type,omitempty"`
	Path         string      `json:"path,omitempty" xml:"path,omitempty"`
	Bucket       string      `json:"bucket,omitempty" xml:"bucket,omitempty"`
	ObjectName   string      `json:"objectName,omitempty" xml:"objectName,omitempty"`
	FileSize     string      `json:"fileSize,omitempty" xml:"fileSize,omitempty"`
	DateAdded    time.Time   `json:"dateAdded,omitempty" xml:"dateAdded,omitempty"`
	DateModified time.Time   `json:"dateModified,omitempty" xml:"dateModified,omitempty"`
	LastUploaded string      `json:"lastUploaded,omitempty" xml:"lastUploaded,omitempty"`
}

type CdnFiles []CdnFile

type CdnFileType struct {
	ID          int    `json:"id,omitempty" xml:"id,omitempty"`
	MimeType    string `json:"mimeType,omitempty" xml:"mimeType,omitempty"`
	Title       string `json:"title,omitempty" xml:"title,omitempty"`
	Description string `json:"description,omitempty" xml:"description,omitempty"`
}

type VideoType struct {
	ID   int    `json:"id,omitempty" xml:"id,omitempty"`
	Name string `json:"name,omitempty" xml:"name,omitempty"`
	Icon string `json:"icon,omitempty" xml:"icon,omitempty"`
}

//TODO categories should be their own entity
type Category struct {
	ID    int    `json:"id,omitempty" xml:"id,omitempty"`
	Title string `json:"title,omitempty" xml:"title,omitempty"`
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
	getVideo         = `SELECT ` + videoFields + `, ` + videoTypeFields + ` FROM VideoNew AS v LEFT JOIN videoType AS vt ON vt.vTypeID = v.subjectTypeID WHERE v.ID = ?`
	getAllVideos     = `SELECT ` + videoFields + `, ` + videoTypeFields + ` FROM VideoNew AS v LEFT JOIN videoType AS vt ON vt.vTypeID = v.subjectTypeID`
	getAllCdnFiles   = `SELECT ` + cdnFileFields + `,` + cdnFileTypeFields + ` FROM CdnFile AS cf LEFT JOIN CdnFileType AS cft ON cft.ID = cf.typeID `
	getAllChannels   = `SELECT ` + channelFields + `, ` + channelTypeFields + ` FROM Channel AS c LEFT JOIN ChannelType AS ct ON ct.ID = c.typeID `
	getAllVideoTypes = `SELECT vt.vTypeID, ` + videoTypeFields + ` FROM videoType AS vt`
	getPartVideos    = `SELECT ` + videoFields + `, ` + videoTypeFields + ` FROM VideoNew AS v LEFT JOIN videoType AS vt ON vt.vTypeID = v.subjectTypeID 
						LEFT JOIN VideoJoin AS vj on vj.videoID = v.ID WHERE vj.partID = ? ORDER BY v.title `
	getVideoChannels = `SELECT ` + channelFields + `, ` + channelTypeFields + ` FROM VideoChannels AS vc 
					  JOIN Channel AS c on c.ID = vc.channelID
					 LEFT JOIN ChannelType AS ct ON ct.ID = c.typeID
					WHERE vc.videoID = ?`
	getVideoCdns = `SELECT ` + cdnFileFields + `, ` + cdnFileTypeFields + ` FROM CdnFile AS cf 
						LEFT JOIN CdnFileType AS cft ON cft.ID = cf.typeID 
						LEFT JOIN VideoCdnFiles AS vcf ON vcf.cdnID = cf.ID 
						WHERE vcf.videoID = ? `

	createVideo             = `INSERT INTO VideoNew (subjectTypeID, title, description, dateAdded, dateModified, isPrimary, thumbnail) VALUES(?, ?, ?, ?, ?, ?, ?)`
	updateVideo             = `UPDATE videoNew SET subjectTypeID = ?, title = ?, description = ?, isPrimary = ?, thumbnail = ? WHERE ID = ?`
	deleteVideo             = `DELETE FROM videoNew WHERE ID = ?`
	joinVideoCdn            = `INSERT INTO VideoCdnFiles(cdnID, videoID) VALUES(?,?)`
	joinVideoChannel        = `INSERT INTO VideoChannels( channelID, videoID) VALUES(?,?)`
	joinVideoPart           = `INSERT INTO VideoJoin(videoID, partID, catID, websiteID, isPrimary) VALUES(?,?,0,?,?)`
	joinVideoCategory       = `INSERT INTO VideoJoin(videoID, partID, catID, websiteID, isPrimary) VALUES(?,0,?,?,?)`
	deleteVideoCdnJoin      = `DELETE FROM VideoCdnFiles WHERE videoID = ?`
	deleteVideoChannelJoin  = `DELETE FROM VideoChannels WHERE videoID = ?`
	deleteVideoPartJoin     = `DELETE FROM VideoJoin WHERE videoID = ? AND partID = ?`
	deleteVideoCategoryJoin = `DELETE FROM VideoJoin WHERE videoID = ? AND catID = ?`

	//crud
	getChannel         = "SELECT ID, typeID, link, embedCode, foriegnID, dateAdded, dateModified, title, `desc` FROM Channel WHERE ID = ?"
	createChannel      = "INSERT INTO Channel (typeID, link, embedCode, foriegnID, dateAdded, title, `desc`) VALUES (?,?,?,?,?,?,?)"
	updateChannel      = "UPDATE Channel SET typeID = ?, link = ?, embedCode = ?, foriegnID = ?, title = ?, `desc` = ? WHERE ID = ?"
	deleteChannel      = "DELETE FROM Channel WHERE ID = ?"
	getCdn             = `SELECT ` + cdnFileFields + `, ` + cdnFileTypeFields + `, cf.dateModified FROM CdnFile AS cf LEFT JOIN CdnFileType AS cft ON cft.ID = cf.typeID WHERE cf.ID = ?`
	createCdn          = `INSERT INTO CdnFile (typeID, path, dateAdded, lastUploaded, bucket, objectName, fileSize) VALUES (?,?,?,?,?,?,?)`
	updateCdn          = `UPDATE CdnFile SET typeID = ?, path = ?, lastUploaded = ?, bucket = ?, objectName = ?, fileSize = ? WHERE ID = ?`
	deleteCdn          = `DELETE FROM CdnFile WHERE ID = ?`
	getCdnType         = `SELECT ID, mimeType, title, description FROM CdnFileType WHERE ID = ?`
	getAllCdnTypes     = `SELECT ID, mimeType, title, description FROM CdnFileType`
	createCdnType      = `INSERT INTO CdnFileType (mimeType, title, description) VALUES(?,?,?)`
	updateCdnType      = `UPDATE CdnFileType SET mimeType = ?, title = ?, description = ? WHERE ID = ?`
	deleteCdnType      = `DELETE FROM CdnFileType WHERE ID = ?`
	getVideoType       = `SELECT vTypeID, name, icon FROM videoType WHERE vTypeID = ?`
	createVideoType    = `INSERT INTO videoType (name, icon) VALUES (?,?)`
	updateVideoType    = `UPDATE videoType SET name = ?, icon = ? WHERE vTypeID = ?`
	deleteVideoType    = `DELETE FROM videoType WHERE vTypeID = ?`
	getChannelType     = `SELECT ID, name, description FROM ChannelType WHERE ID = ?`
	getAllChannelTypes = `SELECT ID, name, description FROM ChannelType `
	createChannelType  = `INSERT INTO ChannelType (name, description) VALUES (?,?)`
	updateChannelType  = `UPDATE ChannelType SET name = ?, description = ? WHERE ID = ?`
	deleteChannelType  = `DELETE FROM ChannelType WHERE ID = ?`
)

//Base Video
func (v *Video) Get() (err error) {
	redis_key := "video:" + strconv.Itoa(v.ID)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &v)
		return err
	}
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
	go redis.Setex(redis_key, v, 86400)
	return err
}

func (v *Video) GetVideoDetails() (err error) {
	redis_key := "video:details:" + strconv.Itoa(v.ID)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &v)
		return err
	}
	baseChan := make(chan bool)
	chanChan := make(chan bool)
	cdnChan := make(chan bool)
	go func() (err error) {
		err = v.Get()
		if err != nil {
			return err
		}
		baseChan <- true
		return err
	}()
	go func() (err error) {
		chs, err := v.GetChannels()
		if err != nil {
			return err
		}
		v.Channels = chs
		chanChan <- true
		return err
	}()
	go func() (err error) {
		cdns, err := v.GetCdnFiles()
		if err != nil {
			return err
		}
		v.Files = cdns
		cdnChan <- true
		return err
	}()

	<-baseChan
	<-chanChan
	<-cdnChan
	go redis.Setex(redis_key, v, 86400)
	return nil
}

func GetAllVideos() (vs Videos, err error) {
	redis_key := "video"
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &vs)
		return vs, err
	}
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
	go redis.Setex(redis_key, vs, 86400)
	return vs, err
}

func GetPartVideos(p products.Part) (vs Videos, err error) {
	redis_key := "video:part:" + strconv.Itoa(p.ID)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &vs)
		return vs, err
	}
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
	go redis.Setex(redis_key, vs, 86400)
	return vs, err
}

func (v *Video) GetChannels() (chs Channels, err error) {
	redis_key := "video:channels:" + strconv.Itoa(v.ID)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &chs)
		return chs, err
	}
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

	if err == nil {
		ch := make(chan Channels)
		go populateChannels(rows, ch)
		chs = <-ch
	}
	go redis.Setex(redis_key, chs, 86400)
	return chs, err
}

func (v *Video) GetCdnFiles() (cdns CdnFiles, err error) {
	redis_key := "video:cdnFiles:" + strconv.Itoa(v.ID)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &cdns)
		return cdns, err
	}
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
	if err == nil {
		ch := make(chan CdnFiles)
		go populateCdns(rows, ch)
		cdns = <-ch
	}
	go redis.Setex(redis_key, cdns, 86400)
	return cdns, err
}

func (v *Video) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(createVideo)
	if err != nil {
		return err
	}
	defer stmt.Close()
	v.DateAdded = time.Now()
	res, err := stmt.Exec(v.VideoType.ID, v.Title, v.Description, v.DateAdded, v.DateModified, v.IsPrimary, v.Thumbnail)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	id, err := res.LastInsertId()
	v.ID = int(id)

	// create joins
	fChan := make(chan int)
	chChan := make(chan int)
	catChan := make(chan int)
	pChan := make(chan int)

	go func() (err error) {
		if len(v.Files) > 0 {
			for _, file := range v.Files {
				err = v.CreateJoinFile(file)
			}
		}
		fChan <- 1
		return err
	}()
	go func() (err error) {
		if len(v.Channels) > 0 {
			for _, channel := range v.Channels {
				err = v.CreateJoinChannel(channel)
			}
		}
		chChan <- 1
		return err
	}()
	go func() (err error) {
		if len(v.Categories) > 0 {
			for _, cat := range v.Categories {
				err = v.CreateJoinCategory(cat)
			}
		}
		catChan <- 1
		return err
	}()
	go func() (err error) {
		if len(v.Parts) > 0 {
			for _, part := range v.Parts {
				err = v.CreateJoinPart(part)
			}
		}
		pChan <- 1
		return err
	}()

	<-fChan
	<-chChan
	<-catChan
	<-pChan

	return err
}

func (v *Video) Update() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(updateVideo)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(v.VideoType.ID, v.Title, v.Description, v.IsPrimary, v.Thumbnail, v.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	// delete and create joins
	fChan := make(chan int)
	chChan := make(chan int)
	catChan := make(chan int)
	pChan := make(chan int)

	go func() (err error) {
		err = v.DeleteJoinFiles()
		if err != nil {
			return err
		}
		if len(v.Files) > 0 {
			for _, file := range v.Files {
				err = v.CreateJoinFile(file)
				if err != nil {
					return err
				}
			}
		}
		fChan <- 1
		return nil
	}()
	go func() (err error) {
		err = v.DeleteJoinChannels()
		if err != nil {
			return err
		}
		if len(v.Channels) > 0 {
			for _, channel := range v.Channels {
				err = v.CreateJoinChannel(channel)
				if err != nil {
					return err
				}
			}

		}
		chChan <- 1
		return nil
	}()
	go func() (err error) {
		if len(v.Categories) > 0 {
			for _, cat := range v.Categories {
				err = v.DeleteJoinCategory(cat)
				if err != nil {
					return err
				}
				err = v.CreateJoinCategory(cat)
				if err != nil {
					return err
				}
			}
		}
		catChan <- 1
		return nil
	}()
	go func() (err error) {
		if len(v.Parts) > 0 {
			for _, part := range v.Parts {
				err = v.DeleteJoinPart(part)
				if err != nil {
					return err
				}
				err = v.CreateJoinPart(part)
				if err != nil {
					return err
				}
			}
		}
		pChan <- 1
		return nil
	}()
	<-fChan
	<-chChan
	<-catChan
	<-pChan

	return err
}

func (v *Video) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteVideo)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(v.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	//delete and create joins
	fChan := make(chan int)
	chChan := make(chan int)
	catChan := make(chan int)
	pChan := make(chan int)

	go func() (err error) {
		if len(v.Files) > 0 {
			err = v.DeleteJoinFiles()
			if err != nil {
				return err
			}
		}
		fChan <- 1
		return nil
	}()
	go func() (err error) {
		if len(v.Channels) > 0 {
			err = v.DeleteJoinChannels()
			if err != nil {
				return err
			}
		}
		chChan <- 1
		return nil
	}()
	go func() (err error) {
		if len(v.Categories) > 0 {
			for _, cat := range v.Categories {
				err = v.DeleteJoinCategory(cat)
				if err != nil {
					return err
				}
			}
		}
		catChan <- 1
		return nil
	}()
	go func() (err error) {
		if len(v.Parts) > 0 {
			for _, part := range v.Parts {
				err = v.DeleteJoinPart(part)
				if err != nil {
					return err
				}
			}
		}
		pChan <- 1
		return nil
	}()
	<-fChan
	<-chChan
	<-catChan
	<-pChan

	return err
}

func (v *Video) CreateJoinFile(f CdnFile) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(joinVideoCdn)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(f.ID, v.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

//Youtube
func (v *Video) CreateJoinChannel(channel Channel) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(joinVideoChannel)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(channel.ID, v.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

func (v *Video) CreateJoinPart(p products.Part) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(joinVideoPart)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(v.ID, p.ID, v.WebsiteId, v.IsPrimary)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

//HTML5
func (v *Video) CreateJoinCategory(c products.Category) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(joinVideoCategory)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(v.ID, c.ID, v.WebsiteId, v.IsPrimary)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

func (v *Video) DeleteJoinFiles() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteVideoCdnJoin)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(v.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

//Youtube
func (v *Video) DeleteJoinChannels() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteVideoChannelJoin)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(v.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

func (v *Video) DeleteJoinPart(p products.Part) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteVideoPartJoin)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(v.ID, p.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

//HTML5
func (v *Video) DeleteJoinCategory(c products.Category) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteVideoCategoryJoin)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(v.ID, c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

func (c *Channel) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getChannel)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var link, foreign, title *string
	err = stmt.QueryRow(c.ID).Scan(
		&c.ID,
		&c.Type.ID,
		&link,
		&c.EmbedCode,
		&foreign,
		&c.DateAdded,
		&c.DateModified,
		&title,
		&c.Description,
	)
	if err != nil {
		return err
	}
	if link != nil {
		c.Link = *link
	}
	if foreign != nil {
		c.ForiegnID = *foreign
	}
	if title != nil {
		c.Title = *title
	}
	return err
}

func GetAllChannels() (cs Channels, err error) {
	redis_key := "video:channels"
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &cs)
		return cs, err
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cs, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllChannels)
	if err != nil {
		return cs, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return cs, err
	}

	ch := make(chan Channels)
	go populateChannels(rows, ch)
	cs = <-ch
	if len(cs) == 0 {
		err = sql.ErrNoRows
	}
	go redis.Setex(redis_key, cs, 86400)
	return cs, err
}

func (c *Channel) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(createChannel)
	if err != nil {
		return err
	}
	defer stmt.Close()
	c.DateAdded = time.Now()
	res, err := stmt.Exec(c.Type.ID, c.Link, c.EmbedCode, c.ForiegnID, c.DateAdded, c.Title, c.Description)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	id, err := res.LastInsertId()
	c.ID = int(id)
	return err
}

func (c *Channel) Update() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(updateChannel)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.Type.ID, c.Link, c.EmbedCode, c.ForiegnID, c.Title, c.Description, c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}
func (c *Channel) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteChannel)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (c *CdnFile) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCdn)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var last, bucket, object, size, desc, tMime, tTitle *string
	err = stmt.QueryRow(c.ID).Scan(
		&c.ID,
		&c.Type.ID,
		&c.Path,
		&c.DateAdded,
		&last,
		&bucket,
		&object,
		&size,
		&tMime,
		&tTitle,
		&desc,
		&c.DateModified,
	)
	if err != nil {
		return err
	}
	if last != nil {
		c.LastUploaded = *last
	}
	if bucket != nil {
		c.Bucket = *bucket
	}
	if object != nil {
		c.ObjectName = *object
	}
	if size != nil {
		c.FileSize = *size
	}
	if desc != nil {
		c.Type.Description = *desc
	}
	if tMime != nil {
		c.Type.MimeType = *tMime
	}
	if tTitle != nil {
		c.Type.Title = *tTitle
	}
	return err
}

func GetAllCdnFiles() (cs CdnFiles, err error) {
	redis_key := "video:cdnFiles"
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &cs)
		return cs, err
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cs, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllCdnFiles)
	if err != nil {
		return cs, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return cs, err
	}

	ch := make(chan CdnFiles)
	go populateCdns(rows, ch)
	cs = <-ch
	if len(cs) == 0 {
		err = sql.ErrNoRows
	}
	go redis.Setex(redis_key, cs, 86400)
	return cs, err
}

func (c *CdnFile) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(createCdn)
	if err != nil {
		return err
	}
	defer stmt.Close()
	c.DateAdded = time.Now()
	res, err := stmt.Exec(c.Type.ID, c.Path, c.DateAdded, c.LastUploaded, c.Bucket, c.ObjectName, c.FileSize)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	id, err := res.LastInsertId()
	c.ID = int(id)
	return err
}

func (c *CdnFile) Update() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(updateCdn)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.Type.ID, c.Path, c.LastUploaded, c.Bucket, c.ObjectName, c.FileSize, c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}
func (c *CdnFile) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteCdn)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (c *CdnFileType) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCdnType)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var desc *string
	err = stmt.QueryRow(c.ID).Scan(
		&c.ID,
		&c.MimeType,
		&c.Title,
		&desc,
	)
	if err != nil {
		return err
	}
	if desc != nil {
		c.Description = *desc
	}
	return err
}

func GetAllCdnFileTypes() (cts []CdnFileType, err error) {
	redis_key := "video:cdnFileTypes"
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &cts)
		return cts, err
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cts, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllCdnTypes)
	if err != nil {
		return cts, err
	}
	defer stmt.Close()

	res, err := stmt.Query()
	if err != nil {
		return cts, err
	}
	var c CdnFileType
	for res.Next() {
		err = res.Scan(
			&c.ID,
			&c.MimeType,
			&c.Title,
			&c.Description,
		)
		if err != nil {
			return cts, err
		}
		cts = append(cts, c)
	}
	defer res.Close()
	go redis.Setex(redis_key, cts, 86400)
	return cts, err
}

func (c *CdnFileType) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(createCdnType)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(c.MimeType, c.Title, c.Description)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	id, err := res.LastInsertId()
	c.ID = int(id)
	return err
}

func (c *CdnFileType) Update() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(updateCdnType)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.MimeType, c.Title, c.Description, c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}
func (c *CdnFileType) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteCdnType)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (c *VideoType) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getVideoType)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(c.ID).Scan(
		&c.ID,
		&c.Name,
		&c.Icon,
	)
	if err != nil {
		return err
	}
	return err
}

func GetAllVideoTypes() (vts []VideoType, err error) {
	redis_key := "video:videoTypes"
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &vts)
		return vts, err
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return vts, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllVideoTypes)
	if err != nil {
		return vts, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return vts, err
	}
	var vt VideoType
	var vName, vIcon *string
	for rows.Next() {
		err = rows.Scan(&vt.ID, &vName, &vIcon)
		if err != nil {
			return vts, err
		}
		if vName != nil {
			vt.Name = *vName
		}
		if vIcon != nil {
			vt.Icon = *vIcon
		}
		vts = append(vts, vt)
	}
	defer rows.Close()
	go redis.Setex(redis_key, vts, 86400)
	return vts, err
}

func (c *VideoType) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(createVideoType)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(c.Name, c.Icon)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	id, err := res.LastInsertId()
	c.ID = int(id)
	return err
}

func (c *VideoType) Update() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(updateVideoType)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.Name, c.Icon, c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}
func (c *VideoType) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteVideoType)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (c *ChannelType) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getChannelType)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(c.ID).Scan(
		&c.ID,
		&c.Name,
		&c.Description,
	)
	if err != nil {
		return err
	}
	return err
}

func GetAllChannelTypes() (cts []ChannelType, err error) {
	redis_key := "video:channelTypes"
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &cts)
		return cts, err
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cts, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllChannelTypes)
	if err != nil {
		return cts, err
	}
	defer stmt.Close()
	var c ChannelType
	res, err := stmt.Query()
	if err != nil {
		return cts, err
	}
	for res.Next() {
		err = res.Scan(
			&c.ID,
			&c.Name,
			&c.Description,
		)
		if err != nil {
			return cts, err
		}
		cts = append(cts, c)
	}
	defer res.Close()
	go redis.Setex(redis_key, cts, 86400)
	return cts, err
}

func (c *ChannelType) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(createChannelType)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(c.Name, c.Description)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	id, err := res.LastInsertId()
	c.ID = int(id)
	return err
}

func (c *ChannelType) Update() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(updateChannelType)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.Name, c.Description, c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}
func (c *ChannelType) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteChannelType)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

//Populates a video + type
func populateVideo(row *sql.Row, ch chan Video) {
	var v Video
	var tName, tIcon *string
	err := row.Scan(
		&v.ID,
		&v.VideoType.ID,
		&v.Title,
		&v.Description,
		&v.DateAdded,
		&v.DateModified,
		&v.IsPrimary,
		&v.Thumbnail,
		&tName,
		&tIcon,
	)
	if err != nil {
		ch <- v
		return
	}
	if tName != nil {
		v.VideoType.Name = *tName
	}
	if tIcon != nil {
		v.VideoType.Icon = *tIcon
	}

	ch <- v
	return
}

//Populates video and video type fields
func populateVideos(rows *sql.Rows, ch chan Videos) {
	var v Video
	var vs Videos
	var tName, tIcon *string
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
			&tName,
			&tIcon,
		)
		if err != nil {
			ch <- Videos{}
			return
		}
		if tName != nil {
			v.VideoType.Name = *tName
		}
		if tIcon != nil {
			v.VideoType.Icon = *tIcon
		}
		vs = append(vs, v)
	}
	defer rows.Close()
	ch <- vs
	return
}

//Populates channels and channel type fields
func populateCdns(rows *sql.Rows, ch chan CdnFiles) {
	var c CdnFile
	var cs CdnFiles
	var last, object, bucket, size, tMime, tTitle, tDesc *string
	for rows.Next() {
		err := rows.Scan(
			&c.ID,
			&c.Type.ID,
			&c.Path,
			&c.DateAdded,
			&last,
			&bucket,
			&object,
			&size,
			&tMime,
			&tTitle,
			&tDesc,
		)
		if err != nil {
			ch <- CdnFiles{}
			return
		}
		if last != nil {
			c.LastUploaded = *last
		}
		if bucket != nil {
			c.Bucket = *bucket
		}
		if object != nil {
			c.ObjectName = *object
		}
		if size != nil {
			c.FileSize = *size
		}
		if tDesc != nil {
			c.Type.Description = *tDesc
		}
		if tMime != nil {
			c.Type.MimeType = *tMime
		}
		if tTitle != nil {
			c.Type.Title = *tTitle
		}

		cs = append(cs, c)
	}
	defer rows.Close()
	ch <- cs
	return
}

//populate channels
func populateChannels(rows *sql.Rows, ch chan Channels) {
	var chs Channels
	var c Channel
	var link, foreign, title, tName, tDesc *string
	var ctypeId *int
	for rows.Next() {
		err := rows.Scan(
			&c.ID,
			&ctypeId,
			&link,
			&c.EmbedCode,
			&foreign,
			&c.DateAdded,
			&c.DateModified,
			&title,
			&c.Description,
			&tName,
			&tDesc,
		)
		if err != nil {
			ch <- Channels{}
			return
		}
		if link != nil {
			c.Link = *link
		}
		if ctypeId != nil {
			c.Type.ID = *ctypeId
		}
		if foreign != nil {
			c.ForiegnID = *foreign
		}
		if title != nil {
			c.Title = *title
		}
		if tName != nil {
			c.Type.Name = *tName
		}
		if tDesc != nil {
			c.Type.Description = *tDesc
		}
		chs = append(chs, c)
	}
	defer rows.Close()
	ch <- chs
	return
}
