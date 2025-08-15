package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username" gorm:"not null;type:string"`
	Password string    `json:"password" gorm:"not null;type:string"`
	Mail     string    `json:"mail"`
	Role     string    `json:"role" gorm:"not null;type:string"`
	Blocked  bool      `json:"blocked"`
}

func (user *User) BeforeCreate(scope *gorm.DB) error {
	user.ID = uuid.New()
	return nil
}
