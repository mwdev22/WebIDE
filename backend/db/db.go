package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DbOpen(connString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	fmt.Println("Connected to PostgreSQL database...")
	return db, nil
}

func InitConn(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get DB from GORM:", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		sqlDB.Close()
		log.Fatal(err)
	}
}
