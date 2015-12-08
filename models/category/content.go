package category

// Content ...
type Content struct {
	Id          int         `bson:"contentid" json:"id" xml:"id"`
	Text        string      `bson:"text" json:"text" xml:"text"`
	ContentType ContentType `json:"contentType" xml:"contentType"`
}

// ContentType ...
type ContentType struct {
	Id         int    `bson:"contenttypeid" json:"id" xml:"id"`
	Type       string `bson:"type" json:"type" xml:"type"`
	AllowsHTML bool   `bson:"allows_html" json:"allows_html" xml:"allows_html"`
	IsPrivate  bool   `bson:"isprivate" json:"isprivate" xml:"isprivate"`
}
