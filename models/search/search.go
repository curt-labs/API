package search

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	elastic "gopkg.in/olivere/elastic.v2"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/mattbaird/elastigo/lib"
)

func newConn() (*elastic.Client, error) {
	hosts := []string{"http://127.0.0.1:9200"}

	if d := os.Getenv("ELASTICSEARCH_IP"); d != "" {
		hosts = []string{}
		urls := strings.Split(d, ",")
		for _, u := range urls {
			hosts = append(
				hosts,
				fmt.Sprintf("http://%s:9200", u),
			)
		}
	}

	user := os.Getenv("ELASTIC_USER")
	pass := os.Getenv("ELASTIC_PASS")

	funcs := []elastic.ClientOptionFunc{
		elastic.SetURL(hosts...),
		elastic.SetMaxRetries(10),
	}

	if user != "" && pass != "" {
		funcs = append(funcs, elastic.SetBasicAuth(user, pass))
	}

	return elastic.NewSimpleClient(funcs...)
}

func Dsl(query string, page int, count int, brand int, dtx *apicontext.DataContext, rawPartNumber string) (*elastic.SearchResult, error) {

	if page == 1 {
		page = 0
	}
	if count == 0 {
		count = 25
	}
	if query == "" {
		return nil, errors.New("cannot execute a search on an empty query")
	}

	c, err := newConn()
	if err != nil {
		return nil, err
	}

	return c.Search(findIndex(brand, dtx)).From(page * count).Size(count).Query(elastic.NewQueryStringQuery(query)).Do()
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
	searchRetrac := false
	searchUWS := false

	if brand == 0 {
		for _, br := range dtx.BrandArray {
			if br == 1 { // search curt
				searchCurt = true
			} else if br == 3 { // search aries
				searchAries = true
			} else if br == 4 { // search luverne
				searchLuverne = true
			} else if br == 5 { // search retrac
				searchRetrac = true
			} else if br == 6 { // search UWS
				searchUWS = true
			}
		}
	} else {
		if brand == 1 {
			searchCurt = true
		} else if brand == 3 {
			searchAries = true
		} else if brand == 4 {
			searchLuverne = true
		} else if brand == 5 {
			searchRetrac = true
		} else if brand == 6 {
			searchUWS = true
		}
	}

	index := "all"

	if searchCurt && !searchAries && !searchLuverne && !searchRetrac && !searchUWS {
		index = "curt"
	}
	if searchAries && !searchCurt && !searchLuverne && !searchRetrac && !searchUWS {
		index = "aries"
	}
	if searchLuverne && !searchAries && !searchCurt && !searchRetrac && !searchUWS {
		index = "luverne"
	}
	if searchRetrac && !searchAries && !searchCurt && !searchLuverne && !searchUWS {
		index = "retrac"
	}
	if searchUWS && !searchAries && !searchCurt && !searchLuverne && !searchRetrac {
		index = "uws"
	}

	return index
}
