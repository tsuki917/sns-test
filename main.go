package main

import (
	"fmt"
	"net/http"
	"sns-test/controllers"
	"sns-test/middlewares"
	"sns-test/models"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

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

// func createcomment(c *gin.Context) {
// 	author := c.Query("author")
// 	post_id_s := c.Query("post_id")
// 	post_id, _ := strconv.Atoi(post_id_s)
// 	post, _ := models.GetPost(post_id)
// 	comment := models.Comment{Content: content, Author: author, Post_id: int(post.ID)}
// 	// comment.Create()
// 	err := models.DB.Create(&comment).Error
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// }

func getallpost(c *gin.Context) {
	posts, err := models.GetAllPost()
	client_userid, _ := strconv.ParseUint(c.Query("userid"), 10, 64)

	fmt.Println("client_userid")
	fmt.Println(client_userid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	type Thread struct {
		Post   models.Post
		User   models.User_Post
		IsFavo bool
	}
	threads := []Thread{}
	for _, post := range posts {
		thread := Thread{}
		thread.Post = post
		u := models.User{}
		err := models.DB.Where("id=?", post.UserId).First(&u).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = models.DB.Model(&models.Favo{}).Where("post_id=?", post.ID).Count(&thread.Post.FavoNum).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		thread.User.Username = u.Username
		thread.User.UserTag = u.UserTag
		thread.User.ImgPath = u.ImgPath
		thread.User.ID = u.ID
		thread.IsFavo = models.IsFavo(uint(client_userid), post.ID)
		threads = append(threads, thread)
	}

	c.JSON(http.StatusOK, gin.H{"threads": threads})
}

func main() {
	router := gin.Default()

	models.ConnectDataBase()
	router.Use(cors.New(cors.Config{
		// アクセスを許可したいアクセス元
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		// アクセスを許可したいHTTPメソッド(以下の例だとPUTやDELETEはアクセスできません)
		AllowMethods: []string{
			"POST",
			"GET",
		},
		// 許可したいHTTPリクエストヘッダ
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Access-Control-Allow-Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
		// cookieなどの情報を必要とするかどうか
		AllowCredentials: false,
		// preflightリクエストの結果をキャッシュする時間

	}))
	public := router.Group("/api")
	// models.AddFavo(2, 1)
	protected := router.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)
	router.GET("/getpost", getpost)
	protected.GET("/user", controllers.CurrentUser)
	router.POST("/createpost", models.Createpost)
	// router.GET("/createcomment", createcomment)
	router.GET("/getallpost", getallpost)
	router.POST("/addfavo", models.AddFavo)
	router.POST("/deletefavo", models.DeleteFavo)
	router.Run("localhost:8080")
}
