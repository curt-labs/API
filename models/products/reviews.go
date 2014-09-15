package products

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
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
	redis_key := fmt.Sprintf("part:%d:reviews", p.PartId)

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

	go redis.Setex(redis_key, p.Reviews, redis.CacheTimeout)

	return nil
}
