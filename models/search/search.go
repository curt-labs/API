package search

import (
	"errors"
	"os"
	"strconv"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/ninnemana/elastigo/lib"
)

func Dsl(query string, page int, count int, brand int, dtx *apicontext.DataContext, rawPartNumber string) (*elastigo.SearchResult, error) {

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

	from := strconv.Itoa(page * count)
	size := strconv.Itoa(count)

	filter := elastigo.Filter()
	if rawPartNumber != "" {
		filter.Terms("raw_part", rawPartNumber)
	}
	index := findIndex(brand, dtx)

	res, err := elastigo.Search(index).Query(
		elastigo.Query().Search(query),
	).Filter(filter).From(from).Size(size).Result(con)
	return res, err
}

func ExactAndCloseDsl(query string, page int, count int, brand int, dtx *apicontext.DataContext) (*elastigo.SearchResult, error) {

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

	from := strconv.Itoa(page * count)
	size := strconv.Itoa(count)

	index := findIndex(brand, dtx)

	args := map[string]interface{}{
		"from": from,
		"size": size,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"match": map[string]interface{}{
						"_all": query,
					},
				},
				"should": map[string]interface{}{
					"match": map[string]interface{}{
						"raw_part": map[string]interface{}{
							"query": query,
							"boost": 10,
						},
					},
				},
			},
		},
	}
	res, err := con.Search(index, "", nil, args)
	return &res, err
}

func findIndex(brand int, dtx *apicontext.DataContext) string {
	searchCurt := false
	searchAries := false
	searchLuverne := false

	if brand == 0 {
		for _, br := range dtx.BrandArray {
			if br == 1 { // search curt
				searchCurt = true
			} else if br == 3 { // search aries
				searchAries = true
			} else if br == 4 { // search luverne
				searchLuverne = true
			}
		}
	} else {
		if brand == 1 {
			searchCurt = true
		} else if brand == 3 {
			searchAries = true
		} else if brand == 4 {
			searchLuverne = true
		}
	}

	index := "all"
	if searchCurt && !searchAries && !searchLuverne {
		index = "curt"
	}
	if searchAries && !searchCurt && !searchLuverne {
		index = "aries"
	}
	if searchLuverne && !searchCurt && !searchAries {
		index = "luverne"
	}
	return index
}
