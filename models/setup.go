package models

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to load env file")
	}

	dsn := os.Getenv("DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: nil,
	})

	if err != nil {
		panic("Failed to connect to database!")
	} else {
		println("Connected to database!")
	}

	// Check if flag `--run-migrations` is passed.
	if len(os.Args) > 1 && os.Args[1] == "--run-migrations" {
		fmt.Println("=====================================")
		fmt.Println("========== RUNNING MIGRATIONS =======")
		fmt.Println("=====================================")

		db.AutoMigrate(&User{}, &Board{}, &List{}, &Card{}, &Attachment{}, &Label{})
	}

	return db, err
}
