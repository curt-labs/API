package forum

import (
	"database/sql"
	"errors"
	"time"

	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllForumGroups = `select * from ForumGroup`
	getForumGroup     = `select * from ForumGroup where forumGroupID = ?`
)

type Groups []Group
type Group struct {
	ID          int
	Name        string
	Description string
	Created     time.Time
	Topics      Topics
}

func GetAllGroups() (groups Groups, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllForumGroups)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	allTopics, err := GetAllTopics()
	allTopicsMap := allTopics.ToMap(MapToGroupID)
	if err != nil {
		return
	}

	for rows.Next() {
		var group Group
		if err = rows.Scan(&group.ID, &group.Name, &group.Description, &group.Created); err == nil {
			if topics, ok := allTopicsMap[group.ID]; ok {
				group.Topics = topics.(Topics)
			}
			groups = append(groups, group)
		}
	}

	return
}

func (g *Group) Get() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getForumGroup)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var group Group
	row := stmt.QueryRow(g.ID)
	err = row.Scan(&group.ID, &group.Name, &group.Description, &group.Created)

	if row == nil || err != nil {
		if row == nil {
			return errors.New("Invalid reference to Forum Group")
		}
		return err
	}

	g.ID = group.ID
	g.Name = group.Name
	g.Description = group.Description
	g.Created = group.Created

	if err = g.GetTopics(); err != nil {
		return err
	}

	return nil
}
