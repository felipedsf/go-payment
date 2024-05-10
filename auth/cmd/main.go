package main

import (
	"database/sql"
	"fmt"
	"github.com/felipedsf/go-payment/auth/internal/impl/auth"
	pb "github.com/felipedsf/go-payment/auth/proto"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	dbDriver = "mysql"
	dbUser   = "root"
	dbPass   = "root"
	dbName   = "users"
)

var db *sql.DB

func main() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", dbUser, dbPass, dbName)

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
	authServer := auth.NewGrpcAuth(db)
	pb.RegisterAuthServiceServer(grpcServer, authServer)

	// listen & serve
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}

	log.Printf("server listening at %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}

	log.Print("Auth is working!!!")
}
