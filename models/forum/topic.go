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
	getAllForumTopics   = `select * from ForumTopic`
	getForumTopic       = `select * from ForumTopic where topicID = ?`
	getForumGroupTopics = `select ft.* from ForumTopic ft
                              inner join ForumGroup fg on ft.TopicGroupID = fg.forumGroupID
                              where ft.TopicGroupID = ?`
	addForumTopic    = `insert into ForumTopic(TopicGroupID, name, description, image, createdDate, active, closed) values (?,?,?,?,UTC_TIMESTAMP(),?,?)`
	updateForumTopic = `update ForumTopic set TopicGroupID = ?, name = ?, description = ?, image = ?, active = ?, closed = ? where topicID = ?`
	deleteForumTopic = `delete from ForumTopic where topicID = ?`
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
	defer rows.Close()

	return
}

func (t *Topic) Get() error {
	if t.ID == 0 {
		return errors.New("Invalid Topic ID")
	}
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
	if err = row.Scan(&topic.ID, &topic.GroupID, &topic.Name, &topic.Description, &topic.Image, &topic.Created, &topic.Active, &topic.Closed); err != nil {
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
	if g.ID == 0 {
		return errors.New("Invalid Group ID")
	}

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

	allThreads, err := GetAllThreads()
	allThreadsMap := allThreads.ToMap(MapToTopicID)
	if err != nil {
		return err
	}

	rows, err := stmt.Query(g.ID)
	if err != nil {
		return err
	}

	for rows.Next() {
		var topic Topic
		if err = rows.Scan(&topic.ID, &topic.GroupID, &topic.Name, &topic.Description, &topic.Image, &topic.Created, &topic.Active, &topic.Closed); err == nil {
			if threads, ok := allThreadsMap[topic.ID]; ok {
				topic.Threads = threads.(Threads)
			}
			g.Topics = append(g.Topics, topic)
		}
	}
	defer rows.Close()

	return nil
}

func (t *Topic) Add() error {
	if t.GroupID == 0 {
		return errors.New("Topic must have a Group ID")
	}
	if len(strings.TrimSpace(t.Name)) == 0 {
		return errors.New("Topic must have a name")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(addForumTopic)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(t.GroupID, t.Name, t.Description, t.Image, t.Active, t.Closed)
	if err != nil {
		return err
	}

	if id, err := res.LastInsertId(); err != nil {
		return err
	} else {
		t.ID = int(id)
	}

	return nil
}

func (t *Topic) Update() error {
	if t.ID == 0 {
		return errors.New("Invalid Topic ID")
	}

	if t.GroupID == 0 {
		return errors.New("Topic must have a Group ID")
	}
	if len(strings.TrimSpace(t.Name)) == 0 {
		return errors.New("Topic must have a name")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(updateForumTopic)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(t.GroupID, t.Name, t.Description, t.Image, t.Active, t.Closed, t.ID); err != nil {
		return err
	}

	return nil
}

func (t *Topic) Delete() error {
	if t.ID == 0 {
		return errors.New("Invalid Topic ID")
	}

	if err := t.DeleteThreads(); err != nil {
		return err
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteForumTopic)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(t.ID); err != nil {
		return err
	}

	return nil
}

func (g *Group) DeleteTopics() error {
	var err error
	if len(g.Topics) == 0 {
		//try getting
		if err = g.Get(); err != nil {
			return err
		}
	}
	for _, topic := range g.Topics {
		if err = topic.Delete(); err != nil {
			return err
		}
	}
	return nil
}
