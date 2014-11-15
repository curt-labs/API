package faq_model

import (
	"database/sql"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/pagination"
	"github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

type Faq struct {
	ID       int    `json:"id,omitempty" xml:"id,omitempty"`
	Question string `json:"question,omitempty" xml:"question,omitempty"`
	Answer   string `json:"answer,omitempty" xml:"answer,omitempty"`
}

type Faqs []Faq

type Pagination struct {
	TotalItems    int `json:"total_items" xml:"total_items"`
	ReturnedCount int `json:"returned_count" xml:"returned_count"`
	Page          int `json:"page" xml:"page"`
	PerPage       int `json:"per_page" xml:"per_page"`
	TotalPages    int `json:"total_pages" xml:"total_pages"`
}

var (
	getFaq       = "SELECT faqID, question, answer FROM FAQ WHERE faqID = ?"
	getAll       = "SELECT faqID, question, answer FROM FAQ"
	create       = "INSERT INTO FAQ (question, answer) VALUES (?,?)"
	update       = "UPDATE FAQ SET question = ?, answer = ? WHERE faqID = ?"
	deleteFaq    = "DELETE FROM FAQ WHERE faqID = ?"
	getQuestions = "SELECT question FROM FAQ"
	getAnswers   = "SELECT answer FROM FAQ"
	search       = "SELECT faqID, question, answer FROM FAQ WHERE question LIKE ? AND answer LIKE ? "
)

func GetAll() (Faqs, error) {
	var fs Faqs
	var err error
	redis_key := "faq"
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
		var f Faq
		res.Scan(&f.ID, &f.Question, &f.Answer)
		if err != nil {
			return fs, err
		}
		fs = append(fs, f)
	}
	defer res.Close()
	go redis.Setex(redis_key, fs, 86400)
	return fs, nil
}

func Search(question, answer, pageStr, resultsStr string) (pagination.Objects, error) {
	var err error
	var fs []interface{}
	var p pagination.Objects

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return p, err
	}
	defer db.Close()

	stmt, err := db.Prepare(search)
	if err != nil {
		return p, err
	}
	defer stmt.Close()

	res, err := stmt.Query("%"+question+"%", "%"+answer+"%")
	for res.Next() {
		var f Faq
		res.Scan(&f.ID, &f.Question, &f.Answer)
		fs = append(fs, f)
	}

	p = pagination.Paginate(pageStr, resultsStr, fs)
	return p, err
}

func (f *Faq) Get() error {
	var err error
	redis_key := "faq:" + strconv.Itoa(f.ID)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &f)
		return err
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getFaq)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(f.ID).Scan(&f.ID, &f.Question, &f.Answer)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return err
	}

	go redis.Setex(redis_key, f, redis.CacheTimeout)

	return nil
}

func (f *Faq) Create() error {
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
	res, err := stmt.Exec(f.Question, f.Answer)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err := res.LastInsertId()
	f.ID = int(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (f *Faq) Update() error {
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
	_, err = stmt.Exec(f.Question, f.Answer, f.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (f *Faq) Delete() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deleteFaq)
	_, err = stmt.Exec(f.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
