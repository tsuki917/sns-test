package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Post struct {
	Id       int
	Content  string
	Author   string
	Comments []Comment
}
type Comment struct {
	Id      int
	Content string
	Author  string
	post_id int
	// Post    *Post
}

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "user=db-test01 dbname=db-sns-test password=Itsuki0530 sslmode=disable")
	if err != nil {
		panic(err)
	}
}

func (comment *Comment) Create() (err error) {
	if comment.post_id == 0 {
		err = errors.New("投稿が見つかりません")
		return
	}
	err = Db.QueryRow("insert into comments (content,author,post_id) values ($1,$2,$3) returning id", comment.Content, comment.Author, comment.post_id).Scan(&comment.Id)
	return
}

func GetPost(id int) (post Post, err error) {
	post = Post{}
	post.Comments = []Comment{}

	err = Db.QueryRow("select id,content,author from posts where id = $1", id).Scan(&post.Id, &post.Content, &post.Author)
	rows, err := Db.Query("select id,content, author from comments where post_id = $1", post.Id)
	if err != nil {
		fmt.Print(err)
		return
	}
	for rows.Next() {
		comment := Comment{post_id: post.Id}
		err = rows.Scan(&comment.Id, &comment.Content, &comment.Author)
		if err != nil {
			return
		}
		post.Comments = append(post.Comments, comment)
	}
	rows.Close()
	return
}

// func (post *Post) Create() (err error) {
// 	fmt.Println("Creating...")

//		statement := "insert into posts (content,author) values ($1,$2) returning id"
//		stmt, err := Db.Prepare(statement)
//		println("stmt")
//		println(stmt)
//		if err != nil {
//			return
//		}
//		defer stmt.Close()
//		err = stmt.QueryRow(post.Content, post.Author).Scan(&post.Id)
//		return
//	}
func (post *Post) Create() (err error) {
	err = Db.QueryRow("insert into posts (content,author) values ($1,$2) returning id", post.Content, post.Author).Scan(&post.Id)
	return
}
func (post *Post) Update() (err error) {
	_, err = Db.Exec("update posts set content = $2,author = $3 where id = $1", post.Id, post.Content, post.Author)
	return
}

func (post *Post) Delete() (err error) {
	_, err = Db.Exec("delete from posts where id = $1", post.Id)
	return
}

func getpost(c *gin.Context) {
	id_s := c.Query("id") // URLパラメータからidを取得する
	id, _ := strconv.Atoi(id_s)

	fmt.Print("getpost")
	post, err := GetPost(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": post})
}

func createpost(c *gin.Context) {

	content := c.Query("content")
	author := c.Query("author")
	post := Post{Content: content, Author: author}
	post.Create()
}

func createcomment(c *gin.Context) {
	content := c.Query("content")
	author := c.Query("author")
	post_id_s := c.Query("post_id")
	post_id, _ := strconv.Atoi(post_id_s)
	post, _ := GetPost(post_id)
	comment := Comment{Content: content, Author: author, post_id: post.Id}
	comment.Create()
}

func main() {
	router := gin.Default()
	router.GET("/getpost", getpost)
	router.GET("/createpost", createpost)
	router.GET("/createcomment", createcomment)
	router.Run("localhost:8080")
}
