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
	getAllForumPosts       = `select * from ForumPost`
	getForumPost           = `select * from ForumPost where postID = ?`
	getForumThreadPosts    = `select * from ForumPost where threadID = ?`
	addForumPost           = `insert into ForumPost(parentID,threadID,createdDate,title,post,name,email,company,notify,approved,active,IPAddress,flag,sticky) values(?,?,UTC_TIMESTAMP(),?,?,?,?,?,?,?,1,?,?,?)`
	updateForumPost        = `update ForumPost set parentID = ?, threadID = ?, title = ?, post = ?, name = ?, email = ?, company = ?, notify = ?, approved = ?, IPAddress = ?, flag = ?, sticky = ? where postID = ?`
	deleteForumPost        = `delete from ForumPost where postID = ?`
	deleteForumThreadPosts = `delete from ForumPost where threadID = ?`
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
	defer rows.Close()

	return
}

func (p *Post) Get() error {
	if p.ID == 0 {
		return errors.New("Invalid Post ID")
	}

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
	if err = row.Scan(&post.ID, &post.ParentID, &post.ThreadID, &post.Created, &post.Title, &post.Post, &post.Name, &post.Email, &post.Company, &post.Notify, &post.Approved, &post.Active, &post.IPAddress, &post.Flag, &post.Sticky); err != nil {
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
	if t.ID == 0 {
		return errors.New("Invalid Thread ID")
	}

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

func (p *Post) Add() error {
	if p.ThreadID == 0 {
		return errors.New("Invalid Thread ID")
	}

	if len(strings.TrimSpace(p.Title)) == 0 {
		return errors.New("Post must have a title")
	}

	if len(strings.TrimSpace(p.Post)) == 0 {
		return errors.New("Post must have a message")
	}

	if p.Notify && len(strings.TrimSpace(p.Email)) == 0 {
		return errors.New("Must have an email address in order to be notified")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(addForumPost)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(p.ParentID, p.ThreadID, p.Title, p.Post, p.Name, p.Email, p.Company, p.Notify, p.Approved, p.IPAddress, p.Flag, p.Sticky)
	if err != nil {
		return err
	}

	if id, err := res.LastInsertId(); err != nil {
		return err
	} else {
		p.ID = int(id)
	}

	return nil
}

func (p *Post) Update() error {
	if p.ID == 0 {
		return errors.New("Invalid Post ID")
	}

	if p.ThreadID == 0 {
		return errors.New("Invalid Thread ID")
	}

	if len(strings.TrimSpace(p.Title)) == 0 {
		return errors.New("Post must have a title")
	}

	if len(strings.TrimSpace(p.Post)) == 0 {
		return errors.New("Post must have a message")
	}

	if p.Notify && len(strings.TrimSpace(p.Email)) == 0 {
		return errors.New("Must have an email address in order to be notified")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(updateForumPost)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(p.ParentID, p.ThreadID, p.Title, p.Post, p.Name, p.Email, p.Company, p.Notify, p.Approved, p.IPAddress, p.Flag, p.Sticky, p.ID); err != nil {
		return err
	}

	return nil
}

func (p *Post) Delete() error {
	if p.ID == 0 {
		return errors.New("Invalid Post ID")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteForumPost)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(p.ID); err != nil {
		return err
	}

	return nil
}

func (t *Thread) DeletePosts() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteForumThreadPosts)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(t.ID); err != nil {
		return err
	}

	return nil
}
