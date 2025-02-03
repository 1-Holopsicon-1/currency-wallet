package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
	"wallet/internal/app/models"
	"wallet/internal/app/storages/entities"
	"wallet/internal/app/utils"
)

type UserService struct {
}

func (us *UserService) Register(db *gorm.DB, name, pw, email string) models.Response {
	userEntity := entities.User{
		Username: name,
		Email:    email,
	}
	err := db.First(&userEntity, "username = ? or email = ?", userEntity.Username, userEntity.Email).Error
	if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Response{Status: http.StatusBadRequest, Message: "Username or email already exists"}
	}
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password")
		return models.Response{Status: http.StatusBadRequest, Message: "Error hashing password"}
	}
	userEntity.Password = string(hashPwd)
	err = db.Create(&entities.Wallet{Usd: 0, Rub: 0, Eur: 0, User: userEntity}).Error
	if err != nil {
		log.Println("UserService.Register error:", err)
		return models.Response{Status: http.StatusBadRequest, Message: err.Error()}
	}
	return models.Response{Status: http.StatusCreated, Message: "User created"}
}

func (us *UserService) Login(db *gorm.DB, username, pw string) models.Response {
	var userEntity entities.User
	err := db.First(&userEntity, "username = ?", username).Error
	if err != nil {
		log.Println("Login error: no such user")
		return models.Response{Status: http.StatusUnauthorized, Message: "No such user"}
	}

	err = bcrypt.CompareHashAndPassword([]byte(userEntity.Password), []byte(pw))
	if err != nil {
		return models.Response{Status: http.StatusUnauthorized, Message: "Wrong password"}
	}
	return models.Response{Status: http.StatusOK, Message: generateStringToken(userEntity.Id)}
}

func generateStringToken(userId int64) string {
	tokenAuth := utils.GenerateToken()
	claims := map[string]interface{}{
		"userId": userId,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	}
	_, tokenString, _ := tokenAuth.Encode(claims)
	return tokenString
}
