package models

import (
	"../helpers/database"
	"strconv"
	"strings"
	"time"
)

var (
	partReviewStmt = `select rating,subject,review_text,name,email,createdDate from Review
				where partID = ? and approved = 1 and active = 1`

	partReviewStmt_ByGroup = `select partID,rating,subject,review_text,name,email,createdDate from Review
				where partID IN ('%s') and approved = 1 and active = 1`
)

type Review struct {
	Rating                           int
	Subject, ReviewText, Name, Email string
	CreatedDate                      time.Time
}

func (p *Part) GetReviews() error {
	qry, err := database.Db.Prepare(partReviewStmt)
	if err != nil {
		return err
	}

	rows, res, err := qry.Exec(p.PartId)
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

func (lookup *Lookup) GetReviews() error {

	var ids []string
	for _, p := range lookup.Parts {
		ids = append(ids, strconv.Itoa(p.PartId))
	}

	rows, res, err := database.Db.Query(partReviewStmt_ByGroup, strings.Join(ids, ","))
	if database.MysqlError(err) {
		return err
	}

	partID := res.Map("partID")
	rating := res.Map("rating")
	subject := res.Map("subject")
	txt := res.Map("review_text")
	name := res.Map("name")
	email := res.Map("email")
	createdDate := res.Map("createdDate")

	reviews := make(map[int][]Review, len(lookup.Parts))

	for _, row := range rows {
		pId := row.Int(partID)
		date_add, _ := time.Parse("2006-01-02 15:04:01", row.Str(createdDate))
		r := Review{
			Rating:      row.Int(rating),
			Subject:     row.Str(subject),
			ReviewText:  row.Str(txt),
			Name:        row.Str(name),
			Email:       row.Str(email),
			CreatedDate: date_add,
		}
		reviews[pId] = append(reviews[pId], r)
	}

	for _, p := range lookup.Parts {
		p.Reviews = reviews[p.PartId]
	}

	return nil
}
