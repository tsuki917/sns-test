package controllers

import (
	"fmt"
	"net/http"
	"test-sns/models"
	"test-sns/utils/token"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type RegisterInput struct {
	gorm.Model
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	UserId   string `json:"userid" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{Username: input.Username, Password: input.Password, UserId: input.UserId, Email: input.Email}
	user.BeforeSave()
	fmt.Print(user)
	user, err := user.Save()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user.PrepareOutput(),
	})
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := models.GenerateToken(input.Username, input.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func CurrentUser(c *gin.Context) {
	userId, err := token.ExtractTokenId(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	err = models.DB.First(&user, userId).Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user.PrepareOutput(),
	})
}
