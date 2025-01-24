package db

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"wallet/internal/app/storages/entities"
)

func Connect() *gorm.DB {
	e := godotenv.Load(".env")
	if e != nil {
		panic("No env")
	}
	user := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	port := os.Getenv("db_port")

	dsn := fmt.Sprintf(`
	host=%s
	user=%s
	password=%s
	dbname=%s
	port=%s`,
		dbHost, user, password, dbName, port)
	db, err := gorm.Open(postgres.Open(string(dsn)),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Println("Fail to connect to DB")
	}

	return db
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&entities.User{})
	if err != nil {
		log.Fatalln("Fail to migrate User error: ", err)
	}
	err = db.AutoMigrate(&entities.Wallet{})
	if err != nil {
		log.Fatalln("Fail to migrate User error: ", err)
	}
}
