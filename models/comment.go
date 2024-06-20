package models

import "github.com/jinzhu/gorm"

type Comment struct {
	gorm.Model
	Content string
	Author  string
	// Post    *Post
	PostId int
}
