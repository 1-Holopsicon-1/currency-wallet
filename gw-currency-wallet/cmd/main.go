package main

import (
	"flag"
	"fmt"
	"log"
	"wallet/internal/app/handler"
	"wallet/internal/app/server"
	"wallet/internal/app/storages/db"
)

// @title Wallet Api
// @version 0.1
// @description All endpoints Wallet (Rates, Balance), User

// @host localhost:5000

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	migr := flag.Bool("migrate", false, fmt.Sprint("Migrating Entity"))
	start := flag.Bool("start", false, fmt.Sprint("Starting server"))

	flag.Parse()

	if *migr {
		session := db.Connect()
		log.Println("Migrating")
		db.Migrate(session)
	}
	if *start {
		srv := new(server.Server)
		log.Println("Starting server")
		defer log.Println("End of Program")
		log.Println("Open the server")
		fmt.Println("Running and Serving on: http://127.0.0.1:5000")
		session := db.Connect()
		mHandler := handler.Handler{DB: session}
		if err := srv.Run(":5000", mHandler.InitRoutes()); err != nil {
			log.Fatalln(err)
		}
	}
}
