package model

import (
	"notification/internal/abstraction"
	"time"
)

type UserEntity struct {
	Username     string `json:"username" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Password     string `json:"password" validate:"required" gorm:"-"`
	ImageProfile string `json:"image_profile" validate:"required"`
	CompanyID    int    `json:"company_id" validate:"required"`
	PasswordHash string `json:"-" gorm:"column:password"`
}

type UserEntityModel struct {
	// abstraction
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`

	// entity
	UserEntity

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

func (UserEntityModel) TableName() string {
	return "users"
}

// func (m *UserEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
// 	m.ModifiedAt = date.DateTodayLocal()
// 	m.ModifiedBy = &m.Context.Auth.ID
// 	return
// }
