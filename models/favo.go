package models

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Favo struct {
	gorm.Model
	PostId int `gorm:"size:255;" json:"postid"`
	UserId int `gorm:"size:255;" json:"userid"`
}

func IsFavo(userid uint, postid uint) bool {
	var favo Favo
	err := DB.First(&favo, "user_id=? AND post_id=?", userid, postid).Error
	// err := DB.First(&favo, "user_id=? AND post_id=?", userid, postid).Error
	fmt.Println(err)
	return err == nil
}

func AddFavo(c *gin.Context) {
	type Input struct {
		UserId int
		PostId int
	}
	var input Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(input)
	favo := Favo{UserId: input.UserId, PostId: input.PostId}

	err := DB.Create(&favo).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func DeleteFavo(c *gin.Context) {

	type Input struct {
		UserId int
		PostId int
	}
	var input Input

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	favo := Favo{UserId: input.UserId, PostId: input.PostId}
	DB.Where("post_id = ? AND user_id = ?", favo.PostId, favo.UserId).Delete(&favo)
}
