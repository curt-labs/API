package news_model

import (
	"database/sql"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/pagination"
	"github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

type News struct {
	ID           int       `json:"id,omitempty" xml:"id,omitempty"`
	Title        string    `json:"title,omitempty" xml:"title,omitempty"`
	Lead         string    `json:"lead,omitempty" xml:"lead,omitempty"`
	Content      string    `json:"content,omitempty" xml:"content,omitempty"`
	PublishStart time.Time `json:"publishStart,omitempty" xml:"publishStart,omitempty"`
	PublishEnd   time.Time `json:"publishEnd,omitempty" xml:"publishEnd,omitempty"`
	Active       bool      `json:"active,omitempty" xml:"active,omitempty"`
	Slug         string    `json:"slug,omitempty" xml:"slug,omitempty"`
}
type Newses []News

var (
	getNews = `SELECT ni.newsItemID, ni.title, ni.lead, ni.content, ni.publishStart, ni.publishEnd, ni.active, ni.slug FROM NewsItem as ni 
									Join NewsItemToBrand as nib on nib.newsItemID = ni.newsItemID
									Join ApiKeyToBrand as akb on akb.brandID = nib.brandID
									Join ApiKey as ak on akb.keyID = ak.id
									where ni.newsItemID = ? && (ak.api_key = ? && (nib.brandID = ? OR 0=?))`
	getAll = `SELECT ni.newsItemID, ni.title, ni.lead, ni.content, ni.publishStart, ni.publishEnd, ni.active, ni.slug FROM NewsItem as ni
									Join NewsItemToBrand as nib on nib.newsItemID = ni.newsItemID
									Join ApiKeyToBrand as akb on akb.brandID = nib.brandID
									Join ApiKey as ak on akb.keyID = ak.id
									where (ak.api_key = ? && (nib.brandID = ? OR 0=?))
	`
	create        = `INSERT INTO NewsItem (title, lead, content, publishStart, publishEnd, active, slug) VALUES (?,?,?,?,?,?,?)`
	createToBrand = `INSERT INTO NewsItemToBrand (newsItemID, brandID) VALUES (?, ?)`
	update        = `UPDATE NewsItem SET title = ?, lead = ?, content = ?, publishStart = ?, publishEnd = ?, active = ?, slug = ? WHERE newsItemID = ?`
	deleteNews    = `DELETE FROM NewsItem WHERE newsItemID = ?`
	deleteToBrand = `DELETE FROM NewsItemToBrand WHERE newsItemID = ?`
	getTitles     = `SELECT ni.title FROM NewsItem as ni 
					Join NewsItemToBrand as nib on nib.newsItemID = ni.newsItemID
					Join ApiKeyToBrand as akb on akb.brandID = nib.brandID
					Join ApiKey as ak on akb.keyID = ak.id
					where (ak.api_key = ? && (nib.brandID = ? OR 0=?))
					`
	getLeads = `SELECT ni.lead FROM NewsItem as ni
					Join NewsItemToBrand as nib on nib.newsItemID = ni.newsItemID
					Join ApiKeyToBrand as akb on akb.brandID = nib.brandID
					Join ApiKey as ak on akb.keyID = ak.id
					where (ak.api_key = ? && (nib.brandID = ? OR 0=?))`
	search = `SELECT ni.newsItemID, ni.title, ni.lead, ni.content, ni.publishStart, ni.publishEnd, ni.active, ni.slug FROM NewsItem as ni
					Join NewsItemToBrand as nib on nib.newsItemID = ni.newsItemID
					Join ApiKeyToBrand as akb on akb.brandID = nib.brandID
					Join ApiKey as ak on akb.keyID = ak.id
					WHERE ni.title LIKE ? AND ni.lead LIKE ? AND ni.content LIKE ? AND ni.publishStart LIKE ? AND ni.publishEnd LIKE ? AND ni.active LIKE ? AND ni.slug LIKE ? &&
					(ak.api_key = ? && (nib.brandID = ? OR 0=?))
					`
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

func (n *News) Get(dtx *apicontext.DataContext) error {
	var err error
	// 1000th commit - Because I can.
	redis_key := "news:" + strconv.Itoa(n.ID) + ":" + dtx.BrandString
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &n)
		return err
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getNews)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Query(n.ID, dtx.APIKey, dtx.BrandID, dtx.BrandID)

	for res.Next() { //Time scanning creates odd driver issue using QueryRow
		res.Scan(&n.ID, &n.Title, &n.Lead, &n.Content, &n.PublishStart, &n.PublishEnd, &n.Active, &n.Slug)
	}
	if err != nil {
		return err
	}
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

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return fs, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getAll)
	if err != nil {
		return fs, err
	}
	defer stmt.Close()

	res, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
	for res.Next() {
		var n News
		res.Scan(&n.ID, &n.Title, &n.Lead, &n.Content, &n.PublishStart, &n.PublishEnd, &n.Active, &n.Slug)
		if err != nil {
			return fs, err
		}
		fs = append(fs, n)
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

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return l, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getTitles)
	if err != nil {
		return l, err
	}
	defer stmt.Close()

	res, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
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
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return l, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getLeads)
	if err != nil {
		return l, err
	}
	defer stmt.Close()

	res, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
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

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return l, err
	}
	defer db.Close()

	stmt, err := db.Prepare(search)
	if err != nil {
		return l, err
	}
	defer stmt.Close()

	res, err := stmt.Query("%"+title+"%", "%"+lead+"%", "%"+content+"%", "%"+publishStart+"%", "%"+publishEnd+"%", "%"+active+"%", "%"+slug+"%", dtx.APIKey, dtx.BrandID, dtx.BrandID)
	for res.Next() {
		var n News
		res.Scan(&n.ID, &n.Title, &n.Lead, &n.Content, &n.PublishStart, &n.PublishEnd, &n.Active, &n.Slug)
		fs = append(fs, n)
	}
	defer res.Close()
	l = pagination.Paginate(pageStr, resultsStr, fs)
	return l, err
}

func (n *News) Create(dtx *apicontext.DataContext) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(create)
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
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(createToBrand)
	_, err = stmt.Exec(n.ID, brand)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (n *News) Update(dtx *apicontext.DataContext) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(update)
	_, err = stmt.Exec(n.Title, n.Lead, n.Content, n.PublishStart, n.PublishEnd, n.Active, n.Slug, n.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (n *News) Delete(dtx *apicontext.DataContext) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deleteNews)
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
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deleteToBrand)
	_, err = stmt.Exec(n.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
