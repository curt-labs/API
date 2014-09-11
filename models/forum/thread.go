package forum

import (
	"database/sql"
	"errors"
	"time"

	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllForumThreads   = `select * from ForumThread`
	getForumThread       = `select * from ForumThread where threadID = ?`
	getForumTopicthreads = `select * from ForumThread where topicID = ?`
)

type Threads []Thread
type Thread struct {
	ID      int
	TopicID int
	Created time.Time
	Active  bool
	Closed  bool
	Posts   Posts
}

func GetAllThreads() (threads Threads, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllForumThreads)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	allPosts, err := GetAllPosts()
	allPostsMap := allPosts.ToMap(MapToThreadID)
	if err != nil {
		return
	}

	for rows.Next() {
		var thread Thread
		if err = rows.Scan(&thread.ID, &thread.TopicID, &thread.Created, &thread.Active, &thread.Closed); err == nil {
			if posts, ok := allPostsMap[thread.ID]; ok {
				thread.Posts = posts.(Posts)
			}
			threads = append(threads, thread)
		}
	}

	return
}

func (t *Thread) Get() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getForumThread)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var thread Thread
	row := stmt.QueryRow(t.ID)
	err = row.Scan(&thread.ID, &thread.TopicID, &thread.Created, &thread.Active, &thread.Closed)

	if row == nil || err != nil {
		if row == nil {
			return errors.New("Invalid reference to Forum Thread")
		}
		return err
	}

	t.ID = thread.ID
	t.TopicID = thread.TopicID
	t.Created = thread.Created
	t.Active = thread.Active
	t.Closed = thread.Closed

	if err = t.GetPosts(); err != nil {
		return err
	}

	return nil
}

func (t *Topic) GetThreads() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getForumTopicthreads)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(t.ID)
	for rows.Next() {
		var thread Thread
		if err = rows.Scan(&thread.ID, &thread.TopicID, &thread.Created, &thread.Active, &thread.Closed); err == nil {
			t.Threads = append(t.Threads, thread)
		}
	}

	return nil
}
