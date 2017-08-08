package products

import (
	// "github.com/curt-labs/API/models/customer"

	"time"
)

// Review ...
type Review struct {
	Rating      int        `bson:"rating" json:"rating" xml:"rating"`
	Subject     string     `bson:"subject" json:"subject" xml:"subject"`
	ReviewText  string     `bson:"review_text" json:"review_text" xml:"review_text"`
	Name        string     `bson:"name" json:"name" xml:"name"`
	Email       string     `bson:"email" json:"email" xml:"email"`
	CreatedDate *time.Time `bson:"created_date" json:"created_date" xml:"created_date"`
	Active      bool       `json:"active,omitempty" xml:"active,omitempty"`
	Approved    bool       `json:"approved,omitempty" xml:"approved,omitempty"`
}
