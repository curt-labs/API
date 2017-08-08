package video

import (
	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/redis"
	"github.com/curt-labs/API/models/brand"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"encoding/json"
	"strconv"
	"time"
)

// Video is a reresentation of a video. It contains information about the video itself, as well as any associated files.
type Video struct {
	ID           int          `json:"id,omitempty" xml:"id,omitempty"`
	Title        string       `json:"title, omitempty" xml:"title,omitempty"`
	SubjectType  string       `bson:"subject_type" json:"subject_type" xml:"subject_type"`
	VideoType    VideoType    `json:"videoType,omitempty" xml:"videoType,omitempty"`
	Description  string       `bson:"description" json:"description" xml:"description"`
	DateAdded    time.Time    `bson:"date_added" json:"date_added" xml:"date_added"`
	DateModified time.Time    `bson:"date_modified" json:"date_modified" xml:"date_modified"`
	Thumbnail    string       `bson:"thumb_nail" json:"thumbnail" xml:"thumbnail"`
	Channels     []Channel    `bson:"channels" json:"channel" xml:"channel"`
	Files        []CdnFile    `bson:"files" json:"cdn_file" xml:"cdn_file"`
	IsPrimary    bool         `json:"isPrimary,omitempty" xml:"isPrimary,omitempty"`
	CategoryIds  []int        `json:"categoryIds,omitempty" xml:"categoryIds,omitempty"`
	PartIds      []int        `json:"partIds,omitempty" xml:"partIds,omitempty"`
	WebsiteId    int          `json:"websiteId,omitempty" xml:"websiteId,omitempty"`
	Brands       brand.Brands `json:"brands,omitempty" xml:"brands,omitempty"`
}

// Videos is just an easier type to work with than using an array of video types.
type Videos []Video

// A Channel type is typicaly the information associated to a online video file such as youtube, vimeo, etc.
type Channel struct {
	ID           int         `json:"id,omitempty" xml:"id,omitempty"`
	Type         ChannelType `json:"type,omitempty" xml:"type,omitempty"`
	Link         string      `bson:"link" json:"link" xml:"link"`
	EmbedCode    string      `bson:"embed_code" json:"embed_code" xml:"embed_code"`
	ForiegnID    string      `bson:"foreign_id" json:"foreign_id" xml:"foreign_id"`
	DateAdded    time.Time   `bson:"date_added" json:"date_added" xml:"date_added"`
	DateModified time.Time   `bson:"date_modified" json:"date_modified" xml:"date_modified"`
	Title        string      `bson:"title" json:"title" xml:"title,attr"`
	Description  string      `bson:"description" json:"description" xml:"description"`
	Duration     string      `bson:"duration" json:"duration" xml:"duration"`
}

// Channels is just an easier type to work with than using an array of Channel types.
type Channels []Channel

// ChannelType is a type of Channel. Channels are online videos, and they have different types such as youtube, vimeo, dailymotion, etc.
type ChannelType struct {
	ID          int    `json:"id,omitempty" xml:"id,omitempty"`
	Name        string `json:"name,omitempty" xml:"name,omitempty"`
	Description string `json:"description,omitempty" xml:"description,omitempty"`
}

// CdnFile or CDN file is a video that is hosted on a content delivery network.
// This is reference to an actual video file rather than an online video type(Channel)
// These CdnFiles are typically used for HTML5 Videos.
type CdnFile struct {
	ID           int         `json:"id,omitempty" xml:"id,omitempty"`
	Type         CdnFileType `bson:"type" json:"type" xml:"type"`
	Path         string      `bson:"path" json:"path" xml:"path"`
	Bucket       string      `bson:"bucket" json:"bucket" xml:"bucket"`
	ObjectName   string      `bson:"object_name" json:"object_name" xml:"object_name"`
	FileSize     string      `bson:"file_size" json:"file_size" xml:"file_size"`
	DateAdded    time.Time   `bson:"date_added" json:"date_added" xml:"date_added"`
	DateModified time.Time   `bson:"date_modified" json:"date_modified" xml:"date_modified"`
	LastUploaded string      `bson:"date_uploaded" json:"date_uploaded" xml:"date_uploaded"`
}

// CdnFiles is just an easier type to work with than using an array of CdnFiles.
type CdnFiles []CdnFile

// CdnFile type specifies what file type the file is. Some examples might be .ogg, .mp4, .avi, etc.
type CdnFileType struct {
	ID          int    `json:"id,omitempty" xml:"id,omitempty"`
	MimeType    string `bson:"mime_type" json:"mime_type" xml:"mime_type"`
	Title       string `bson:"title" json:"title" xml:"title,attr"`
	Description string `bson:"description" json:"description" xml:"description"`
}

// Video Type specifies what kind of video it is. Some examples might be, Product video, Howto, or Instructional Video.
type VideoType struct {
	ID   int    `json:"id,omitempty" xml:"id,omitempty"`
	Name string `bson:"name" json:"name" xml:"name"`
	Icon string `bson:"icon" json:"icon" xml:"icon"`
}

// Videos can be associated to categories. This is the most basic information about a category.
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
	getAllVideos = `SELECT ` + videoFields + `, ` + videoTypeFields + ` FROM VideoNew AS v
			LEFT JOIN videoType AS vt ON vt.vTypeID = v.subjectTypeID
			JOIN VideoNewToBrand AS vtb ON vtb.videoID = v.ID
		WHERE vtb.brandID = ?`
	getBrands        = `select brandID from VideoNewToBrand where videoID = ?`
	getAllCdnFiles   = `SELECT ` + cdnFileFields + `,` + cdnFileTypeFields + ` FROM CdnFile AS cf LEFT JOIN CdnFileType AS cft ON cft.ID = cf.typeID `
	getAllChannels   = `SELECT ` + channelFields + `, ` + channelTypeFields + ` FROM Channel AS c LEFT JOIN ChannelType AS ct ON ct.ID = c.typeID `
	getAllVideoTypes = `SELECT vt.vTypeID, ` + videoTypeFields + ` FROM videoType AS vt`
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

// Retrieves a base video file. This does not grab all the associated channels or CDN files.
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

// GetVideoDetails grabs a video's more advance information such as, Brands, CDN files, associated channels(youtube videos), and any products associated with the video.
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

// GetAllVideos This grabs all the videos given a certain Brand. Videos are  Base Videos and do not have advanced information.
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

// GetBrands Gets all the brands associated to a specific base video.
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

// GetChannels Gets all the Video's channels.
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

// GetParts Gets all the Video's associated products.
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
	defer rows.Close()

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

// GetCdnFiles Gets all of the CdnFiles for the specific video.
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

// Get Get a gven Channel
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

// GetAllChannels Retrieves all Channels from the DB
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

// Get Retrieves a given CdnFile
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

// GetAllCdnFiles Retrieves all CdnFiles
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

// Get Retrieves a given CdnFileType
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

// GetAllCdnFileTypes Retrieves all CdnFileTypes
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

// Get Retrieves a given VideoType
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

// GetAllVideoTypes Retrieves all VideoTypes
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

// Get Retrieves a given ChannelType
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

// GetAllChannelTypes Retrieves all ChannelType
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

// Populates a video + type
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
