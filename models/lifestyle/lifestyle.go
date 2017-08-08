package lifestyle

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
)

type Lifestyle struct {
	ID          int       `json:"id,omitempty" xml:"id,omitempty`
	DateAdded   time.Time `json:"dateAdded,omitempty" xml:"dateAdded,omitempty`
	ParentID    int       `json:"parentID,omitempty" xml:"parentID,omitempty"`
	Name        string    `json:"name,omitempty" xml:"name,omitempty"`
	ShortDesc   string    `json:"shortDesc,omitempty" xml:"shortDesc,omitempty`
	LongDesc    string    `json:"longDesc,omitempty" xml:"longDesc,omitempty"`
	Image       string    `json:"image,omitempty" xml:"image,omitempty"`
	IsLifestyle int       `json:"isLifestyle,omitempty" xml:"isLifestyle,omitempty"`
	Sort        int       `json:"sort,omitempty" xml:"sort,omitempty"`
	Contents    Contents  `json:"contents,omitempty" xml:"contents,omitempty"`
	Towables    Towables  `json:"towables,omitempty" xml:"towables,omitempty`
}

type Lifestyles []Lifestyle

type Content struct {
	ID          int         `json:"id,omitempty" xml:"id,omitempty`
	UserID      int         `json:"userID,omitempty" xml:"userID,omitempty"`
	Text        string      `json:"content,omitempty" xml:"content,omitempty`
	ContentType ContentType `json:"contentType,omitempty" xml:"contentType,omitempty"`
	Deleted     bool        `json:"deleted,omitempty" xml:"deleted,omitempty"`
	PartID      int
}
type Contents []Content

type ContentType struct {
	ID   int    `json:"id,omitempty" xml:"id,omitempty`
	Name string `json:"name,omitempty" xml:"name,omitempty`
	HTML bool   `json:"html,omitempty" xml:"html,omitempty`
}
type Towable struct {
	ID         int    `json:"id,omitempty" xml:"id,omitempty`
	CatId      int    `json:"catId,omitempty" xml:"catId,omitempty`
	Name       string `json:"name,omitempty" xml:"name,omitempty"`
	ShortDesc  string `json:"shortDesc,omitempty" xml:"shortDesc,omitempty"`
	Image      string `json:"image,omitempty" xml:"image,omitempty"`
	HitchClass string `json:"hitchClass,omitempty" xml:"hitchClass,omitempty"`
	TW         int    `json:"TW,omitempty" xml:"TW,omitempty`
	GTW        int    `json:"GTW,omitempty" xml:"GTW,omitempty"`
	Message    string `json:"message,omitempty" xml:"message,omitempty"`
}
type Towables []Towable

