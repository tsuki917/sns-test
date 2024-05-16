package models

import (
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
)

type PostImg struct {
	gorm.Model
	PostId   int64
	FileName string
}

func GetImageURL(PostId int) (imageURL []string) {
	var images []PostImg
	err := DB.Where("post_id=?", PostId).Find(&images).Error
	if err != nil {
		fmt.Print(err)
		return
	}
	parent := "https://storage.googleapis.com/sns-image-storage.appspot.com/postImg/" + strconv.Itoa(int(PostId))
	for _, image := range images {
		path := parent + "/" + image.FileName
		imageURL = append(imageURL, path)
	}
	fmt.Println(imageURL)
	return
}
