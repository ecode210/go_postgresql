package config

import (
	"fmt"
	"go_postgresql/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

var DB *gorm.DB

// ConnectDB - Connects app to postgresql database
func ConnectDB() {
	// Postgresql DSN
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=admin dbname=postgres port=5432"
		fmt.Println("LOCAL DATABASE")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(any("Unable to connect to database"))
	}
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		panic(any("Unable to migrate User struct"))
	}
	DB = db
}
