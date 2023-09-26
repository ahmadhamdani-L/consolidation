package model

import (
	"worker/internal/abstraction"
	"worker/pkg/util/date"

	"gorm.io/gorm"
)

type FormatterBridgesEntity struct {
	TrxRefID      int `json:"trx_ref_id" `
	Source    string    `json:"source" `
	FormatterID int    `json:"formatter_id" `
}

type FormatterBridgesFilter struct {
	Source      *string `query:"source"`
	TrxRefID    *int    `query:"trx_ref_id"`
	FormatterID *int    `query:"formatter_id"`
}


type FormatterBridgesEntityModel struct {
	// abstraction
	abstraction.EntityFormatter

	// entity
	FormatterBridgesEntity

	// relations

	Formatter          FormatterEntityModel            `json:"formatter,omitempty" gorm:"foreignKey:FormatterID"`
	TrialBalanceDetail []TrialBalanceDetailEntityModel `json:"trial_balance_detail" gorm:"foreignKey:FormatterBridgesID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type FormatterBridgesFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	FormatterBridgesFilter
	// CompanyCustomFilter
}

func (FormatterBridgesEntityModel) TableName() string {
	return "formatter_bridges"
}

func (m *FormatterBridgesEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}
