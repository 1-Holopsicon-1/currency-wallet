package entities

type Wallet struct {
	Id     int `gorm:"primary_key;AUTO_INCREMENT"`
	Usd    float64
	Rub    float64
	Eur    float64
	UserId int64
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
