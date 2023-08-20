package entity

import (
	"time"

	"github.com/google/uuid"
)

type Phone struct {
	ID          uuid.UUID `json:"-"`
	Number      string    `json:"phone" binding:"required"`
	Description string    `json:"description,omitempty"`
	User        User      `json:"-"`
	UserId      uuid.UUID `json:"user_id"`
	IsFax       bool      `json:"is_fax"`
	CreatedAt   time.Time `json:"-"`
}

func NewPhone() *Phone {
	return &Phone{}
}
