package model

import (
	"time"
	"worker/internal/abstraction"
	"worker/pkg/util/date"

	"gorm.io/gorm"
)

type FormatterBridgesEntity struct {
	TrxRefID    int    `json:"trx_ref_id" validate:"required" example:"1"`
	Source      string `json:"source" validate:"required" example:"TRIAL-BALANCE"`
	FormatterID int    `json:"formatter_id" validate:"required" example:"1"`
}

type FormatterBridgesFilter struct {
	TrxRefID          *int       `query:"trx_ref_id" example:"1"`
	Source            *string    `query:"source" example:"TRIAL-BALANCE"`
	FormatterID       *int       `query:"formatter_id" example:"1"`
	CreatedAt         *time.Time `query:"created_at" filter:"DATE" example:"2022-08-17T15:04:05Z"`
	CreatedBy         *int       `query:"created_by" example:"1"`
	UserCreatedString *string    `query:"user_created" filter:"user_created" example:"Lutfi Ramadhan"`
}

type FormatterBridgesEntityModel struct {
	// abstraction
	ID                int       `json:"id" gorm:"primaryKey;autoIncrement;"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         int       `json:"created_by"`
	UserCreatedString string    `json:"user_created" gorm:"-"`

	// entity
	FormatterBridgesEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	UserCreated UserEntityModel `json:"-" gorm:"foreignKey:CreatedBy"`
	// TrialBalance TrialBalanceEntityModel `json:"-" gorm:"foreignKey:TrxRefId"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type FormatterBridgesFilterModel struct {
	// abstraction

	// filter
	FormatterBridgesFilter
}

func (FormatterBridgesEntityModel) TableName() string {
	return "formatter_bridges"
}

func (m *FormatterBridgesEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}
