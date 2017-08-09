package forum

import (
	"errors"
	"strings"
	"time"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllForumTopics = `select FT.topicID, FT.TopicGroupID, FT.name, FT.description, FT.image, FT.createdDate, FT.active, FT.closed
						 from ForumTopic FT
						 join ForumGroup FG on FG.forumGroupID = FT.TopicGroupID
						 join WebsiteToBrand WTB on WTB.WebsiteID = FG.WebsiteID
						 join ApiKeyToBrand AKB on AKB.brandID = WTB.brandID
						 join ApiKey AK on AK.id = AKB.keyID
						 where AK.api_key = ? && (FG.websiteID = ? || 0 = ?)`
	getForumTopic = `select FT.topicID, FT.TopicGroupID, FT.name, FT.description, FT.image, FT.createdDate, FT.active, FT.closed
					 from ForumTopic FT
					 join ForumGroup FG on FG.forumGroupID = FT.TopicGroupID
					 join WebsiteToBrand WTB on WTB.WebsiteID = FG.WebsiteID
					 join ApiKeyToBrand AKB on AKB.brandID = WTB.brandID
					 join ApiKey AK on AK.id = AKB.keyID
					 where AK.api_key = ? && (FG.websiteID = ? || 0 = ?) && FT.topicID = ?`
	getForumGroupTopics = `select FT.topicID, FT.TopicGroupID, FT.name, FT.description, FT.image, FT.createdDate, FT.active, FT.closed
						   from ForumTopic FT
						   join ForumGroup FG on FG.forumGroupID = FT.TopicGroupID
						   where FT.TopicGroupID = ?`
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

func GetAllTopics(dtx *apicontext.DataContext) (topics Topics, err error) {
	err = database.Init()
	if err != nil {
		return
	}

	stmt, err := database.DB.Prepare(getAllForumTopics)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(dtx.APIKey, dtx.WebsiteID, dtx.WebsiteID)
	if err != nil {
		return
	}

	allThreads, err := GetAllThreads(dtx)
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

func (t *Topic) Get(dtx *apicontext.DataContext) error {
	if t.ID == 0 {
		return errors.New("Invalid Topic ID")
	}
	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(getForumTopic)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var topic Topic
	row := stmt.QueryRow(dtx.APIKey, dtx.WebsiteID, dtx.WebsiteID, t.ID)
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

	if err = t.GetThreads(dtx); err != nil {
		return err
	}

	return nil
}

func (g *Group) GetTopics(dtx *apicontext.DataContext) error {
	if g.ID == 0 {
		return errors.New("Invalid Group ID")
	}

	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(getForumGroupTopics)
	if err != nil {
		return err
	}
	defer stmt.Close()

	allThreads, err := GetAllThreads(dtx)
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

	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(addForumTopic)
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

	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(updateForumTopic)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(t.GroupID, t.Name, t.Description, t.Image, t.Active, t.Closed, t.ID); err != nil {
		return err
	}

	return nil
}

func (t *Topic) Delete(dtx *apicontext.DataContext) error {
	if t.ID == 0 {
		return errors.New("Invalid Topic ID")
	}

	if err := t.DeleteThreads(dtx); err != nil {
		return err
	}

	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(deleteForumTopic)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(t.ID); err != nil {
		return err
	}

	return nil
}

func (g *Group) DeleteTopics(dtx *apicontext.DataContext) error {
	var err error
	if len(g.Topics) == 0 {
		//try getting
		if err = g.Get(dtx); err != nil {
			return err
		}
	}
	for _, topic := range g.Topics {
		if err = topic.Delete(dtx); err != nil {
			return err
		}
	}
	return nil
}
