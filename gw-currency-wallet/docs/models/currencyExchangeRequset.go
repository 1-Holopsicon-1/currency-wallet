package models

type CurrencyExchangeRequset struct {
	FromCurrency string  `json:"from_currency" example:"USD"`
	ToCurrency   string  `json:"to_currency" example:"EUR"`
	Amount       float64 `json:"amount" example:"100"`
}
