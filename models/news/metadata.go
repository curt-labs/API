package news_model

import (
	"database/sql"
	"github.com/curt-labs/API/helpers/database"
)

var (
	GetMetadataStmt = `select mt.id, mt.type, mt.value from MetaTag as mt
						join NewsToMetaTag as nm on mt.id = nm.meta_id
						where nm.news_id = ?
						order by mt.type, mt.value`
)

type Metadata struct {
	ID    int    `json:"-" xml:"-"`
	Type  string `json:"type" xml:"type"`
	Value string `json:"value" xml:"value"`
}

func GetMetadata(newsID int) ([]Metadata, error) {
	var data []Metadata
	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return data, err
	}
	defer db.Close()

	stmt, err := db.Prepare(GetMetadataStmt)
	if err != nil {
		return data, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(newsID)
	if err != nil {
		return data, err
	}

	var t, v *string
	var id *int
	for rows.Next() {
		err = rows.Scan(&id, &t, &v)
		if err != nil || id == nil || t == nil || v == nil {
			continue
		}

		mt := Metadata{
			ID:    *id,
			Type:  *t,
			Value: *v,
		}

		data = append(data, mt)
	}

	return data, nil
}
