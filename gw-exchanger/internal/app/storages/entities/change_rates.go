package entities

type Rates struct {
	Id             int `gorm:"primary_key;column:id"`
	CurrencyFrom   uint
	CurrencyTo     uint
	Rate           float64  `gorm:"column:rate"`
	FromCurrencyID Currency `gorm:"foreignKey:CurrencyFrom"`
	ToCurrencyID   Currency `gorm:"foreignKey:CurrencyTo"`
}
