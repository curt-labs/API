package products

import (
	"github.com/curt-labs/GoAPI/models/customer/content"
)

type Content struct {
	ID          int                     `json:"id,omitempty" xml:"id,omitempty"`
	Text        string                  `json:"text,omitempty" xml:"text,omitempty"`
	ContentType custcontent.ContentType `json:"contentType,omitempty" xml:"contentType,omitempty"`
	UserID      string                  `json:"userId,omitempty" xml:"userId,omitempty"`
	Deleted     bool                    `json:"deleted,omitempty" xml:"deleted,omitempty"`
}
