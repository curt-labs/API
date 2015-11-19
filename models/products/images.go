package products

import (
	"net/url"
)

type Image struct {
	ID     int      `json:"id,omitempty" xml:"id,omitempty"`
	Size   string   `json:"size,omitempty" xml:"size,omitempty"`
	Sort   string   `json:"sort,omitempty" xml:"sort,omitempty"`
	Height int      `json:"height,omitempty" xml:"height,omitempty"`
	Width  int      `json:"width,omitempty" xml:"width,omitempty"`
	Path   *url.URL `json:"path,omitempty" xml:"path,omitempty"`
	PartID int      `json:"partId,omitempty" xml:"partId,omitempty"`
}
