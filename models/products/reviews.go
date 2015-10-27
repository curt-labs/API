package products

import (
	"github.com/curt-labs/GoAPI/models/customer"

	"time"
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
