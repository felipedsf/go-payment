package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const (
	dbDriver = "mysql"
	dbUser   = "root"
	dbPass   = "root"
	dbName   = "money_movement"
)

var db *sql.DB

func main() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3307)/%s", dbUser, dbPass, dbName)

	db, err = sql.Open(dbDriver, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Fatalf("Error closing db: %s", err)
		}
	}()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// grpc server setup
	//grpcServer := grpc.NewServer()

	log.Print("Money movement is working!!!")
}
