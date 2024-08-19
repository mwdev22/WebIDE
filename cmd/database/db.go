package database

import (
	"fmt"
	"log"

	"github.com/mwdev22/WebIDE/cmd/storage"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DbOpen(connString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to the database:", err)
	}

	fmt.Println("connected to PostgreSQL database...")
	return db, nil
}

func InitConn(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get DB from GORM:", err)
	}
	db.AutoMigrate(models...)
	err = sqlDB.Ping()
	if err != nil {
		sqlDB.Close()
		log.Fatal(err)
	}
}

var models = []interface{}{
	&storage.User{},
	&storage.File{},
	&storage.Repository{},
}
