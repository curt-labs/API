package search

import (
	"errors"
	"github.com/mattbaird/elastigo/lib"
	"os"
)

func Dsl(query string, fields []string) (elastigo.SearchResult, error) {

	var con *elastigo.Conn
	if host := os.Getenv("ELASTICSEARCH_IP"); host != "" {
		con = &elastigo.Conn{
			Protocol: elastigo.DefaultProtocol,
			Domain:   host,
			Port:     os.Getenv("ELASTICSEARCH_PORT"),
			Username: os.Getenv("ELASTICSEARCH_USER"),
			Password: os.Getenv("ELASTICSEARCH_PASS"),
		}
	}
	if con == nil {
		return elastigo.SearchResult{}, errors.New("failed to connect to elasticsearch")
	}

	qry := map[string]interface{}{
		"query": map[string]interface{}{
			"query_string": map[string]interface{}{
				"query":  query,
				"fields": fields,
			},
		},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"*": map[string]interface{}{},
			},
		},
	}

	// Get the brands that this customer is subscribed to
	// and iterate over those brands hitting the search endpoint for all.
	// This will require us to join all brand results together into one
	// big SearchResult.

	var args map[string]interface{}
	return con.Search("*", "", args, qry)
}
