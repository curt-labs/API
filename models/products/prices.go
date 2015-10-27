package products

import (
	"time"
)

type Price struct {
	Id           int       `json:"id,omitempty" xml:"id,omitempty"`
	PartId       int       `json:"partId,omitempty" xml:"partId,omitempty"`
	Type         string    `json:"type,omitempty" xml:"type,omitempty"`
	Price        float64   `json:"price" xml:"price"`
	Enforced     bool      `json:"enforced,omitempty", xml:"enforced, omitempty"`
	DateModified time.Time `json:"dateModified,omitempty" xml:"dateModified,omitempty"`
}
