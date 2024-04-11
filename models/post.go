package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model
	Content  string    `gorm:"size:255;" json:"content"`
	FavoNum  int64     `gorm:"size:255;" json:"favonum"`
	UserId   int       `gorm:"size:255;" json:"userid"`
	Comments []Comment `gorm:"size:255;" json:"comennts"`
}

func GetPost(id int) (post Post, err error) {
	post = Post{}
	// post.Comments = []Comment{}
	err = DB.First(&post).Error

	if err != nil {
		return
	}
	// err = DB.Find(&post.Comments, post.Id).Error
	// if err != nil {
	// 	return
	// }

	return
}
func GetAllPost() (posts []Post, err error) {
	// p := []Post{}
	err = DB.Limit(5).Find(&posts).Error
	// for p, _ := range posts {
	// 	fmt.Println(p)
	// }

	// for _, post := range p {
	// 	Commnets := []Comment{}
	// 	DB.Where("post_id=?", post.ID).Find(&Commnets)
	// 	post.Comments = append(post.Comments, Commnets...)
	// 	posts = append(posts, post)
	// }
	return

}

func Createpost(c *gin.Context) {
	type Input struct {
		Content string
		UserId  int
	}
	var input Input

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post := Post{Content: input.Content, UserId: input.UserId}
	err := DB.Create(&post).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
