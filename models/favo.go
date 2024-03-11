package models

import (
	"github.com/jinzhu/gorm"
)

type Favo struct {
	gorm.Model
	PostId int
	UserId int
}
