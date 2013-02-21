package models

import (
	"net/url"
)

type Image struct {
	Size, Sort            string
	PartId, Height, Width int
	Path                  url.URL
}
