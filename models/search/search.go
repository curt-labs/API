package search

import (
	"errors"
	"os"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/ninnemana/elastigo/lib"
)

type SearchRequest struct {
	Query Query `json:"query"`
}

type Query struct {
	MultiMatch MultiMatch `json:"multi_match"`
}

type MultiMatch struct {
	Query  string   `json:"query"`
	Fields []string `json:"fields"`
}

func Dsl(query string, page int, count int, brand int, dtx *apicontext.DataContext) (*elastigo.SearchResult, error) {

	if page == 1 {
		page = 0
	}
	if count == 0 {
		count = 25
	}
	if query == "" {
		return nil, errors.New("cannot execute a search on an empty query")
	}

	var con *elastigo.Conn
	if host := os.Getenv("ELASTICSEARCH_IP"); host != "" {
		con = &elastigo.Conn{
			Protocol: elastigo.DefaultProtocol,
			Domain:   host,
			Port:     os.Getenv("ELASTIC_PORT"),
			Username: os.Getenv("ELASTIC_USER"),
			Password: os.Getenv("ELASTIC_PASS"),
		}
	}
	if con == nil {
		return nil, errors.New("failed to connect to elasticsearch")
	}

	qry := SearchRequest{
		Query: Query{
			MultiMatch: MultiMatch{
				Query:  query,
				Fields: []string{"part_number^1", "_all"},
			},
		},
	}

	searchCurt := false
	searchAries := false

	if brand == 0 {
		for _, br := range dtx.BrandArray {
			if br == 1 { // search curt
				searchCurt = true
			} else if br == 3 { // search aries
				searchAries = true
			}
		}
	} else {
		if brand == 1 {
			searchCurt = true
			searchAries = false
		} else if brand == 3 {
			searchAries = true
			searchCurt = false
		}
	}

	var index string
	if searchAries && searchCurt {
		index = "mongo_all"
	} else if searchAries {
		index = "mongo_aries"
	} else if searchCurt {
		index = "mongo_curt"
	}

	args := map[string]interface{}{
		"from": page * count,
		"size": count,
	}

	out, err := con.Search(index, "", args, qry)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
