package products

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
	"github.com/curt-labs/GoAPI/models/customer"
	_ "github.com/go-sql-driver/mysql"
	// "log"
	"strconv"
	"time"
)

var (
	activeApprovedReviews = `select rating,subject,review_text,name,email,createdDate from Review
				where partID = ? and approved = 1 and active = 1`

	// partReviewStmt_ByGroup = `select partID,rating,subject,review_text,name,email,createdDate from Review
	// 			where partID IN (%s) and approved = 1 and active = 1`
	getAllReviews = `SELECT reviewID, partID, rating, subject, review_text, name, email, active, approved, createdDate, cust_id FROM Review`
	getReview     = `SELECT reviewID, partID, rating, subject, review_text, name, email, active, approved, createdDate, cust_id FROM Review WHERE reviewID = ?`
	createReview  = `INSERT INTO Review (partID, rating, subject, review_text, name, email, active, approved, createdDate, cust_id) VALUES (?,?,?,?,?,?,?,?,?,?)`
	updateReview  = `UPDATE Review SET partID = ?, rating = ?, subject = ?, review_text = ?, name = ?, email = ?, active = ?, approved = ?, createdDate = ?, cust_id = ? WHERE reviewID = ?`
	deleteReview  = `DELETE FROM Review WHERE reviewID = ?`
	deleteReviews = `DELETE FROM Review WHERE partID = ?`
)

type Review struct {
	Id          int               `json:"id,omitempty" xml:"id,omitempty"`
	PartID      int               `json:"partId,omitempty" xml:"partId,omitempty"`
	Rating      int               `json:"rating,omitempty" xml:"rating,omitempty"`
	Subject     string            `json:"subject,omitempty" xml:"subject,omitempty"`
	ReviewText  string            `json:"reviewText,omitempty" xml:"reviewText,omitempty"`
	Name        string            `json:"name,omitempty" xml:"name,omitempty"`
	Email       string            `json:"email,omitempty" xml:"email,omitempty"`
	Active      bool              `json:"active,omitempty" xml:"active,omitempty"`
	Approved    bool              `json:"approved,omitempty" xml:"approved,omitempty"`
	CreatedDate time.Time         `json:"createdDate,omitempty" xml:"createdDate,omitempty"`
	Customer    customer.Customer `json:"customer,omitempty" xml:"customer,omitempty"`
}
type Reviews []Review

//gets kosher reviews by part
func (p *Part) GetActiveApprovedReviews() error {
	redis_key := fmt.Sprintf("part:%d:%d:reviews", p.BrandID, p.ID)

	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		if err = json.Unmarshal(data, &p.Reviews); err != nil {
			return nil
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	qry, err := db.Prepare(activeApprovedReviews)
	if err != nil {
		return err
	}
	defer qry.Close()

	rows, err := qry.Query(p.ID)
	if err != nil {
		return err
	}

	var reviews []Review
	var ratingCounter int
	for rows.Next() {
		var r Review
		err = rows.Scan(
			&r.Rating,
			&r.Subject,
			&r.ReviewText,
			&r.Name,
			&r.Email,
			&r.CreatedDate)
		if err == nil {
			reviews = append(reviews, r)
			ratingCounter = ratingCounter + r.Rating
		}
	}
	defer rows.Close()

	p.Reviews = reviews
	if len(reviews) > 0 {
		avg_str := strconv.Itoa(ratingCounter / len(reviews))
		p.AverageReview, _ = strconv.ParseFloat(avg_str, 64)
	} else {
		p.AverageReview = 0
	}

	go redis.Setex(redis_key, p.Reviews, redis.CacheTimeout)

	return nil
}

//get all reveiws, ever
func GetAllReviews(dtx *apicontext.DataContext) (revs Reviews, err error) {
	redis_key := fmt.Sprintf("reviews:%s", dtx.BrandString)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &revs)
		return revs, err
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return revs, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllReviews)
	if err != nil {
		return revs, err
	}
	defer stmt.Close()

	res, err := stmt.Query()

	var subject, text, name, email *string

	for res.Next() {
		var r Review
		err = res.Scan(&r.Id, &r.PartID, &r.Rating, &subject, &text, &name, &email, &r.Active, &r.Approved, &r.CreatedDate, &r.Customer.Id)
		if err != nil {
			return revs, err
		}
		if subject != nil {
			r.Subject = *subject
		}
		if text != nil {
			r.ReviewText = *text
		}
		if name != nil {
			r.Name = *name
		}
		if email != nil {
			r.Email = *email
		}

		revs = append(revs, r)
	}
	defer res.Close()
	if dtx.BrandString != "" {
		go redis.Setex(redis_key, revs, 86400)
	}
	return revs, err
}

func (r *Review) Get(dtx *apicontext.DataContext) (err error) {
	redis_key := fmt.Sprintf("reviews:%d:%s", r.Id, dtx.BrandString)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &r)
		return err
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getReview)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var subject, text, name, email *string
	var partId *int

	err = stmt.QueryRow(r.Id).Scan(&r.Id, &partId, &r.Rating, &subject, &text, &name, &email, &r.Active, &r.Approved, &r.CreatedDate, &r.Customer.Id)
	if subject != nil {
		r.Subject = *subject
	}
	if text != nil {
		r.ReviewText = *text
	}
	if name != nil {
		r.Name = *name
	}
	if email != nil {
		r.Email = *email
	}
	if partId != nil {
		r.PartID = *partId
	}

	if err != nil {
		return err
	}
	if dtx.BrandString != "" {
		go redis.Setex(redis_key, r, 86400)
	}
	return nil
}

func (r *Review) Create(dtx *apicontext.DataContext) (err error) {
	go redis.Delete("reviews:" + dtx.BrandString)
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()

	stmt, err := tx.Prepare(createReview)
	if err != nil {
		return err
	}
	defer stmt.Close()

	r.CreatedDate = time.Now()
	res, err := stmt.Exec(r.PartID, r.Rating, r.Subject, r.ReviewText, r.Name, r.Email, r.Active, r.Approved, r.CreatedDate, r.Customer.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err := res.LastInsertId()
	r.Id = int(id)
	err = tx.Commit()

	return err
}

func (r *Review) Update(dtx *apicontext.DataContext) (err error) {
	go redis.Delete("reviews:" + dtx.BrandString)
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()

	stmt, err := tx.Prepare(updateReview)
	if err != nil {

		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(r.PartID, r.Rating, r.Subject, r.ReviewText, r.Name, r.Email, r.Active, r.Approved, r.CreatedDate, r.Customer.Id, r.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

func (r *Review) Delete(dtx *apicontext.DataContext) (err error) {
	go redis.Delete("reviews:" + dtx.BrandString)
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()

	stmt, err := tx.Prepare(deleteReview)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(r.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()

	return err
}
func (r *Review) DeletebyPart(dtx *apicontext.DataContext) (err error) {
	go redis.Delete("reviews:" + dtx.BrandString)
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()

	stmt, err := tx.Prepare(deleteReviews)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(r.PartID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()

	return err
}
