package models

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDataBase() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Tokyo", dbHost, dbUser, dbPass, dbName)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Could not connect to the database", err)
	}
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Post{})
	DB.AutoMigrate(&Comment{})
}
