package tests

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"testing"
	"time"
	"wallet/internal/app/services"
	"wallet/internal/app/storages/db"
	"wallet/internal/app/storages/entities"
	"wallet/internal/app/utils"
)

func TestTestRegister(t *testing.T) {
	err := godotenv.Load("tests.env")
	if err != nil {
		panic("Error loading .env file")
	}
	session := db.Connect()
	db.ReInit(session)
	session.Debug()
	tester := entities.User{
		Username: "tester",
		Password: "tester1",
		Email:    "tester1@gmail.com",
	}
	tester2 := entities.User{}
	userService := &services.UserService{}
	resp := userService.Register(session, tester.Username, tester.Password, tester.Email)
	session.First(&tester2, "username = ?", tester.Username)
	assert.Equal(t, http.StatusCreated, resp.Status)
	assert.Equal(t, "User created", resp.Message)
	assert.NotEqual(t, tester.Id, tester2.Id)
	assert.Equal(t, nil, bcrypt.CompareHashAndPassword([]byte(tester2.Password), []byte(tester.Password)))
	resp = userService.Register(session, tester.Username, tester.Password, tester.Email)
	assert.Equal(t, http.StatusBadRequest, resp.Status)
	assert.Equal(t, "Username or email already exists", resp.Message)
}

func TestTestLogin(t *testing.T) {
	err := godotenv.Load("tests.env")
	if err != nil {
		panic("Error loading .env file")
	}
	session := db.Connect()
	db.ReInit(session)
	session.Debug()
	tester := entities.User{
		Username: "tester",
		Password: "tester1",
		Email:    "tester1@gmail.com",
	}
	notRegTester := entities.User{
		Username: "tester123123123",
		Password: "tester1123123",
		Email:    "tester1123123123@gmail.com",
	}
	userService := &services.UserService{}
	userService.Register(session, tester.Username, tester.Password, tester.Email)
	session.Select("id").First(&tester, "username = ?", tester.Username)
	resp := userService.Login(session, tester.Username, tester.Password)
	assert.Equal(t, http.StatusOK, resp.Status)
	assert.Equal(t, generateStringToken(tester.Id), resp.Message)
	resp = userService.Login(session, tester.Username, notRegTester.Password)
	assert.Equal(t, http.StatusUnauthorized, resp.Status)
	assert.Equal(t, "Wrong password", resp.Message)
	resp = userService.Login(session, notRegTester.Username, notRegTester.Password)
	assert.Equal(t, http.StatusUnauthorized, resp.Status)
	assert.Equal(t, "No such user", resp.Message)
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
