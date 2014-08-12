package part

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

var (
	partReviewStmt = `select rating,subject,review_text,name,email,createdDate from Review
				where partID = ? and approved = 1 and active = 1`

	partReviewStmt_ByGroup = `select partID,rating,subject,review_text,name,email,createdDate from Review
				where partID IN (%s) and approved = 1 and active = 1`
)

type Review struct {
	Rating                           int
	Subject, ReviewText, Name, Email string
	CreatedDate                      time.Time
}

func (p *Part) GetReviews() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	qry, err := db.Prepare(partReviewStmt)
	if err != nil {
		return err
	}
	defer qry.Close()

	rows, err := qry.Query(p.PartId)
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

	p.Reviews = reviews
	if len(reviews) > 0 {
		avg_str := strconv.Itoa(ratingCounter / len(reviews))
		p.AverageReview, _ = strconv.ParseFloat(avg_str, 64)
	} else {
		p.AverageReview = 0
	}

	return nil
}
