package main

import (
	"database/sql"
	"fmt"
	mm "github.com/felipedsf/go-payment/money_movement/internal/impl"
	pb "github.com/felipedsf/go-payment/money_movement/proto"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"log"
	"net"
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
	grpcServer := grpc.NewServer()
	pb.RegisterMoneyMovementServiceServer(grpcServer, mm.NewGrpcMoneyMovement(db))

	// listen & serve
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Failed to listen on port 8000: %v", err)
	}

	log.Printf("server listening at %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}

	log.Print("Money movement is working!!!")
}
