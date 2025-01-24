package entities

type Currency struct {
	Id   int `gorm:"primaryKey;"`
	Name string
}
