package models

import "time"

type User struct {
    ID        uint      `json:"id" example:"1"`
    Username  string    `json:"username" example:"johndoe"`
    Password  string    `json:"password,omitempty" example:"secret123"`
    CreatedAt time.Time `json:"created_at" example:"2026-03-07T00:00:00Z"`
    UpdatedAt time.Time `json:"updated_at" example:"2026-03-07T00:00:00Z"`
}