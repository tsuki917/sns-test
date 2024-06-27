package models

import (
	"fmt"
	"io"
	"net/http"
	fire "sns-test/utils/firebase"
	"strconv"

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
	err = DB.Limit(10).Find(&posts).Error
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
func GetAllMyPost(myID int) (posts []Post, err error) {
	// p := []Post{}
	err = DB.Limit(10).Where("user_id=?", myID).Find(&posts).Error
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
		Content    string
		UserId     int
		FileLength int64
	}
	var input Input

	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(input)
	var lastValue int64

	// シーケンス番号を取得するクエリを実行
	DB.Raw("SELECT last_value FROM posts_id_seq").Scan(&lastValue)
	lastValue += 1
	// fire.ImagefileUP(c, input.FileLength, lastValue)

	post := Post{Content: input.Content, UserId: input.UserId}
	err := DB.Create(&post).Error
	ImagefileUP(c, input.FileLength, lastValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
func ImagefileUP(c *gin.Context, fileLen int64, seq int64) {
	for i := 0; i < int(fileLen); i++ {
		filedata, header, err := c.Request.FormFile("image[" + strconv.Itoa(i) + "]")
		if err != nil {
			print(err)
			return
		}
		p := PostImg{}
		p.FileName = header.Filename

		p.PostId = seq
		err = DB.Create(&p).Error
		fmt.Println(p)
		if err != nil {
			fmt.Println(err)
		}
		binaryData, err := io.ReadAll(filedata)
		if err != nil {
			print(err)
			return
		}
		path := "postImg/" + strconv.Itoa(int(seq)) + "/" + header.Filename
		fire.SavePostImg(binaryData, path)
	}
}

// func GetImgBase64Array(PostId int) (base64Array []string) {
// 	// var images []PostImg

// 	// err := DB.Where("id=?", PostId).Find(&images).Error
// 	// fmt.Println("GetImage")
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }
// 	// for _, image := range images {
// 	// 	base64Array = append(base64Array, fire.GetPostImg(PostId, image.FileName))
// 	// }
// 	return
// }