var (
	getAllLifestyles = `select c.catID, c.catTitle, c.dateAdded, c.parentID,
							c.shortDesc, c.longDesc, c.image, c.isLifestyle,
							c.sort from Categories as c
							Join ApiKeyToBrand as akb on akb.brandID = c.brandID
							Join ApiKey as ak on akb.keyID = ak.id
							where c.isLifestyle = 1 && (ak.api_key = ? && (c.brandID = ? OR 0=?))
							order by c.sort`
	getLifestyle = `select
						c.catID, c.catTitle, c.dateAdded, c.parentID,
						c.shortDesc, c.longDesc, c.image, c.isLifestyle,
						c.sort
						from Categories as c
						Join ApiKeyToBrand as akb on akb.brandID = c.brandID
						Join ApiKey as ak on akb.keyID = ak.id
						where c.catID = ? && (ak.api_key = ? && (c.brandID = ? OR 0=?))
						limit 1`
	getLifestyleContent = `select ct.allowHTML, ct.type, c.text from Content as c
							join ContentBridge as cb on c.contentID = cb.contentID
							join ContentType as ct on c.cTypeID = ct.cTypeID
							where cb.catID = ?`
	getAllLifestyleContent = `select cb.catID, ct.allowHTML, ct.type, c.text from Content as c
							join ContentBridge as cb on c.contentID = cb.contentID
							join ContentType as ct on c.cTypeID = ct.cTypeID
							join Category as cat on cat.catID = cb.catID
							Join ApiKeyToBrand as akb on akb.brandID = cat.brandID
							Join ApiKey as ak on akb.keyID = ak.id
							where cb.catID > 0 && (ak.api_key = ? && (cat.brandID = ? OR 0=?))`
	getLifestyleTowables = `select
								t.trailerID, t.name, t.shortDesc, t.hitchClass, t.image, t.TW, t.GTW, t.message
								from Trailer as t
								join Lifestyle_Trailer as lt on t.trailerID = lt.trailerID
								where lt.catID = ?
								order by t.TW`

	getAllLifestyleTowables = `select
								t.trailerID, lt.catId, t.name, t.shortDesc, t.hitchClass, t.image, t.TW, t.GTW, t.message
								from Trailer as t
								join Lifestyle_Trailer as lt on t.trailerID = lt.trailerID
								join Category as cat on cat.catID = lt.catID
								Join ApiKeyToBrand as akb on akb.brandID = cat.brandID
								Join ApiKey as ak on akb.keyID = ak.id
								where (ak.api_key = ? && (cat.brandID = ? OR 0=?))
								order by t.TW`

	createLifestyle = `INSERT INTO Categories (dateAdded, parentID, catTitle, shortDesc, longDesc, image, isLifestyle, sort, brandID) VALUES (?,?,?,?,?,?,?,?,?)`
	updateLifestyle = `UPDATE Categories SET dateAdded = ?, parentID = ?, catTitle = ?, shortDesc = ?, longDesc = ?, image = ?, isLifestyle = ?, sort = ?, brandID = ? WHERE catID = ?`
	deleteLifestyle = `DELETE FROM Categories WHERE catID = ?`
	deleteContents  = `DELETE FROM ContentBridge WHERE catID = ?`
	deleteTowables  = `DELETE FROM Lifestyle_Trailer WHERE catID = ?`
	insertContent   = `INSERT INTO ContentBridge (catID,  contentID) VALUES (?,?)`
	insertTowable   = `INSERT INTO Lifestyle_Trailer (catID, trailerID) VALUES (?,?)`
	createContent   = `INSERT INTO Content (text, cTypeID, userID, deleted) VALUES (?,?,?,?)`
	createTowable   = `INSERT INTO Trailers (image, name, TW, GTW, hitchClass, shortDesc, message) VALUES (?,?,?,?,?,?,?)`
	getContent      = `SELECT c.contentID, c.text, c.cTypeID, c.userID, c.deleted, ct.type, ct.allowHTML FROM Content AS c LEFT JOIN ContentType AS ct ON ct.cTypeId = c.cTypeId WHERE c.contentID = ?`
	getTowable      = `SELECT trailerID, image, name, TW, GTW, hitchClass, shortDesc, message FROM Trailer WHERE trailerID = ?`
)

func GetAll(dtx *apicontext.DataContext) (ls Lifestyles, err error) {
	redis_key := "lifestyle:all:" + dtx.BrandString
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &ls)
		return ls, err
	}

	err = database.Init()
	if err != nil {
		return ls, err
	}

	stmt, err := database.DB.Prepare(getAllLifestyles)
	if err != nil {
		return ls, err
	}
	defer stmt.Close()
	//get content and towables
	cs, err := getAllContent(dtx)
	contentMap := cs.ToMap()
	ts, err := getAllTowables(dtx)
	towMap := ts.ToMap()

	res, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
	for res.Next() {
		var l Lifestyle
		err = res.Scan(&l.ID, &l.Name, &l.DateAdded, &l.ParentID, &l.ShortDesc, &l.LongDesc, &l.Image, &l.IsLifestyle, &l.Sort)
		if err != nil {
			return ls, err
		}
		//bind content and towables
		cChan := make(chan int)
		tChan := make(chan int)

		go func() {
			for _, val := range contentMap {
				if val.ID == l.ID {
					l.Contents = append(l.Contents, val)
				}
			}
			cChan <- 1
		}()

		go func() {
			for _, val := range towMap {
				if val.CatId == l.ID {
					l.Towables = append(l.Towables, val)
				}
			}
			tChan <- 1
		}()
		<-cChan
		<-tChan

		ls = append(ls, l)
	}
	defer res.Close()
	go redis.Setex(redis_key, ls, 86400)
	return ls, err
}

