package forum

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllForumGroups = `select * from ForumGroup`
	getForumGroup     = `select * from ForumGroup where forumGroupID = ?`
	addForumGroup     = `insert ForumGroup(name,description,createdDate) values(?,?,UTC_TIMESTAMP())`
	updateForumGroup  = `update ForumGroup set name = ?, description = ? where forumGroupID = ?`
	deleteForumGroup  = `delete from ForumGroup where forumGroupID = ?`
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

func (g *Group) Add() error {
	if len(strings.TrimSpace(g.Name)) == 0 {
		return errors.New("Group must have a name")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(addForumGroup)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(g.Name, g.Description)
	if err != nil {
		return err
	}

	if id, err := res.LastInsertId(); err != nil {
		return err
	} else {
		g.ID = int(id)
	}

	return nil
}

func (g *Group) Update() error {
	if g.ID == 0 {
		return errors.New("Invalid Group ID")
	}

	if len(strings.TrimSpace(g.Name)) == 0 {
		return errors.New("Group must have a name")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(updateForumGroup)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(g.Name, g.Description, g.ID); err != nil {
		return err
	}

	return nil
}

func (g *Group) Delete() error {
	if err := g.DeleteTopics(); err != nil {
		return err
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteForumGroup)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(g.ID); err != nil {
		return err
	}

	return nil
}
