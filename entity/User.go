package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"-" binding:"required,email"`
	Name      string    `json:"name" binding:"required"`
	Age       int       `json:"age"`
	Password  string    `json:"-" binding:"required"`
	CreatedAt time.Time `json:"-"`
}

func NewUser() *User {
	return &User{}
}
