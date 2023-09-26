package model

import (
	"worker/internal/abstraction"
	"worker/pkg/util/date"

	"gorm.io/gorm"
)

type FormatterEntity struct {
	FormatterFor string `json:"formatter_for" validate:"required"`
	Description  string `json:"description" validate:"required"`
}

type FormatterFilter struct {
	FormatterFor *string `query:"formatter_for" filter:"ILIKE"`
	Description  *string `query:"description" filter:"ILIKE"`
}

type FormatterEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	FormatterEntity

	// relations
	FormatterDetail []FormatterDetailEntityModel `json:"formatter_detail" gorm:"ForeignKey:FormatterID"`
	// TrialBalance    []TrialBalanceEntityModel    `json:"trial_balance" gorm:"foreignKey:FormatterID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type FormatterFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	FormatterFilter
}

func (FormatterEntityModel) TableName() string {
	return "m_formatter"
}

func (m *FormatterEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *FormatterEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
