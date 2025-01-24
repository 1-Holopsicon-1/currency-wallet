package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
	_ "wallet/docs/models"
	"wallet/internal/app/services"
)

var walletService services.WalletService

// @Summary Get user balance
// @Description Retrieves the current balance of the user
// @Tags Balance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Success 200 {object} map[string]interface{} "User balance"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/user/balance [get]
func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := make(map[string]interface{})
	_, claims, _ := jwtauth.FromContext(r.Context())
	data, err := walletService.GetBalance(h.DB, claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("message: " + "Internal Server Error")
		return
	}
	w.WriteHeader(http.StatusOK)
	response["balance"] = data
	json.NewEncoder(w).Encode(response)

}

// @Summary Deposit money
// @Description Adds money to the user's balance
// @Tags Balance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param deposit body map[string]float64 true "Deposit amount" default(Bearer )
// @Param Authorization header string true "Bearer token"
// @Success 200 {string} string "Deposit successful"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/user/balance/deposit [post]
func (h *Handler) Deposit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, claims, _ := jwtauth.FromContext(r.Context())
	body := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("error: " + err.Error())
		return
	}
	data, err := walletService.Deposit(h.DB, claims, body["amount"].(float64), body["currency"].(string))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("error: " + "Internal Server Error")
		return
	}
	response := map[string]interface{}{
		"message":     "Account topped up successfully",
		"new_balance": data,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// @Summary Withdraw money
// @Description Subtracts money from the user's balance
// @Tags Balance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param withdraw body map[string]float64 true "Withdrawal amount" default(Bearer )
// @Param Authorization header string true "Bearer token"
// @Success 200 {string} string "Withdrawal successful"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/user/balance/withdraw [post]
func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, claims, _ := jwtauth.FromContext(r.Context())
	body := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("error: " + err.Error())
	}
	data, err := walletService.Withdraw(h.DB, claims, body["amount"].(float64), body["currency"].(string))
	fmt.Println(errors.Is(err, fmt.Errorf("insufficient funds")))
	if err != nil {
		if err.Error() == "insufficient funds" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("error: " + err.Error())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("error: " + "Internal Server Error")
		return
	}
	response := map[string]interface{}{
		"message":     "Withdrawal successful",
		"new_balance": data,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// @Summary Get exchange rates
// @Description Retrieves current exchange rates
// @Tags Exchange
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Success 200 {object} map[string]interface{} "Exchange rates"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/exchange/rates [get]
func (h *Handler) GetExchangeRates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := make(map[string]interface{})

	_, _, _ = jwtauth.FromContext(r.Context())
	data, err := walletService.GetExchangeRates()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("error: " + "Internal Server Error")
		return
	}
	response["rates"] = data.Rates
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

// @Summary Exchange currencies
// @Description Exchange entered currencies
// @Tags Exchange
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Param withdraw body models.CurrencyExchangeRequset true "Withdrawal amount"
// @Success 200 {object} map[string]interface{} "Exchange rates"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Failure 400 {string} string "Insufficient funds or invalid currencies"
// @Router /api/v1/exchange [post]
func (h *Handler) ExchangeCurrencies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := make(map[string]interface{})
	body := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("error: " + err.Error())
	}
	_, claims, _ := jwtauth.FromContext(r.Context())
	data, err := walletService.ChangeCurrency(h.DB, claims, body["from_currency"].(string),
		body["to_currency"].(string), body["amount"].(float64))
	if err != nil {
		if err.Error() == "insufficient funds" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("error: " + "Insufficient funds or invalid currencies")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("error: " + "Internal Server Error")
	}

	w.WriteHeader(http.StatusOK)
	response = data
	response["message"] = "Exchange successful"
	json.NewEncoder(w).Encode(response)
}
