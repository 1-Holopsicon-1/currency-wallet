package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	proto "proto-exchange/gen/exchange"
	"strings"
	"time"
	mychace "wallet/internal/app/cache"
	"wallet/internal/app/storages/dto"
	"wallet/internal/app/storages/entities"
	protoExchange "wallet/internal/proto/client"
)

type WalletService struct{}

func (service *WalletService) GetBalance(db *gorm.DB, claims map[string]interface{}) (dto.WalletDto, error) {
	var wallet entities.Wallet

	if err := db.First(&wallet, "user_id = ?", claims["userId"].(string)).Error; err != nil {
		log.Printf("Error getting balance: %v", err)
		return dto.WalletDto{}, fmt.Errorf("error getting balance")
	}
	walletDto := dto.WalletDto{
		Usd: wallet.Usd,
		Rub: wallet.Rub,
		Eur: wallet.Eur,
	}
	return walletDto, nil
}

func (service *WalletService) Deposit(db *gorm.DB, claims map[string]interface{}, amount float64, currency string) (dto.WalletDto, error) {
	var walletOld map[string]interface{}
	query := `
    select * from wallets where user_id = ?   
    `
	if err := db.Raw(query, claims["userId"].(string)).Scan(&walletOld).Error; err != nil {
		log.Printf("Error getting balance: %v", err)
		return dto.WalletDto{}, fmt.Errorf("error getting balance")
	}

	currency = strings.ToLower(currency)
	fmt.Println(walletOld)
	walletOld[currency] = walletOld[currency].(float64) + amount
	walletNew := entities.Wallet{Id: int(walletOld["id"].(int64)),
		Usd:    walletOld["usd"].(float64),
		Rub:    walletOld["rub"].(float64),
		Eur:    walletOld["eur"].(float64),
		UserId: walletOld["user_id"].(int64),
	}
	db.Save(&walletNew)
	walletDto := dto.WalletDto{Usd: walletNew.Usd, Rub: walletNew.Rub, Eur: walletNew.Eur}
	return walletDto, nil
}

func (service *WalletService) Withdraw(db *gorm.DB, claims map[string]interface{}, amount float64, currency string) (dto.WalletDto, error) {
	var walletOld map[string]interface{}

	query := `
    select * from wallets where user_id = ?   
    `
	if err := db.Raw(query, claims["userId"].(string)).Scan(&walletOld).Error; err != nil {
		return dto.WalletDto{}, err
	}
	currency = strings.ToLower(currency)
	if walletOld[currency].(float64)-amount < 0 {
		return dto.WalletDto{}, fmt.Errorf("insufficient funds")
	}
	walletOld[currency] = walletOld[currency].(float64) - amount
	walletNew := entities.Wallet{Id: int(walletOld["id"].(int64)),
		Usd:    walletOld["usd"].(float64),
		Rub:    walletOld["rub"].(float64),
		Eur:    walletOld["eur"].(float64),
		UserId: walletOld["user_id"].(int64),
	}
	db.Save(&walletNew)
	walletDto := dto.WalletDto{Usd: walletNew.Usd, Rub: walletNew.Rub, Eur: walletNew.Eur}
	return walletDto, nil
}

func (service *WalletService) GetExchangeRates() (*proto.ExchangeRatesResponse, error) {
	ctx := context.Background()
	var CacheManager = cache.New[any](mychace.InitRedis())
	if value, err := CacheManager.Get(ctx, "rates"); !errors.Is(err, redis.Nil) {
		response := &proto.ExchangeRatesResponse{}
		err := json.Unmarshal([]byte(value.(string)), &response.Rates)
		if err != nil {
			log.Printf("Error of Unmarshal data", err)
			return nil, fmt.Errorf("Cash error")
		}
		return response, nil
	}
	client, err := protoExchange.Connect()
	if err != nil {
		log.Println("Error connecting to exchange")
		return nil, fmt.Errorf("error connecting to exchange")
	}
	rates, err := protoExchange.GetAllExchange(client)
	if err != nil {
		log.Println("Error getting exchange rates")
		return nil, fmt.Errorf("error getting exchange rates")
	}
	cacheData, _ := json.Marshal(rates.Rates)
	err = CacheManager.Set(ctx, "rates", cacheData, store.WithExpiration(5*time.Minute))
	if err != nil {
		log.Printf("Error set cache rates: %v", err)
	}

	return rates, nil
}

func (service *WalletService) ChangeCurrency(db *gorm.DB, claims map[string]interface{},
	currencyFrom, currencyTo string, amount float64) (map[string]interface{}, error) {
	ctx := context.Background()
	var CacheManager = cache.New[any](mychace.InitRedis())
	var balance map[string]interface{}
	currencyFrom = strings.ToLower(currencyFrom)
	currencyTo = strings.ToLower(currencyTo)
	keyCache := fmt.Sprintf("currency_from_%s_to_%s", currencyFrom, currencyTo)
	value, err := CacheManager.Get(ctx, keyCache)
	if errors.Is(err, redis.Nil) {
		client, err := protoExchange.Connect()
		if err != nil {
			log.Println("Error connecting to exchange")
			return nil, fmt.Errorf("error connecting to exchange")
		}
		currency, err := client.GetExchangeRateForCurrency(ctx, &proto.CurrencyRequest{FromCurrency: currencyFrom, ToCurrency: currencyTo})
		fmt.Println(currency)
		if err != nil {
			log.Printf("Error getting exchange rates: %v", err)
			return nil, fmt.Errorf("error getting exchange rates")
		}
		cacheData, err := json.Marshal(currency)
		if err != nil {
			log.Printf("Error of Marshal data", err)
			return nil, fmt.Errorf("error marshaling data")
		}
		err = CacheManager.Set(ctx, keyCache, cacheData, store.WithExpiration(5*time.Minute))
		if err != nil {
			log.Printf("Error set cache rates: %v", err)
			return nil, fmt.Errorf("error rates")
		}
		value = string(cacheData)
	}
	var ratesFromTo proto.ExchangeRateResponse
	err = json.Unmarshal([]byte(value.(string)), &ratesFromTo)
	if err != nil {
		log.Printf("Error of Unmarshal data %v", err)
		return nil, fmt.Errorf("Error getting exchange rates")
	}
	db.Model(entities.Wallet{}).First(&balance, "user_id = ?", claims["userId"].(string))
	if balance[currencyFrom] = balance[currencyFrom].(float64) - amount; balance[currencyFrom].(float64) < 0 {
		return nil, fmt.Errorf("insufficient funds")
	}
	balance[currencyTo] = balance[currencyTo].(float64) + (amount * ratesFromTo.Rate)
	db.Model(entities.Wallet{}).Where("user_id = ?", balance["user_id"]).Updates(balance)
	response := make(map[string]interface{})
	response["new_balance"] = map[string]interface{}{
		currencyFrom: balance[currencyFrom].(float64),
		currencyTo:   balance[currencyTo].(float64),
	}
	response["exchanged_amount"] = balance[currencyTo].(float64)
	return response, nil
}
