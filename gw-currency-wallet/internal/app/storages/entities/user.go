package entities

type User struct {
	Id       int64  `gorm:"primary_key;AUTO_INCREMENT"`
	Username string `gorm:"unique"`
	Password string `gorm:"size:255"`
	Email    string `gorm:"unique"`
}
