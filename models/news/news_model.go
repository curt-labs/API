package news_model

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/pagination"
	"github.com/curt-labs/API/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
)

type News struct {
	ID           int        `json:"id,omitempty" xml:"id,omitempty"`
	Title        string     `json:"title,omitempty" xml:"title,omitempty"`
	Lead         string     `json:"lead,omitempty" xml:"lead,omitempty"`
	Content      string     `json:"content,omitempty" xml:"content,omitempty"`
	PublishStart time.Time  `json:"publishStart,omitempty" xml:"publishStart,omitempty"`
	PublishEnd   time.Time  `json:"publishEnd,omitempty" xml:"publishEnd,omitempty"`
	Active       bool       `json:"active,omitempty" xml:"active,omitempty"`
	Slug         string     `json:"slug,omitempty" xml:"slug,omitempty"`
	Metadata     []Metadata `json:"metadata" xml:"metadata"`
}
type Newses []News

type Scanner interface {
	Scan(...interface{}) error
}

var (
	getNews = `SELECT ni.newsItemID, ni.title, ni.lead, ni.content, ni.publishStart, ni.publishEnd, ni.active, ni.slug
		FROM NewsItem AS ni WHERE ni.newsItemID = ?`
	getAll = `SELECT ni.newsItemID, ni.title, ni.lead, ni.content, ni.publishStart, ni.publishEnd, ni.active, ni.slug
		FROM NewsItem AS ni
		JOIN NewsItemToBrand AS nib ON nib.newsItemID = ni.newsItemID
		WHERE nib.brandID = ?`
	create        = `INSERT INTO NewsItem (title, lead, content, publishStart, publishEnd, active, slug) VALUES (?,?,?,?,?,?,?)`
	createToBrand = `INSERT INTO NewsItemToBrand (newsItemID, brandID) VALUES (?, ?)`
	update        = `UPDATE NewsItem SET title = ?, lead = ?, content = ?, publishStart = ?, publishEnd = ?, active = ?, slug = ? WHERE newsItemID = ?`
	deleteNews    = `DELETE FROM NewsItem WHERE newsItemID = ?`
	deleteToBrand = `DELETE FROM NewsItemToBrand WHERE newsItemID = ?`
	getTitles     = `SELECT ni.title FROM NewsItem AS ni JOIN NewsItemToBrand AS nib ON nib.newsItemID = ni.newsItemID WHERE nib.brandID = ?`
	getLeads = `SELECT ni.lead FROM NewsItem AS ni
		JOIN NewsItemToBrand AS nib ON nib.newsItemID = ni.newsItemID
		WHERE nib.brandID = ?`
	search = `SELECT ni.newsItemID, ni.title, ni.lead, ni.content, ni.publishStart, ni.publishEnd, ni.active, ni.slug FROM NewsItem AS ni
		JOIN NewsItemToBrand AS nib ON nib.newsItemID = ni.newsItemID
		WHERE ni.title LIKE ? AND ni.lead LIKE ? AND ni.content LIKE ? AND ni.publishStart LIKE ? AND ni.publishEnd LIKE ? AND
		ni.active LIKE ? AND ni.slug LIKE ? && nib.brandID = ?`
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

func (n *News) Get(dtx *apicontext.DataContext) error {
	var err error

	redis_key := "news:" + strconv.Itoa(n.ID) + ":" + dtx.BrandString
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &n)
		return err
	}

	err = database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(getNews)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRow(n.ID)
	item, err := scanItem(row)
	if err != nil {
		return err
	}

	n.copy(item)

	go redis.Setex(redis_key, n, 86400)

	return nil
}

func GetAll(dtx *apicontext.DataContext) (Newses, error) {
	var fs Newses
	var err error
	redis_key := "news" + ":" + dtx.BrandString
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &fs)
		return fs, err
	}

	err = database.Init()
	if err != nil {
		return fs, err
	}

	stmt, err := database.DB.Prepare(getAll)
	if err != nil {
		return fs, err
	}
	defer stmt.Close()

	res, err := stmt.Query(dtx.BrandID)
	for res.Next() {
		n, err := scanItem(res)
		if err == nil {
			fs = append(fs, n)
		}
	}
	defer res.Close()

	go redis.Setex(redis_key, fs, 86400)
	return fs, nil
}

func GetTitles(pageStr, resultsStr string, dtx *apicontext.DataContext) (pagination.Objects, error) {
	var err error
	var fs []interface{}
	var l pagination.Objects

	redis_key := "news:titles" + ":" + dtx.BrandString
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &l)
		return l, err
	}

	err = database.Init()
	if err != nil {
		return l, err
	}

	stmt, err := database.DB.Prepare(getTitles)
	if err != nil {
		return l, err
	}
	defer stmt.Close()

	res, err := stmt.Query(dtx.BrandID)
	for res.Next() {
		var f News
		res.Scan(&f.Title)
		fs = append(fs, f)
	}
	defer res.Close()
	l = pagination.Paginate(pageStr, resultsStr, fs)

	go redis.Setex(redis_key, l, 86400)
	return l, err
}

