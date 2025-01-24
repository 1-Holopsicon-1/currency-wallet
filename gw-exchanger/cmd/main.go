package main

import (
	"exchanger/internal/app/service"
	"exchanger/internal/app/storages/db"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	proto "proto-exchange/gen/exchange"
)

func main() {
	migr := flag.Bool("migrate", false, fmt.Sprint("Migrating Entity"))
	startGrpc := flag.Bool("startGrpc", false, fmt.Sprint("Starting server"))
	flag.Parse()

	if *migr {
		session := db.Connect()
		log.Println("Migrating")
		db.Migrate(session)

	}

	if *startGrpc {
		session := db.Connect()
		exchangeService := service.ExchangeService{Db: session}
		grpcServer := grpc.NewServer()
		proto.RegisterExchangeServiceServer(grpcServer, &exchangeService)
		defer log.Println("End of Program")
		listener, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen on port 50051: %v", err)
		}
		log.Println("gRPC server is running on port 50051")
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}

}
