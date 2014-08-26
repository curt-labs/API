package news_model

import (
	"database/sql"
	"encoding/json"
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
	getNews    = "SELECT newsItemID, title, lead, content, publishStart, publishEnd, active, slug FROM NewsItem WHERE newsItemID= ?"
	getAll     = "SELECT newsItemID, title, lead, content, publishStart, publishEnd, active, slug FROM NewsItem"
	create     = "INSERT INTO NewsItem (title, lead, content, publishStart, publishEnd, active, slug) VALUES (?,?,?,?,?,?,?)"
	update     = "UPDATE NewsItem SET title = ?, lead = ?, content = ?, publishStart = ?, publishEnd = ?, active = ?, slug = ? WHERE newsItemID = ?"
	deleteNews = "DELETE FROM NewsItem WHERE newsItemID = ?"
	getTitles  = "SELECT title FROM NewsItem"
	getLeads   = "SELECT lead FROM NewsItem"
	search     = "SELECT newsItemID, title, lead, content, publishStart, publishEnd, active, slug FROM NewsItem WHERE title LIKE ? AND lead LIKE ? AND content LIKE ? AND publishStart LIKE ? AND publishEnd LIKE ? AND active LIKE ? AND slug LIKE ?"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

func (n *News) Get() error {
	var err error

	redis_key := "news:" + strconv.Itoa(n.ID)
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
	res, err := stmt.Query(n.ID)

	for res.Next() { //Time scanning creates odd driver issue using QueryRow
		res.Scan(&n.ID, &n.Title, &n.Lead, &n.Content, &n.PublishStart, &n.PublishEnd, &n.Active, &n.Slug)
	}
	if err != nil {
		return err
	}
	go redis.Setex(redis_key, n, 86400)
	return nil
}

func GetAll() (Newses, error) {
	var fs Newses
	var err error
	redis_key := "news"
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

	res, err := stmt.Query()
	for res.Next() {
		var n News
		res.Scan(&n.ID, &n.Title, &n.Lead, &n.Content, &n.PublishStart, &n.PublishEnd, &n.Active, &n.Slug)
		if err != nil {
			return fs, err
		}
		fs = append(fs, n)
	}
	go redis.Setex(redis_key, fs, 86400)
	return fs, nil
}

func GetTitles(pageStr, resultsStr string) (pagination.Objects, error) {
	var err error
	var fs []interface{}
	var l pagination.Objects

	redis_key := "news:titles"
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

	res, err := stmt.Query()
	for res.Next() {
		var f News
		res.Scan(&f.Title)
		fs = append(fs, f)
	}
	l = pagination.Paginate(pageStr, resultsStr, fs)

	go redis.Setex(redis_key, l, 86400)
	return l, err
}

func GetLeads(pageStr, resultsStr string) (pagination.Objects, error) {
	var err error
	var fs []interface{}
	var l pagination.Objects
	redis_key := "news:leads"
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

	res, err := stmt.Query()
	for res.Next() {
		var f News
		res.Scan(&f.Lead)
		fs = append(fs, f)
	}
	l = pagination.Paginate(pageStr, resultsStr, fs)

	go redis.Setex(redis_key, l, 86400)
	return l, err
}

func Search(title, lead, content, publishStart, publishEnd, active, slug, pageStr, resultsStr string) (pagination.Objects, error) {
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

	res, err := stmt.Query("%"+title+"%", "%"+lead+"%", "%"+content+"%", "%"+publishStart+"%", "%"+publishEnd+"%", "%"+active+"%", "%"+slug+"%")
	for res.Next() {
		var n News
		res.Scan(&n.ID, &n.Title, &n.Lead, &n.Content, &n.PublishStart, &n.PublishEnd, &n.Active, &n.Slug)
		fs = append(fs, n)
	}
	l = pagination.Paginate(pageStr, resultsStr, fs)
	return l, err
}

func (n *News) Create() error {
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

	tx.Commit()
	return nil
}

func (n *News) Update() error {
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

func (n *News) Delete() error {
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
	tx.Commit()
	return nil
}
