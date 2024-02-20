package main

import (
	"fmt"
	"net/http"
	"strconv"
	"test-sns/models"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// var Db *sql.DB

// func init() {
// 	var err error
// 	Db, err = sql.Open("postgres", "user=db-test01 dbname=db-sns-test password=Itsuki0530 sslmode=disable")
// 	if err != nil {
// 		panic(err)
// 	}
// }

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

func getpost(c *gin.Context) {
	id_s := c.Query("id") // URLパラメータからidを取得する
	id, _ := strconv.Atoi(id_s)

	fmt.Print("getpost")
	post, err := models.GetPost(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": post})
}

// func createpost(c *gin.Context) {

// 	// content := c.Query("content")
// 	// author := c.Query("author")
// 	post := models.Post{Content: "content", Author: "author"}
// 	err := models.DB.Create(&post).Error
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// }

func createcomment(c *gin.Context) {
	content := c.Query("content")
	author := c.Query("author")
	post_id_s := c.Query("post_id")
	post_id, _ := strconv.Atoi(post_id_s)
	post, _ := models.GetPost(post_id)
	comment := models.Comment{Content: content, Author: author, Post_id: int(post.ID)}
	// comment.Create()
	err := models.DB.Create(&comment).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func getallpost(c *gin.Context) {
	posts, _ := models.GetAllPost()

	c.JSON(http.StatusOK, gin.H{"post": posts})
}

func main() {
	router := gin.Default()
	models.ConnectDataBase()
	router.GET("/getpost", getpost)
	// router.GET("/createpost", createpost)
	router.GET("/createcomment", createcomment)
	router.GET("/getallpost", getallpost)

	router.Run("localhost:8080")
}