func (l *Lifestyle) Get(dtx *apicontext.DataContext) (err error) {
	redis_key := "lifestyle:get:" + strconv.Itoa(l.ID) + ":" + dtx.BrandString
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &l)
		return err
	}

	err = database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(getLifestyle)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(l.ID, dtx.APIKey, dtx.BrandID, dtx.BrandID).Scan(&l.ID, &l.Name, &l.DateAdded, &l.ParentID, &l.ShortDesc, &l.LongDesc, &l.Image, &l.IsLifestyle, &l.Sort)
	if err != nil {
		return err
	}
	err = l.GetContents()
	if err != nil {
		return err
	}
	err = l.GetTowables()
	if err != nil {
		return err
	}
	go redis.Setex(redis_key, l, 86400)
	return nil
}

func getAllContent(dtx *apicontext.DataContext) (cs Contents, err error) {
	err = database.Init()
	if err != nil {
		return cs, err
	}

	stmt, err := database.DB.Prepare(getAllLifestyleContent)
	if err != nil {
		return cs, err
	}
	defer stmt.Close()

	res, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
	for res.Next() {
		var c Content
		err = res.Scan(&c.ID, &c.ContentType.HTML, &c.ContentType.Name, &c.Text)
		if err != nil {
			return cs, err
		}
		cs = append(cs, c)
	}
	defer res.Close()
	return cs, err
}

func getAllTowables(dtx *apicontext.DataContext) (ts Towables, err error) {
	err = database.Init()
	if err != nil {
		return ts, err
	}

	stmt, err := database.DB.Prepare(getAllLifestyleTowables)
	if err != nil {
		return ts, err
	}
	defer stmt.Close()
	res, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
	for res.Next() {
		var t Towable
		err = res.Scan(&t.ID, &t.CatId, &t.Name, &t.ShortDesc, &t.HitchClass, &t.Image, &t.TW, &t.GTW, &t.Message)
		if err != nil {
			return ts, err
		}
		ts = append(ts, t)
	}
	defer res.Close()
	return ts, err
}

func (l *Lifestyle) GetContents() (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(getLifestyleContent)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Query(l.ID)
	for res.Next() {
		var c Content
		err = res.Scan(&c.ContentType.HTML, &c.ContentType.Name, &c.Text)
		if err != nil {
			return err
		}
		l.Contents = append(l.Contents, c)
	}
	defer res.Close()
	return err
}

func (l *Lifestyle) GetTowables() (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(getLifestyleTowables)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Query(l.ID)
	for res.Next() {
		var t Towable
		err = res.Scan(&t.ID, &t.Name, &t.ShortDesc, &t.HitchClass, &t.Image, &t.TW, &t.GTW, &t.Message)
		if err != nil {
			return err
		}
		l.Towables = append(l.Towables, t)
	}
	defer res.Close()
	return err
}

