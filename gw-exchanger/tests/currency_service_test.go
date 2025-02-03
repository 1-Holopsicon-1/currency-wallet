package tests

import (
	"context"
	"exchanger/internal/app/service"
	"exchanger/internal/app/storages/db"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	proto "proto-exchange/gen/exchange"
	"testing"
)

func TestGetExchageRates(t *testing.T) {
	ctx := context.Background()
	session := connTest()
	var empty *proto.Empty
	exchangeServie := service.ExchangeService{Db: session}
	_, err := exchangeServie.GetExchangeRates(ctx, empty)
	assert.Nil(t, err)
}

func TestGetExchangeRateForCurrency(t *testing.T) {
	ctx := context.Background()
	session := connTest()
	req := proto.CurrencyRequest{
		FromCurrency: "USD",
		ToCurrency:   "EUR",
	}
	exchangeService := service.ExchangeService{Db: session}
	_, err := exchangeService.GetExchangeRateForCurrency(ctx, &req)
	assert.Nil(t, err)
	req.FromCurrency = "NOT_EXIST"
	answer, err := exchangeService.GetExchangeRateForCurrency(ctx, &req)
	assert.Equal(t, 0.0, answer.Rate)
}

func connTest() *gorm.DB {
	err := godotenv.Load("test.env")
	if err != nil {
		panic("Error loading .env file")
	}
	session := db.Connect()
	db.Migrate(session)
	return session
}
