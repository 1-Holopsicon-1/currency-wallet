package dto

type ExchangeDto struct {
	NameFrom string  `json:"nameFrom" gorm:"column:currency_from"`
	NameTo   string  `json:"nameTo" gorm:"column:currency_to"`
	Rate     float64 `json:"rate"`
}