func (l *Lifestyle) Create(dtx *apicontext.DataContext) (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()

	stmt, err := tx.Prepare(createLifestyle)
	if err != nil {
		return err
	}
	defer stmt.Close()
	l.DateAdded = time.Now()
	res, err := stmt.Exec(l.DateAdded, l.ParentID, l.Name, l.ShortDesc, l.LongDesc, l.Image, l.IsLifestyle, l.Sort, dtx.BrandID)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err := res.LastInsertId()
	l.ID = int(id)
	err = tx.Commit()

	cChan := make(chan int)
	tChan := make(chan int)

	//insert content and/or joins
	go func() {
		cChan <- 1
		for _, content := range l.Contents {
			err = content.Get()
			if err != nil {
				err = content.Create()
				if err != nil {
					return
				}
			}
			err = content.insertContent(*l)
			if err != nil {
				return
			}
		}
	}()

	//insert towable and/or joins
	go func() {
		tChan <- 1
		for _, towable := range l.Towables {
			err = towable.Get()
			if err != nil {
				err = towable.Create()
				if err != nil {
					return
				}
			}
			err = towable.insertTowable()
			if err != nil {
				return
			}
		}

	}()
	<-cChan
	<-tChan

	return err
}

func (l *Lifestyle) Update(dtx *apicontext.DataContext) (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()

	stmt, err := tx.Prepare(updateLifestyle)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(l.DateAdded, l.ParentID, l.Name, l.ShortDesc, l.LongDesc, l.Image, l.IsLifestyle, l.Sort, dtx.BrandID, l.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()

	//delete content and towable joins
	err = l.deleteContents()
	err = l.deleteTowables()
	if err != nil {
		return err
	}

	//insert content and/or joins
	for _, content := range l.Contents {
		err = content.Get()
		if err != nil {
			err = content.Create()
			if err != nil {
				return err
			}
		}
		err = content.insertContent(*l)
		if err != nil {
			return err
		}
	}
	//insert towable and/or joins
	for _, towable := range l.Towables {
		err = towable.Get()
		if err != nil {
			err = towable.Create()
			if err != nil {
				return err
			}
		}
		err = towable.insertTowable()
		if err != nil {
			return err
		}
	}

	redis_key := "lifestyle:get:" + strconv.Itoa(l.ID) + ":" + dtx.BrandString
	go redis.Setex(redis_key, l, redis.CacheTimeout)
	return err
}

func (l *Lifestyle) Delete() (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()

	stmt, err := tx.Prepare(deleteLifestyle)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(l.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	err = l.deleteContents()
	err = l.deleteTowables()
	return err
}

func (l *Lifestyle) deleteContents() (err error) {
	if l.ID == 0 {
		return errors.New("Lifestyle content ID is zero.")
	}

	err = database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()

	stmt, err := tx.Prepare(deleteContents)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(l.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (l *Lifestyle) deleteTowables() (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()

	stmt, err := tx.Prepare(deleteTowables)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(l.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (c *Content) insertContent(l Lifestyle) (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()

	stmt, err := tx.Prepare(insertContent)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(l.ID, c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (t *Towable) insertTowable() (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()

	stmt, err := tx.Prepare(insertTowable)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(t.CatId, t.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (c *Content) Get() (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(getContent)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(c.ID).Scan(&c.ID, &c.Text, &c.ContentType.ID, &c.UserID, &c.Deleted, &c.ContentType.Name, &c.ContentType.HTML)

	return err
}

func (t *Towable) Get() (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(getTowable)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(t.ID).Scan(&t.ID, &t.Image, &t.Name, &t.TW, &t.GTW, &t.HitchClass, &t.ShortDesc, &t.Message)

	return err
}

func (c *Content) Create() (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()

	stmt, err := tx.Prepare(createContent)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(c.Text, c.ContentType.ID, c.UserID, c.Deleted)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err := res.LastInsertId()
	c.ID = int(id)
	return err
}

func (t *Towable) Create() (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()

	stmt, err := tx.Prepare(createTowable)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(t.Image, t.Name, t.TW, t.GTW, t.HitchClass, t.ShortDesc, t.Message)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err := res.LastInsertId()
	t.ID = int(id)
	err = tx.Commit()
	return err
}
