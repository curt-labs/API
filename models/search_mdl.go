package models

import (
	"github.com/curt-labs/GoAPI/helpers/database"
	"net/url"
)

type SearchResult struct {
	NextPage, Request SearchQuery
	Items             []SearchResultItem
}

type PartSearchResult struct {
	NextPage, Request SearchQuery
	Items             []*Part
}

type SearchQuery struct {
	Title                         string
	TotalResults                  int
	SearchTerms                   string
	Count, StartIndex             int
	InputEncoding, OutputEncoding string
}

type SearchResultItem struct {
	Kind                 string
	Title, HtmlTitle     string
	Link                 string
	Snippet, HtmlSnippet string
	Image                *url.URL
}

func (q *PartSearchResult) SearchParts(key string) error {

	q.Request.StartIndex = 1
	q.Request.Count = 10

	partChan := make(chan int)

	parts := make([]*Part, 0)

	go func() {
		qry, err := database.GetStatement("SearchPart")
		if !database.MysqlError(err) {
			rows, res, err := qry.Exec(q.Request.SearchTerms, q.Request.SearchTerms, q.Request.SearchTerms, q.Request.StartIndex, q.Request.Count)
			if !database.MysqlError(err) {
				pId := res.Map("partID")
				var lookup Lookup

				for _, row := range rows {
					lookup.Parts = append(lookup.Parts, &Part{
						PartId: row.Int(pId),
					})
				}
				lookup.Get(key)
				parts = append(parts, lookup.Parts...)
			}
		}
		partChan <- 1
	}()

	qry, err := database.GetStatement("SearchPartAttributes")
	if !database.MysqlError(err) {
		rows, res, err := qry.Exec(q.Request.SearchTerms, q.Request.SearchTerms, q.Request.StartIndex, q.Request.Count)
		if !database.MysqlError(err) {
			pId := res.Map("partID")
			var lookup Lookup

			for _, row := range rows {
				lookup.Parts = append(lookup.Parts, &Part{
					PartId: row.Int(pId),
				})
			}
			lookup.Get(key)
			parts = append(parts, lookup.Parts...)
		}
	}

	<-partChan

	q.Items = append(q.Items, parts...)
	return nil

}

func SearchCategories() (cats []Category, err error) {
	return
}
