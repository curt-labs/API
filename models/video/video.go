package video

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
	"github.com/curt-labs/GoAPI/models/brand"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"encoding/json"
	"strconv"
	"time"
)

type Video struct {
	ID           int          `json:"id,omitempty" xml:"id,omitempty"`
	Title        string       `json:"title, omitempty" xml:"title,omitempty"`
	VideoType    VideoType    `json:"videoType,omitempty" xml:"v,omitempty"`
	Description  string       `json:"description,omitempty" xml:"description,omitempty"`
	DateAdded    time.Time    `json:"dateAdded,omitempty" xml:"dateAdded,omitempty"`
	DateModified time.Time    `json:"dateModified,omitempty" xml:"dateModified,omitempty"`
	IsPrimary    bool         `json:"isPrimary,omitempty" xml:"v,omitempty"`
	Thumbnail    string       `json:"thumbnail,omitempty" xml:"thumbnail,omitempty"`
	Channels     Channels     `json:"channels,omitempty" xml:"channels,omitempty"`
	Files        CdnFiles     `json:"files,omitempty" xml:"files,omitempty"`
	CategoryIds  []int        `json:"categoryIds,omitempty" xml:"categoryIds,omitempty"`
	PartIds      []int        `json:"partIds,omitempty" xml:"partIds,omitempty"`
	WebsiteId    int          `json:"websiteId,omitempty" xml:"websiteId,omitempty"` //TODO
	Brands       brand.Brands `json:"brands,omitempty" xml:"brands,omitempty"`
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
	videoFields             = ` v.ID, v.subjectTypeID, v.title, v.description, v.dateAdded, v.dateModified, v.isPrimary, v.thumbnail `
	videoTypeFields         = `  vt.name, vt.icon `
	channelFields           = ` c.ID, c.typeID, c.link, c.embedCode, c.foriegnID, c.dateAdded, c.dateModified, c.title, c.desc `
	channelTypeFields       = ` ct.name, ct.description `
	cdnFileFields           = `cf.ID, cf.typeID, cf.path, cf.dateAdded, cf.lastUploaded, cf.bucket, cf.objectName, cf.fileSize `
	cdnFileTypeFields       = ` cft.mimeType, cft.title, cft.description `
	AllCdnFileTypeRedisKey  = "video:cdnFileTypes"
	AllVideoTypesRedisKey   = "video:videoTypes"
	AllChannelTypesRedisKey = "video:channelTypes"
	AllCdnFilesRedisKey     = "video:cdnFiles"
	AllChannelsRedisKey     = "video:channels"
	AllVideosRedisKey       = "video"
)

var (
	getVideo     = `SELECT ` + videoFields + `, ` + videoTypeFields + ` FROM VideoNew AS v LEFT JOIN videoType AS vt ON vt.vTypeID = v.subjectTypeID WHERE v.ID = ?`
	getAllVideos = `SELECT ` + videoFields + `, ` + videoTypeFields + ` FROM VideoNew AS v LEFT JOIN videoType AS vt ON vt.vTypeID = v.subjectTypeID
			join VideoNewToBrand as vtb on vtb.videoID = v.ID
			join ApiKeyToBrand as akb on akb.brandID = vtb.brandID
			join ApiKey as ak on ak.id = akb.keyID
            && ak.api_key = ? && (vtb.brandID = ? or 0 = ?)`
	getBrands        = `select brandID from VideoNewToBrand where videoID = ?`
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
	getVideoParts      = `select partID from VideoJoin where videoID = ?`
	getVideoIdFromPart = `select videoID from VideoJoin where partID = ?`
	getChannel         = "SELECT ID, typeID, link, embedCode, foriegnID, dateAdded, dateModified, title, `desc` FROM Channel WHERE ID = ?"
	getCdn             = `SELECT ` + cdnFileFields + `, ` + cdnFileTypeFields + `, cf.dateModified FROM CdnFile AS cf LEFT JOIN CdnFileType AS cft ON cft.ID = cf.typeID WHERE cf.ID = ?`
	getCdnType         = `SELECT ID, mimeType, title, description FROM CdnFileType WHERE ID = ?`
	getAllCdnTypes     = `SELECT ID, mimeType, title, description FROM CdnFileType`
	getVideoType       = `SELECT vTypeID, name, icon FROM videoType WHERE vTypeID = ?`
	getChannelType     = `SELECT ID, name, description FROM ChannelType WHERE ID = ?`
	getAllChannelTypes = `SELECT ID, name, description FROM ChannelType `
)

