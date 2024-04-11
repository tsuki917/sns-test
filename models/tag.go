package models

import (
	"github.com/jinzhu/gorm"
)

type Tag struct {
	gorm.Model
	Title  string `gorm:"size:255;" json:"title"`
	PostId int    `gorm:"size:255;" json:"postid"`
}
