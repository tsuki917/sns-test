package models

import (
	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model
	Content string
	UserTag string
	FavoId  int
	UserId  int
	// Comments []Comment
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
	// post.Comments = []Comment{}
	p := []Post{}
	DB.Limit(5).Find(&p)

	// for _, post := range p {
	// 	Commnets := []Comment{}
	// 	DB.Where("post_id=?", post.ID).Find(&Commnets)
	// 	post.Comments = append(post.Comments, Commnets...)
	// 	posts = append(posts, post)
	// }
	return

}
