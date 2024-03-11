package models

import (
	"github.com/jinzhu/gorm"
)

type Follow struct {
	gorm.Model
	FollowedId  int
	FollowingId int
	// Comments []Comment
}
