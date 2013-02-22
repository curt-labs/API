package models

import (
	"../helpers/database"
	"strconv"
	"time"
)

type Review struct {
	Rating                           int
	Subject, ReviewText, Name, Email string
	CreatedDate                      time.Time
}

func (p *Part) GetReviews() error {
	db := database.Db

	rows, res, err := db.Query(partReviewStmt, p.PartId)
	if database.MysqlError(err) {
		return err
	}

	rating := res.Map("rating")
	subject := res.Map("subject")
	txt := res.Map("review_text")
	name := res.Map("name")
	email := res.Map("email")
	createdDate := res.Map("createdDate")

	var reviews []Review
	var ratingCounter int
	for _, row := range rows {
		date_add, _ := time.Parse("2006-01-02 15:04:01", row.Str(createdDate))
		r := Review{
			Rating:      row.Int(rating),
			Subject:     row.Str(subject),
			ReviewText:  row.Str(txt),
			Name:        row.Str(name),
			Email:       row.Str(email),
			CreatedDate: date_add,
		}
		reviews = append(reviews, r)

		ratingCounter = ratingCounter + r.Rating
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
