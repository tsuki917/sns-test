package models

import (
	"github.com/jinzhu/gorm"
)

type Tag struct {
	gorm.Model
	Title  string
	PostId int
}
