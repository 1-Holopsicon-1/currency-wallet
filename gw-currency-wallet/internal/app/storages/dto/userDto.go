package dto

type UserDto struct {
	Username string `json:"username" example:"User1"`
	Password string `json:"password" example:"password1"`
	Email    string `json:"email"    example:"user1@example.com"`
}
