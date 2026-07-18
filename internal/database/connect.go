package database

import (
	"fmt"
	"os"
	"strconv"

	"github.com/imrishabk/goauth/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() {
	var err error
	p := os.Getenv("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		panic("failed to parse database port")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		port,
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to database")
	if err := DB.AutoMigrate(&model.User{}); err == nil {
		fmt.Println("Database migrated")
	} else {
		fmt.Println("Failed to migrate database:", err)
	}
}