//Base Video
func (v *Video) Get() error {
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
	if v != nil {
		go redis.Setex(redis_key, *v, redis.CacheTimeout)
	}
	return err
}

func (v *Video) GetVideoDetails() error {
	redis_key := "video:details:" + strconv.Itoa(v.ID)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &v)
		return err
	}

	brandChan := make(chan error)
	chanChan := make(chan error)
	cdnChan := make(chan error)
	partChan := make(chan error)

	err = v.Get()
	if err != nil {
		return err
	}

	go func() {
		err = v.GetBrands()
		if err != nil {
			brandChan <- err
		}
		brandChan <- nil

	}()
	go func() {
		chs, err := v.GetChannels()
		if err != nil {
			chanChan <- err
		}
		v.Channels = chs
		chanChan <- nil

	}()
	go func() {
		cdns, err := v.GetCdnFiles()
		if err != nil {
			cdnChan <- err
		}
		v.Files = cdns
		cdnChan <- nil

	}()

	go func() {
		err := v.GetParts()
		if err != nil {
			partChan <- err
		}
		partChan <- nil

	}()

	err = <-brandChan
	if err != nil {
		return err
	}
	err = <-chanChan
	if err != nil {
		return err
	}
	err = <-cdnChan
	if err != nil {
		return err
	}
	err = <-partChan
	if err != nil {
		return err
	}

	close(brandChan)
	close(chanChan)
	close(cdnChan)
	close(partChan)

	if v != nil {
		go redis.Setex(redis_key, v, redis.CacheTimeout)
	}
	return nil
}

func GetAllVideos(dtx *apicontext.DataContext) (vs Videos, err error) {
	data, err := redis.Get(AllVideosRedisKey + ":" + dtx.BrandString)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &vs)
		return
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllVideos)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
	if err != nil {
		return
	}

	ch := make(chan Videos)
	go populateVideos(rows, ch)
	vs = <-ch

	close(ch)

	if len(vs) == 0 {
		err = sql.ErrNoRows
		return
	}

	go redis.Setex(AllVideosRedisKey+":"+dtx.BrandString, vs, 86400)

	return
}

// // TODO - This is very slow...
// func GetPartVideos(partId int) (vs Videos, err error) {
// 	vs = make([]Video, 0)
// 	redis_key := "video:part:" + strconv.Itoa(partId)

// 	data, err := redis.Get(redis_key)
// 	if err == nil && len(data) > 0 {
// 		if err = json.Unmarshal(data, &vs); err == nil {
// 			return
// 		}
// 	}

// 	db, err := sql.Open("mysql", database.ConnectionString())
// 	if err != nil {
// 		return
// 	}
// 	defer db.Close()

// 	stmt, err := db.Prepare(getVideoIdFromPart)
// 	if err != nil {
// 		return
// 	}
// 	defer stmt.Close()

// 	res, err := stmt.Query(partId)
// 	if err != nil {
// 		return
// 	}

// 	var v Video
// 	for res.Next() {
// 		err = res.Scan(&v.ID)
// 		if err != nil {
// 			return
// 		}
// 		err = v.GetVideoDetails()
// 		if err != nil {
// 			return
// 		}
// 		vs = append(vs, v)
// 	}
// 	defer res.Close()

// 	if vs != nil {
// 		go redis.Setex(redis_key, vs, redis.CacheTimeout)
// 	}

// 	return vs, err
// }

func (v *Video) GetBrands() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getBrands)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Query(v.ID)
	if err != nil {
		return err
	}
	var b brand.Brand
	for res.Next() {
		err = res.Scan(&b.ID)
		if err != nil {
			return err
		}
		v.Brands = append(v.Brands, b)
	}
	return err
}

