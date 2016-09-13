package faq_model

import (
	"database/sql"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/pagination"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllStmt = `select F.faqID, F.question, F.answer, F.brandID
                  from FAQ F
                  join ApiKeyToBrand as AKB on AKB.brandID = F.brandID
                  join ApiKey AK on AK.id = AKB.keyID
                  where AK.api_key = ? && (F.brandID = ? || 0 = ?)`
	getFaqStmt = `select F.faqID, F.question, F.answer, F.brandID
                  from FAQ F
                  join ApiKeyToBrand as AKB on AKB.brandID = F.brandID
                  join ApiKey AK on AK.id = AKB.keyID
                  where AK.api_key = ? && F.brandID = ? && F.faqID = ?`
	searchFaqStmt = `select F.faqID, F.question, F.answer, F.brandID
                     from FAQ F
                     join ApiKeyToBrand as AKB on AKB.brandID = F.brandID
                     join ApiKey AK on AK.id = AKB.keyID
                     where AK.api_Key = ? && F.brandID = ? && question like ? && answer like ?`
	createFaqStmt = `insert into FAQ(question,answer,brandID) values(?,?,?)`
	updateFaqStmt = `update FAQ set question = ?, answer = ?, brandID = ? where faqID = ?`
	deleteFaqStmt = `delete from FAQ where faqID = ?`
)

type Faqs []Faq
type Faq struct {
	ID       int    `json:"id,omitempty" xml:"id,omitempty"`
	BrandID  int    `json:"brandId,omitempty" xml:"brandId,omitempty"`
	Question string `json:"question,omitempty" xml:"question,omitempty"`
	Answer   string `json:"answer,omitempty" xml:"answer,omitempty"`
}

type Pagination struct {
	TotalItems    int `json:"total_items" xml:"total_items"`
	ReturnedCount int `json:"returned_count" xml:"returned_count"`
	Page          int `json:"page" xml:"page"`
	PerPage       int `json:"per_page" xml:"per_page"`
	TotalPages    int `json:"total_pages" xml:"total_pages"`
}

func GetAll(dtx *apicontext.DataContext) (Faqs, error) {
	var fs Faqs
	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return fs, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllStmt)
	if err != nil {
		return fs, err
	}
	defer stmt.Close()

	res, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
	for res.Next() {
		var f Faq
		res.Scan(&f.ID, &f.Question, &f.Answer, &f.BrandID)
		if err != nil {
			return fs, err
		}
		fs = append(fs, f)
	}
	defer res.Close()

	return fs, nil
}

func Search(dtx *apicontext.DataContext, question, answer, pageStr, resultsStr string) (pagination.Objects, error) {
	var err error
	var fs []interface{}
	var p pagination.Objects

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return p, err
	}
	defer db.Close()

	stmt, err := db.Prepare(searchFaqStmt)
	if err != nil {
		return p, err
	}
	defer stmt.Close()

	res, err := stmt.Query(dtx.APIKey, dtx.BrandID, "%"+question+"%", "%"+answer+"%")
	for res.Next() {
		var f Faq
		res.Scan(&f.ID, &f.Question, &f.Answer, &f.BrandID)
		fs = append(fs, f)
	}

	p = pagination.Paginate(pageStr, resultsStr, fs)
	return p, err
}

func (f *Faq) Get(dtx *apicontext.DataContext) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getFaqStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRow(dtx.APIKey, dtx.BrandID, f.ID)
	err = row.Scan(&f.ID, &f.Question, &f.Answer, &f.BrandID)

	if err != nil {
		return err
	}

	return nil
}

func (f *Faq) Create(dtx *apicontext.DataContext) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(createFaqStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(f.Question, f.Answer, dtx.BrandID)
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

func (f *Faq) Update(dtx *apicontext.DataContext) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(updateFaqStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(f.Question, f.Answer, dtx.BrandID, f.ID)
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
	stmt, err := tx.Prepare(deleteFaqStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(f.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
