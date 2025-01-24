package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"wallet/internal/app/services"
	_ "wallet/internal/app/storages/dto"
	"wallet/internal/app/utils"
)

var userService services.UserService

// @Summary Register a new user
// @Description Creates a new user in the system
// @Tags User
// @Accept json
// @Produce json
// @Param user body dto.UserDto true "Registration data"
// @Success 201 {string} string "User registered"
// @Failure 401 {string} string "Unauthorized"
// @Router /api/v1/user/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	body := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	response := userService.Register(h.DB, body["username"].(string), body["password"].(string), body["email"].(string))
	w.WriteHeader(response.Status)
	json.NewEncoder(w).Encode(map[string]string{"message": response.Message.(string)})
	return
}

// @Summary Login User
// @Description Authenticates a user and provides a token
// @Tags User
// @Accept json
// @Produce json
// @Param user body dto.UserDto true "User credentials"
// @Success 201 {string} string "User registered"
// @Failure 401 {string} string "Unauthorized"
// @Router /api/v1/user/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	body := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}
	response := userService.Login(h.DB, body["username"].(string), body["password"].(string))
	w.WriteHeader(response.Status)
	json.NewEncoder(w).Encode(map[string]string{"message": response.Message.(string)})
	return
}

func (h *Handler) initToken() {
	h.tokenAuth = utils.GenerateToken()
}
