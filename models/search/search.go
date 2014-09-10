package search

import (
	"errors"
	"github.com/mattbaird/elastigo/lib"
	"os"
)

func Dsl(query string, fields []string) (interface{}, error) {

	var con *elastigo.Conn
	if host := os.Getenv("ELASTICSEARCH_HOST"); host != "" {
		con = &elastigo.Conn{
			Protocol: elastigo.DefaultProtocol,
			Domain:   host,
			Port:     os.Getenv("ELASTICSEARCH_PORT"),
			Username: os.Getenv("ELASTICSEARCH_USERNAME"),
			Password: os.Getenv("ELASTICSEARCH_PASSWORD"),
		}
	}
	if con == nil {
		return nil, errors.New("failed to connect to elasticsearch")
	}

	qry := map[string]interface{}{
		"query": map[string]interface{}{
			"query_string": map[string]interface{}{
				"query":  query,
				"fields": fields,
			},
		},
	}

	var args map[string]interface{}
	res, e := con.Search("curt", "", args, qry)
	if e != nil {
		return nil, e
	}

	return res, nil
}
