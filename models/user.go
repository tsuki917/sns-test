package models

import (
	"strings"

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

func (u User) PrepareOutput() User {
	u.Password = ""
	return u
}

func GenerateToken(email string, password string) (string, error) {
	var user User

	err := DB.Where("username = ?", email).First(&user).Error

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
