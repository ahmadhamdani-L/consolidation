package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"
	"time"

	"gorm.io/gorm"
)

type ValidationDetailEntity struct {
	CompanyID      int    `json:"company_id" validate:"required,numeric"`
	Period         string `json:"period" validate:"required"`
	Versions       int    `json:"versions" validate:"required,numeric"`
	ValidateBy     int    `json:"validate_by" validate:"required,numeric"`
	UserValidateBy string `json:"user_validate" gorm:"-"`
	Name           string `json:"name"`
	Note           string `json:"note"`
	Status         int    `json:"status" validate:"required,numeric"`
}

type ValidationDetailFilter struct {
	ID         *int    `query:"id" validate:"numeric"`
	Period     *string `query:"period" filter:"DATESTRING"`
	Name       *string `query:"name"`
	Versions   *int    `query:"versions" validate:"numeric"`
	ValidateBy *int    `query:"validate_by" validate:"numeric"`
	Status     *int    `query:"status" validate:"numeric"`
}

type ValidationDetailEntityModel struct {
	// abstraction
	ID           int             `json:"id" gorm:"primaryKey;autoIncrement;"`
	ModifiedAt   *time.Time      `json:"modified_at"`
	UserValidate UserEntityModel `json:"-" gorm:"foreignKey:ValidateBy"`
	// entity
	ValidationDetailEntity

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

func (ValidationDetailEntityModel) TableName() string {
	return "validation_detail"
}

type ValidationDetailFilterModel struct {
	// abstraction
	// abstraction.Filter

	// filter
	ValidationDetailFilter
	CompanyCustomFilter
}

func (m *ValidationDetailEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	return
}

func (m *ValidationDetailEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	return
}
