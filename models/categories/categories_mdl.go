package categories

import (
	"net/url"
	"time"
)

type Category struct {
	CategoryId, ParentId, Sort   int
	DateAdded                    time.Time
	Title, ShortDesc, LongDesc   string
	ColorCode, FontCode          string
	Image                        url.URL
	IsLifestyle, VehicleSpecific bool
}
