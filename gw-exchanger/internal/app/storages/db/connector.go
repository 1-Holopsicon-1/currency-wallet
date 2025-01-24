package db

import (
	"exchanger/internal/app/storages/entities"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

func Connect() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
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
	err := db.AutoMigrate(&entities.Rates{})
	err = db.AutoMigrate(&entities.Currency{})
	if err != nil {
		log.Println(err)
		panic("Fail to migrate")
	}
	db.Save(&entities.Currency{Name: "usd"})
	db.Save(&entities.Currency{Name: "eur"})
	db.Save(&entities.Currency{Name: "rub"})
	db.Save(&entities.Rates{CurrencyFrom: 1, CurrencyTo: 2, Rate: 0.01})
	db.Save(&entities.Rates{CurrencyFrom: 1, CurrencyTo: 3, Rate: 0.05})
	db.Save(&entities.Rates{CurrencyFrom: 2, CurrencyTo: 3, Rate: 0.5})

}
