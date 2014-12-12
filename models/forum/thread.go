package forum

import (
	"database/sql"
	"errors"
	"time"

	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllForumThreads = `select FTH.threadID, FTH.topicID, FTH.createdDate, FTH.active, FTH.closed
						  from ForumThread FTH
						  join ForumTopic FT on FT.topicID = FTH.topicID
						  join ForumGroup FG on FG.forumGroupID = FT.TopicGroupID
						  join WebsiteToBrand WTB on WTB.WebsiteID = FG.WebsiteID
						  join ApiKeyToBrand AKB on AKB.brandID = WTB.brandID
						  join ApiKey AK on AK.id = AKB.keyID
						  where AK.api_key = ? && (FG.websiteID = ? || 0 = ?)`
	getForumThread = `select FTH.threadID, FTH.topicID, FTH.createdDate, FTH.active, FTH.closed
					  from ForumThread FTH
					  join ForumTopic FT on FT.topicID = FTH.topicID
					  join ForumGroup FG on FG.forumGroupID = FT.TopicGroupID
					  join WebsiteToBrand WTB on WTB.WebsiteID = FG.WebsiteID
					  join ApiKeyToBrand AKB on AKB.brandID = WTB.brandID
					  join ApiKey AK on AK.id = AKB.keyID
					  where AK.api_key = ? && (FG.websiteID = ? || 0 = ?) && FTH.threadID = ?`
	getForumTopicThreads = `select FTH.threadID, FTH.topicID, FTH.createdDate, FTH.active, FTH.closed
							from ForumThread FTH
							join ForumTopic FT on FT.topicID = FTH.topicID
							join ForumGroup FG on FG.forumGroupID = FT.TopicGroupID
							join WebsiteToBrand WTB on WTB.WebsiteID = FG.WebsiteID
							join ApiKeyToBrand AKB on AKB.brandID = WTB.brandID
							join ApiKey AK on AK.id = AKB.keyID
							where AK.api_key = ? && (FG.websiteID = ? || 0 = ?) && FTH.topicID = ?`
	addForumThread    = `insert into ForumThread(topicID,createdDate,active,closed) values(?,UTC_TIMESTAMP(), ?, ?)`
	updateForumThread = `update ForumThread set topicID = ?, active = ?, closed = ? where threadID = ?`
	deleteForumThread = `delete from ForumThread where threadID = ?`
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

func GetAllThreads(dtx *apicontext.DataContext) (threads Threads, err error) {
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

	rows, err := stmt.Query(dtx.APIKey, dtx.WebsiteID, dtx.WebsiteID)
	if err != nil {
		return
	}

	allPosts, err := GetAllPosts(dtx)
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
	defer rows.Close()

	return
}

func (t *Thread) Get(dtx *apicontext.DataContext) error {
	if t.ID == 0 {
		return errors.New("Invalid Thread ID")
	}

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
	row := stmt.QueryRow(dtx.APIKey, dtx.WebsiteID, dtx.WebsiteID, t.ID)
	if err = row.Scan(&thread.ID, &thread.TopicID, &thread.Created, &thread.Active, &thread.Closed); err != nil {
		return err
	}

	t.ID = thread.ID
	t.TopicID = thread.TopicID
	t.Created = thread.Created
	t.Active = thread.Active
	t.Closed = thread.Closed

	if err = t.GetPosts(dtx); err != nil {
		return err
	}

	return nil
}

func (t *Topic) GetThreads(dtx *apicontext.DataContext) error {
	if t.ID == 0 {
		return errors.New("Invalid Topic ID")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getForumTopicThreads)
	if err != nil {
		return err
	}
	defer stmt.Close()

	allPosts, err := GetAllPosts(dtx)
	allPostsMap := allPosts.ToMap(MapToThreadID)
	if err != nil {
		return err
	}

	rows, err := stmt.Query(dtx.APIKey, dtx.WebsiteID, dtx.WebsiteID, t.ID)
	if err != nil {
		return err
	}

	for rows.Next() {
		var thread Thread
		if err = rows.Scan(&thread.ID, &thread.TopicID, &thread.Created, &thread.Active, &thread.Closed); err == nil {
			if posts, ok := allPostsMap[thread.ID]; ok {
				thread.Posts = posts.(Posts)
			}
			t.Threads = append(t.Threads, thread)
		}
	}
	defer rows.Close()

	return nil
}

func (t *Thread) Add() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(addForumThread)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(t.TopicID, t.Active, t.Closed)
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

func (t *Thread) Update() error {
	if t.ID == 0 {
		return errors.New("Invalid Thread ID")
	}

	if t.TopicID == 0 {
		return errors.New("Invalid Topic ID")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(updateForumThread)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(t.TopicID, t.Active, t.Closed, t.ID); err != nil {
		return err
	}

	return nil
}

func (t *Thread) Delete() error {
	if t.ID == 0 {
		return errors.New("Invalid Thread ID")
	}

	if err := t.DeletePosts(); err != nil {
		return err
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteForumThread)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(t.ID); err != nil {
		return err
	}

	return nil
}

func (t *Topic) DeleteThreads(dtx *apicontext.DataContext) error {
	var err error
	if len(t.Threads) == 0 {
		//try getting
		if err = t.Get(dtx); err != nil {
			return err
		}
	}
	for _, thread := range t.Threads {
		if err = thread.Delete(); err != nil {
			return err
		}
	}

	return nil
}
