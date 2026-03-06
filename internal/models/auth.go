package models

type LoginRequest struct {
    Username string `json:"username" example:"johndoe"`
    Password string `json:"password" example:"password123"`
}