func GetLeads(pageStr, resultsStr string, dtx *apicontext.DataContext) (pagination.Objects, error) {
	var err error
	var fs []interface{}
	var l pagination.Objects
	redis_key := "news:leads" + ":" + dtx.BrandString
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &l)
		return l, err
	}

	err = database.Init()
	if err != nil {
		return l, err
	}

	stmt, err := database.DB.Prepare(getLeads)
	if err != nil {
		return l, err
	}
	defer stmt.Close()

	res, err := stmt.Query(dtx.BrandID)
	for res.Next() {
		var f News
		res.Scan(&f.Lead)
		fs = append(fs, f)
	}
	defer res.Close()
	l = pagination.Paginate(pageStr, resultsStr, fs)

	go redis.Setex(redis_key, l, 86400)
	return l, err
}

func Search(title, lead, content, publishStart, publishEnd, active, slug, pageStr, resultsStr string, dtx *apicontext.DataContext) (pagination.Objects, error) {
	var err error
	var l pagination.Objects
	var fs []interface{}

	err = database.Init()
	if err != nil {
		return l, err
	}

	stmt, err := database.DB.Prepare(search)
	if err != nil {
		return l, err
	}
	defer stmt.Close()

	res, err := stmt.Query("%"+title+"%", "%"+lead+"%", "%"+content+"%", "%"+publishStart+"%", "%"+publishEnd+"%", "%"+active+"%", "%"+slug+"%", dtx.BrandID)
	for res.Next() {
		n, err := scanItem(res)
		if err == nil {
			fs = append(fs, n)
		}
	}
	defer res.Close()
	l = pagination.Paginate(pageStr, resultsStr, fs)
	return l, err
}

func (n *News) Create(dtx *apicontext.DataContext) error {
	err := database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(create)
	if err != nil {
		return err
	}

	defer stmt.Close()
	res, err := stmt.Exec(n.Title, n.Lead, n.Content, n.PublishStart, n.PublishEnd, n.Active, n.Slug)
	if err != nil {
		tx.Rollback()
		return err
	}

	id, err := res.LastInsertId()
	n.ID = int(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	//createToBrand
	brands, err := dtx.GetBrandsFromKey()
	if err != nil {
		tx.Rollback()
		return err
	}
	bChan := make(chan int)
	go func() (err error) {
		if len(brands) > 0 {
			for _, brand := range brands {
				err = n.CreateJoinBrand(brand)
			}
		}
		bChan <- 1
		return err
	}()
	<-bChan

	tx.Commit()
	return nil
}

func (n *News) CreateJoinBrand(brand int) error {
	err := database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(createToBrand)
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(n.ID, brand)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (n *News) Update(dtx *apicontext.DataContext) error {
	err := database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(update)
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(n.Title, n.Lead, n.Content, n.PublishStart, n.PublishEnd, n.Active, n.Slug, n.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (n *News) Delete(dtx *apicontext.DataContext) error {
	err := database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(deleteNews)
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(n.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = n.DeleteJoinBrand()
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (n *News) DeleteJoinBrand() error {
	err := database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deleteToBrand)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(n.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (n *News) copy(item News) {
	n.ID = item.ID
	n.Title = item.Title
	n.Lead = item.Lead
	n.Content = item.Content
	n.PublishEnd = item.PublishEnd
	n.PublishStart = item.PublishStart
	n.Active = item.Active
	n.Slug = item.Slug
	n.Metadata = item.Metadata
}

func scanItem(res Scanner) (News, error) {
	var n News
	var id *int
	var title, lead, content, slug *string
	var pubStart, pubEnd *time.Time
	var active *bool
	err := res.Scan(
		&id,
		&title,
		&lead,
		&content,
		&pubStart,
		&pubEnd,
		&active,
		&slug)

	if err != nil || id == nil {
		return n, err
	}
	n.ID = *id
	if title != nil {
		n.Title = *title
	}
	if lead != nil {
		n.Lead = *lead
	}
	if content != nil {
		n.Content = *content
	}
	if slug != nil {
		n.Slug = *slug
	}
	if pubStart != nil {
		n.PublishStart = *pubStart
	}
	if pubEnd != nil {
		n.PublishEnd = *pubEnd
	}
	if active != nil {
		n.Active = *active
	}

	n.Metadata, err = GetMetadata(n.ID)

	return n, err
}
