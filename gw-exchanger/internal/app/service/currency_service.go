package service

import (
	"context"
	"exchanger/internal/app/storages/dto"
	"fmt"
	"gorm.io/gorm"
	"math"
	proto "proto-exchange/gen/exchange"
	"strconv"
)

type ExchangeService struct {
	proto.UnimplementedExchangeServiceServer
	Db *gorm.DB
}

func (exchangeService *ExchangeService) GetExchangeRates(ctx context.Context, empty *proto.Empty) (*proto.ExchangeRatesResponse, error) {
	output := make(map[string]float64)
	exchDto, err := exchangeService.currentAllExchange(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get All exchange", err)
	}
	for _, v := range exchDto {
		output[v.NameFrom+"_to_"+v.NameTo] = float64(v.Rate)
	}
	answer := proto.ExchangeRatesResponse{Rates: output}
	return &answer, err
}

func (exchangeService *ExchangeService) GetExchangeRateForCurrency(ctx context.Context, request *proto.CurrencyRequest) (*proto.ExchangeRateResponse, error) {
	exchange, err := exchangeService.currentExchange(ctx, request.FromCurrency, request.ToCurrency)
	if err != nil {
		return nil, err
	}
	answer := proto.ExchangeRateResponse{FromCurrency: exchange.NameFrom, ToCurrency: exchange.NameTo, Rate: exchange.Rate}
	return &answer, nil
}

func (exchangeService *ExchangeService) currentAllExchange(ctx context.Context) ([]dto.ExchangeDto, error) {
	var currenciesDto []dto.ExchangeDto
	query := `
	SELECT cf.name as currency_from, ct.name as currency_to, rt.rate
	FROM rates rt
         join currencies cf on rt.currency_from = cf.id
         join currencies ct on rt.currency_to = ct.id
	`
	err := exchangeService.Db.WithContext(ctx).Raw(query).Scan(&currenciesDto).Error
	if err != nil {
		return currenciesDto, err
	}
	return currenciesDto, nil
}

func (exchangeService *ExchangeService) currentExchange(ctx context.Context, fcur, tcur string) (dto.ExchangeDto, error) {
	var currenciesDto dto.ExchangeDto

	query := `SELECT name as currency_from, 
    (select name from currencies where name = ?) as currency_to, rate
    from currencies cur 
	join rates on cur.id = rates.currency_from 
	where cur.name = ? and rates.currency_to = (select id from currencies where name = ?);`

	err := exchangeService.Db.WithContext(ctx).Raw(query, tcur, fcur, tcur).Scan(&currenciesDto).Error
	if err != nil {
		return dto.ExchangeDto{}, fmt.Errorf("can't get current exchange")
	}

	if (dto.ExchangeDto{}) == currenciesDto {
		err = exchangeService.Db.WithContext(ctx).Raw(query, fcur, tcur, fcur).Scan(&currenciesDto).Error
		v, err := strconv.ParseFloat(fmt.Sprintf("%.2f", 1/currenciesDto.Rate), 64)
		if err != nil {
			return dto.ExchangeDto{}, fmt.Errorf("Convert error")
		}
		currenciesDto.Rate = v
	}
	if currenciesDto.Rate == math.Inf(0) {
		currenciesDto.Rate = 0
	}
	return currenciesDto, nil
}
