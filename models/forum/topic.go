package forum

import (
	"database/sql"
	"errors"
	"time"

	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllForumTopics   = `select * from ForumTopic`
	getForumTopic       = `select * from ForumTopic where topicID = ?`
	getForumGroupTopics = `select ft.* from ForumTopic ft
						   inner join ForumGroup fg on ft.TopicGroupID = fg.forumGroupID
						   where ft.TopicGroupID = ?`
)

type Topics []Topic
type Topic struct {
	ID          int
	GroupID     int
	Name        string
	Description string
	Image       string
	Created     time.Time
	Active      bool
	Closed      bool
	Threads     Threads
}

func GetAllTopics() (topics Topics, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllForumTopics)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	allThreads, err := GetAllThreads()
	allThreadsMap := allThreads.ToMap(MapToTopicID)
	if err != nil {
		return
	}

	for rows.Next() {
		var topic Topic
		if err = rows.Scan(&topic.ID, &topic.GroupID, &topic.Name, &topic.Description, &topic.Image, &topic.Created, &topic.Active, &topic.Closed); err == nil {
			if threads, ok := allThreadsMap[topic.ID]; ok {
				topic.Threads = threads.(Threads)
			}
			topics = append(topics, topic)
		}
	}

	return
}

func (t *Topic) Get() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getForumTopic)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var topic Topic
	row := stmt.QueryRow(t.ID)
	err = row.Scan(&topic.ID, &topic.GroupID, &topic.Name, &topic.Description, &topic.Image, &topic.Created, &topic.Active, &topic.Closed)
	if row == nil || err != nil {
		if row == nil {
			return errors.New("Invalid reference to Forum Topic")
		}
		return err
	}

	t.ID = topic.ID
	t.GroupID = topic.GroupID
	t.Name = topic.Name
	t.Description = topic.Description
	t.Image = topic.Description
	t.Created = topic.Created
	t.Active = topic.Active
	t.Closed = topic.Closed

	if err = t.GetThreads(); err != nil {
		return err
	}

	return nil
}

func (g *Group) GetTopics() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getForumGroupTopics)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(g.ID)
	for rows.Next() {
		var topic Topic
		if err = rows.Scan(&topic.ID, &topic.GroupID, &topic.Name, &topic.Description, &topic.Image, &topic.Created, &topic.Active, &topic.Closed); err == nil {
			g.Topics = append(g.Topics, topic)
		}
	}

	return nil
}
