package search

import (
	"errors"
	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/ninnemana/elastigo/lib"
	"os"
	"strconv"
)

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

	searchCurt := false
	searchAries := false
	searchAll := false

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

	if searchAries && searchCurt {
		searchAll = true
	}

	from := strconv.Itoa(page * count)
	size := strconv.Itoa(count)

	if searchAll {
		return elastigo.Search("mongo_all").Query(
			elastigo.Query().Search(query),
		).From(from).Size(size).Result(con)
	} else if searchCurt {
		return elastigo.Search("mongo_curt").Query(
			elastigo.Query().Search(query),
		).From(from).Size(size).Result(con)
	} else if searchAries {
		return elastigo.Search("mongo_aries").Query(
			elastigo.Query().Search(query),
		).From(from).Size(size).Result(con)
	}

	return nil, errors.New("no index for determined brands")
}
