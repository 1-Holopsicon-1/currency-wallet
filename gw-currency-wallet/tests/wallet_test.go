package tests

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"strconv"
	"testing"
	"wallet/internal/app/services"
	"wallet/internal/app/storages/db"
	"wallet/internal/app/storages/dto"
	"wallet/internal/app/storages/entities"
)

var usService services.UserService

func TestGetBalance(t *testing.T) {
	session := connTest()
	userCr := userGen(session)
	var walletService services.WalletService
	resp, err := walletService.GetBalance(session, userCr)
	walletDto := dto.WalletDto{
		Usd: 0,
		Rub: 0,
		Eur: 0,
	}
	assert.Equal(t, walletDto, resp)
	assert.Equal(t, nil, err)
	userCr["userId"] = "123"
	resp, err = walletService.GetBalance(session, userCr)
	assert.Equal(t, fmt.Errorf("error getting balance"), err)
	assert.Equal(t, dto.WalletDto{}, resp)
}

func TestDeposit(t *testing.T) {
	session := connTest()
	userCr := userGen(session)
	var walletService services.WalletService
	amount := 100.00
	name := "USD"
	walletDto := dto.WalletDto{
		Usd: amount,
		Rub: 0,
		Eur: 0,
	}
	session.Debug()
	resp, err := walletService.Deposit(session, userCr, amount, name)
	assert.Equal(t, walletDto, resp)
	assert.Equal(t, nil, err)

}

func TestWithdraw(t *testing.T) {
	session := connTest()
	userCr := userGen(session)
	session.Debug()
	amount := 100.00
	name := "USD"
	walletDto := dto.WalletDto{
		Usd: 0,
		Rub: 0,
		Eur: 0,
	}
	session.Model(&entities.Wallet{}).Where("id = ?", userCr["userId"]).Update("usd", amount)
	var walletService services.WalletService
	resp, err := walletService.Withdraw(session, userCr, amount, name)
	assert.Equal(t, walletDto, resp)
	assert.Equal(t, nil, err)
	resp, err = walletService.Withdraw(session, userCr, amount, name)
	assert.Equal(t, dto.WalletDto{}, resp)
	assert.Equal(t, fmt.Errorf("insufficient funds"), err)
}

func TestGetExchangeRates(t *testing.T) {
	err := godotenv.Load("tests.env")
	if err != nil {
		panic("Error loading .env file")
	}
	var walletService services.WalletService
	_, err = walletService.GetExchangeRates()
	assert.Equal(t, nil, err)

}

func TestChangeCurrency(t *testing.T) {
	session := connTest()
	userCr := userGen(session)
	session.Debug()
	amount := 100.00
	session.Model(&entities.Wallet{}).Where("id = ?", userCr["userId"]).Update("usd", amount)
	var walletService services.WalletService
	resp, err := walletService.ChangeCurrency(session, userCr, "USD", "EUR", 100)
	assert.Equal(t, nil, err)
	eqCheck := resp["new_balance"].(map[string]interface{})
	assert.NotEqual(t, eqCheck["currencyFrom"], amount)
	assert.NotEqual(t, eqCheck["currencyTo"], 0)
	_, err = walletService.ChangeCurrency(session, userCr, "USD", "EUR", 100)
	assert.Equal(t, fmt.Errorf("insufficient funds"), err)

}

func connTest() *gorm.DB {
	err := godotenv.Load("tests.env")
	if err != nil {
		panic("Error loading .env file")
	}
	session := db.Connect()
	db.ReInit(session)
	return session
}

func userGen(session *gorm.DB) map[string]interface{} {
	tester := entities.User{
		Username: "tester",
		Password: "tester1",
		Email:    "tester1@gmail.com",
	}
	ret := make(map[string]interface{}, 0)
	usService.Register(session, tester.Username, tester.Password, tester.Email)
	session.Model(&entities.User{}).First(&ret, "username = ?", tester.Username)
	ret["userId"] = strconv.FormatInt(ret["id"].(int64), 10)
	return ret
}
