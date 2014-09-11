package forum

import (
	"database/sql"
	"errors"
	"time"

	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllForumPosts    = `select * from ForumPost`
	getForumPost        = `select * from ForumPost where postID = ?`
	getForumThreadPosts = `select * from ForumPost where threadID = ?`
)

type Posts []Post
type Post struct {
	ID        int
	ParentID  int
	ThreadID  int
	Created   time.Time
	Title     string
	Post      string
	Name      string
	Email     string
	Company   string
	Notify    bool
	Approved  bool
	Active    bool
	IPAddress string
	Flag      bool
	Sticky    bool
}

func GetAllPosts() (posts Posts, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllForumPosts)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	for rows.Next() {
		var post Post
		if err = rows.Scan(&post.ID, &post.ParentID, &post.ThreadID, &post.Created, &post.Title, &post.Post, &post.Name, &post.Email, &post.Company, &post.Notify, &post.Approved, &post.Active, &post.IPAddress, &post.Flag, &post.Sticky); err == nil {
			posts = append(posts, post)
		}
	}

	return
}

func (p *Post) Get() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getForumPost)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var post Post
	row := stmt.QueryRow(p.ID)
	err = row.Scan(&post.ID, &post.ParentID, &post.ThreadID, &post.Created, &post.Title, &post.Post, &post.Name, &post.Email, &post.Company, &post.Notify, &post.Approved, &post.Active, &post.IPAddress, &post.Flag, &post.Sticky)

	if row == nil || err != nil {
		if row == nil {
			return errors.New("Invalid reference to Forum Post")
		}
		return err
	}

	p.ID = post.ID
	p.ParentID = post.ParentID
	p.ThreadID = post.ThreadID
	p.Created = post.Created
	p.Title = post.Title
	p.Post = post.Post
	p.Name = post.Name
	p.Email = post.Email
	p.Company = post.Company
	p.Notify = post.Notify
	p.Approved = post.Approved
	p.Active = post.Active
	p.IPAddress = post.IPAddress
	p.Flag = post.Flag
	p.Sticky = post.Sticky

	return nil
}

func (t *Thread) GetPosts() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getForumThreadPosts)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(t.ID)
	for rows.Next() {
		var post Post
		if err = rows.Scan(&post.ID, &post.ParentID, &post.ThreadID, &post.Created, &post.Title, &post.Post, &post.Name, &post.Email, &post.Company, &post.Notify, &post.Approved, &post.Active, &post.IPAddress, &post.Flag, &post.Sticky); err == nil {
			t.Posts = append(t.Posts, post)
		}
	}

	return nil
}
