package database

import (
	"blockchain/models"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	db, err := gorm.Open(mysql.Open(os.Getenv("DATABASE_URL")), &gorm.Config{
		SkipDefaultTransaction: true, // Gereksiz transactionları önle
		PrepareStmt:            true, // Statement cache aktif
	})

	if err != nil {
		fmt.Println("failed to connect database" + err.Error())
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
}

func Migrate() {
	DB.AutoMigrate(&models.Block{}, &models.Transaction{})
}