func (v *Video) GetChannels() (chs Channels, err error) {
	redis_key := "video:channels:" + strconv.Itoa(v.ID)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &chs)
		return
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare(getVideoChannels)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(v.ID)
	if err != nil {
		return
	}

	ch := make(chan Channels)
	go populateChannels(rows, ch)
	chs = <-ch

	if chs != nil {
		go redis.Setex(redis_key, chs, redis.CacheTimeout)
	}

	return
}

func (v *Video) GetParts() (err error) {
	redis_key := "video:parts:" + strconv.Itoa(v.ID)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &v.PartIds)
		return
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare(getVideoParts)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(v.ID)
	if err != nil {
		return
	}
	var i *int
	for rows.Next() {
		err = rows.Scan(&i)
		if err != nil {
			return err
		}
		if i != nil {
			v.PartIds = append(v.PartIds, *i)
		}
	}
	if len(v.PartIds) > 0 {
		go redis.Setex(redis_key, v.PartIds, redis.CacheTimeout)
	}
	return
}

func (v *Video) GetCdnFiles() (cdns CdnFiles, err error) {
	redis_key := "video:cdnFiles:" + strconv.Itoa(v.ID)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &cdns)
		return
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare(getVideoCdns)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(v.ID)
	if err != nil {
		return
	}

	ch := make(chan CdnFiles)
	go populateCdns(rows, ch)
	cdns = <-ch

	if cdns != nil {
		go redis.Setex(redis_key, cdns, redis.CacheTimeout)
	}

	return
}

func (c *Channel) Get() error {
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
	data, err := redis.Get(AllChannelsRedisKey)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &cs)
		return
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllChannels)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return
	}

	ch := make(chan Channels)
	go populateChannels(rows, ch)
	cs = <-ch

	if len(cs) == 0 {
		err = sql.ErrNoRows
		return
	}
	go redis.Setex(AllChannelsRedisKey, cs, 86400)
	return
}

func (c *CdnFile) Get() error {
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
	data, err := redis.Get(AllCdnFilesRedisKey)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &cs)
		return
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllCdnFiles)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return
	}

	ch := make(chan CdnFiles)
	go populateCdns(rows, ch)
	cs = <-ch
	if len(cs) == 0 {
		err = sql.ErrNoRows
		return
	}
	go redis.Setex(AllCdnFilesRedisKey, cs, 86400)
	return
}

func (c *CdnFileType) Get() error {
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
	data, err := redis.Get(AllCdnFileTypeRedisKey)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &cts)
		return
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllCdnTypes)
	if err != nil {
		return
	}
	defer stmt.Close()

	res, err := stmt.Query()
	if err != nil {
		return
	}
	var c CdnFileType
	var desc *string
	for res.Next() {
		err = res.Scan(
			&c.ID,
			&c.MimeType,
			&c.Title,
			&desc,
		)
		if err != nil {
			return
		}
		if desc != nil {
			c.Description = *desc
		}
		cts = append(cts, c)
	}
	defer res.Close()

	go redis.Setex(AllCdnFileTypeRedisKey, cts, 86400)
	return
}

func (c *VideoType) Get() error {
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
	data, err := redis.Get(AllVideoTypesRedisKey)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &vts)
		return
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllVideoTypes)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return
	}
	var vt VideoType
	var vName, vIcon *string
	for rows.Next() {
		err = rows.Scan(&vt.ID, &vName, &vIcon)
		if err != nil {
			return
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
	go redis.Setex(AllVideoTypesRedisKey, vts, 86400)
	return
}

func (c *ChannelType) Get() error {
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
	data, err := redis.Get(AllChannelTypesRedisKey)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &cts)
		return
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllChannelTypes)
	if err != nil {
		return
	}
	defer stmt.Close()
	var c ChannelType
	res, err := stmt.Query()
	if err != nil {
		return
	}
	for res.Next() {
		err = res.Scan(
			&c.ID,
			&c.Name,
			&c.Description,
		)
		if err != nil {
			return
		}
		cts = append(cts, c)
	}
	defer res.Close()
	go redis.Setex(AllChannelTypesRedisKey, cts, 86400)
	return cts, err
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
