package model

import (
	"time"
	"worker-validation/internal/abstraction"
	"worker-validation/pkg/util/date"

	"gorm.io/gorm"
)

type ValidationDetailEntity struct {
	CompanyID  int    `json:"company_id"`
	Period     string `json:"period"`
	Versions   int    `json:"versions"`
	ValidateBy int    `json:"validate_by"`
	Name       string `json:"name"`
	Note       string `json:"note"`
	Status     int    `json:"status"`
}

type ValidationDetailFilter struct {
	ID             *int    `query:"id"`
	TrialBalanceID *int    `query:"validation_id"`
	CompanyID      *int    `query:"company_id"`
	Period         *string `query:"period"`
	Versions       *int    `query:"versions"`
	ValidateBy     *int    `query:"validate_by"`
	Status         *int    `query:"status"`
	Name           *string `query:"name"`
}

type ValidationDetailEntityModel struct {
	// abstraction
	ID         int        `json:"id" gorm:"primaryKey;autoIncrement;"`
	ModifiedAt *time.Time `json:"modified_at"`

	// entity
	ValidationDetailEntity

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type ValidationDetailFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	ValidationDetailFilter
	CompanyCustomFilter
}

func (ValidationDetailEntityModel) TableName() string {
	return "validation_detail"
}

func (m *ValidationDetailEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	return
}

func (m *ValidationDetailEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	return
}
