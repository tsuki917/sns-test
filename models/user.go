package models

import (
	"fmt"
	"io"
	fire "sns-test/utils/firebase"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/wazyiz/jwt-gin/utils/token"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null;" json:"password"`
	UserTag  string `gorm:"size:255;not null;" json:"usertag"`
	ImgPath  string `gorm:"size:255;" json:"img_path"`
	Profile  string `gorm:"size:255;" json:"profile"`
}

type User_Post struct {
	Username string `gorm:"size:255;not null;unique" json:"username"`
	UserTag  string `gorm:"size:255;not null;" json:"usertag"`
	ImgPath  string `gorm:"size:255;" json:"img_path"`
	ID       uint   `gorm:"size:255;" json:"id"`
}

func (u User) Save() (User, error) {
	err := DB.Create(&u).Error

	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (u *User) BeforeSave() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	u.Username = strings.ToLower(u.Username)

	return nil
}

type output struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	UserTag  string `gorm:"size:255;not null;" json:"usertag"`
	ImgPath  string `gorm:"size:255;" json:"img_path"`
	Profile  string `gorm:"size:255;" json:"profile"`
}

func GetOtherUser(UserId int) (userData output) {
	var user User
	err := DB.Model(&user).Where("id=?", UserId).First(&userData)
	if err != nil {
		fmt.Print(err)
	}
	return
}

func (u User) UpdateUser(UserId int, c *gin.Context, isFileUpdate bool) {
	fmt.Println("UpdateUser")
	var err error
	fmt.Println(isFileUpdate)
	newUser := User{Username: u.Username, UserTag: u.UserTag, Profile: u.Profile}
	fmt.Println("is fifile up load?")

	if isFileUpdate {
		fmt.Println("file up load")
		filedata, header, err := c.Request.FormFile("image")
		if err != nil {
			fmt.Println(err)
		}
		// filename := header.Filename
		preUser := User{}
		err = DB.Where("id=?", UserId).Find(&preUser).Error
		if err != nil {
			fmt.Println(err)
		}
		h := header.Header["Content-Type"][0]
		t := strings.Replace(h, "image/", "", 1)
		fmt.Println(h)
		binaryData, err := io.ReadAll(filedata)
		if err != nil {
			fmt.Println(err)
		}

		// prepath := strings.Replace(preUser.ImgPath, "https://storage.googleapis.com/sns-image-storage.appspot.com/", "", 1)
		nextpath := "IconImg/" + strconv.Itoa(UserId) + "/" + "prof." + t
		// fire.DeleteIconImg(prepath)
		fire.SavePostImg(binaryData, nextpath)
		p := "https://storage.googleapis.com/sns-image-storage.appspot.com/" + nextpath
		newUser.ImgPath = p

	}

	err = DB.Model(&User{}).Where("id=?", UserId).Updates(newUser).Error
	if err != nil {
		fmt.Println(err)
	}
}

func (u User) PrepareOutput() User {
	u.Password = ""
	return u
}

func GenerateToken(username string, password string) (string, error) {
	var user User

	err := DB.Where("username = ?", username).First(&user).Error

	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return "", err
	}

	token, err := token.GenerateToken(user.ID)

	if err != nil {
		return "", err
	}

	return token, nil
}
