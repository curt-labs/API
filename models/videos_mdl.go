package models

import (
	"net/url"
)

type Video struct {
	YouTubeVideoId, Type string
	IsPrimary            bool
	TypeId               int
	TypeIcon             url.URL
}
