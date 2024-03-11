package main

import (
	"fmt"
	"net/http"
	"strconv"
	"test-sns/controllers"
	"test-sns/middlewares"
	"test-sns/models"

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

func createpost(c *gin.Context) {
	content := c.Query("content")
	post := models.Post{Content: content}
	err := models.DB.Create(&post).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// func createcomment(c *gin.Context) {
// 	content := c.Query("content")
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
	posts, _ := models.GetAllPost()

	c.JSON(http.StatusOK, gin.H{"posts": posts})
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

	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)
	router.GET("/getpost", getpost)
	protected := router.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", controllers.CurrentUser)
	router.GET("/createpost", createpost)
	// router.GET("/createcomment", createcomment)
	router.GET("/getallpost", getallpost)

	router.Run("localhost:8080")

}
