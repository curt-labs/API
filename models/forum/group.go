package forum

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllForumGroups = `select FG.forumGroupID, FG.name, FG.description, FG.createdDate, FG.websiteID
						 from ForumGroup FG
						 join WebsiteToBrand WTB on WTB.WebsiteID = FG.websiteID
						 join ApiKeyToBrand AKB on AKB.brandID = WTB.brandID
						 join ApiKey AK on AK.id = AKB.keyID
						 where AK.api_key = ? && (FG.websiteID = ? || 0 = ?)`
	getForumGroup = `select FG.forumGroupID, FG.name, FG.description, FG.createdDate, FG.websiteID
	                 from ForumGroup FG
	                 join WebsiteToBrand WTB on WTB.WebsiteID = FG.websiteID
	                 join ApiKeyToBrand AKB on AKB.brandID = WTB.brandID
	                 join ApiKey AK on AK.id = AKB.keyID
	                 where AK.api_key = ? && (FG.websiteID = ? || 0 = ?) && FG.forumGroupID = ?`
	addForumGroup    = `insert ForumGroup(createdDate,name,description,websiteID) values(UTC_TIMESTAMP(), ?, ?, ?)`
	updateForumGroup = `update ForumGroup set name = ?, description = ?, websiteID = ? where forumGroupID = ?`
	deleteForumGroup = `delete from ForumGroup where forumGroupID = ?`
)

type Groups []Group
type Group struct {
	ID          int
	WebsiteID   int
	Name        string
	Description string
	Created     time.Time
	Topics      Topics
}

func GetAllGroups(dtx *apicontext.DataContext) (groups Groups, err error) {
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

	rows, err := stmt.Query(dtx.APIKey, dtx.WebsiteID, dtx.WebsiteID)
	if err != nil {
		return
	}

	allTopics, err := GetAllTopics(dtx)
	allTopicsMap := allTopics.ToMap(MapToGroupID)
	if err != nil {
		return
	}

	for rows.Next() {
		var group Group
		if err = rows.Scan(&group.ID, &group.Name, &group.Description, &group.Created, &group.WebsiteID); err == nil {
			if topics, ok := allTopicsMap[group.ID]; ok {
				group.Topics = topics.(Topics)
			}
			groups = append(groups, group)
		}
	}
	defer rows.Close()

	return
}

func (g *Group) Get(dtx *apicontext.DataContext) error {
	if g.ID == 0 {
		return errors.New("Invalid Group ID")
	}

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

	row := stmt.QueryRow(dtx.APIKey, dtx.WebsiteID, dtx.WebsiteID, g.ID)
	if err = row.Scan(&group.ID, &group.Name, &group.Description, &group.Created, &group.WebsiteID); err != nil {
		return err
	}

	g.ID = group.ID
	g.WebsiteID = group.WebsiteID
	g.Name = group.Name
	g.Description = group.Description
	g.Created = group.Created

	if err = g.GetTopics(dtx); err != nil {
		return err
	}

	return nil
}

func (g *Group) Add(dtx *apicontext.DataContext) error {
	if len(strings.TrimSpace(g.Name)) == 0 {
		return errors.New("Group must have a name")
	}

	if g.WebsiteID != dtx.WebsiteID {
		g.WebsiteID = dtx.WebsiteID
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

	res, err := stmt.Exec(g.Name, g.Description, g.WebsiteID)
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

func (g *Group) Update(dtx *apicontext.DataContext) error {
	if g.ID == 0 {
		return errors.New("Invalid Group ID")
	}

	if len(strings.TrimSpace(g.Name)) == 0 {
		return errors.New("Group must have a name")
	}

	if g.WebsiteID != dtx.WebsiteID {
		g.WebsiteID = dtx.WebsiteID
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

	if _, err := stmt.Exec(g.Name, g.Description, g.WebsiteID, g.ID); err != nil {
		return err
	}

	return nil
}

func (g *Group) Delete(dtx *apicontext.DataContext) error {
	if g.ID == 0 {
		return errors.New("Invalid Group ID")
	}

	if err := g.DeleteTopics(dtx); err != nil {
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
