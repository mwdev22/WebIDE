package main

import (
	"fmt"
	"log"

	"github.com/mwdev22/WebIDE/api"
	database "github.com/mwdev22/WebIDE/db"
	"github.com/mwdev22/WebIDE/utils"
)

func main() {

	dbCfg := utils.GetDbCfg()
	utils.LoadEnv()

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=5432 sslmode=disable",
		dbCfg.User, dbCfg.Name, dbCfg.Pass, dbCfg.Host)

	db, err := database.DbOpen(connStr)
	if err != nil {
		fmt.Printf("db open failed: %v", err)
		return
	}
	database.InitConn(db)

	server := api.NewServer(":8080")
	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
