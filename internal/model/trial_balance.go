package model

import (
	"worker/internal/abstraction"
	"worker/pkg/util/date"

	"gorm.io/gorm"
)

type TrialBalanceEntity struct {
	Period          string `json:"period" validate:"required" example:"2022-01-01"`
	Versions        int    `json:"versions" validate:"required" example:"1"`
	CompanyID       int    `json:"company_id" validate:"required" example:"1"`
	Status          int    `json:"status" validate:"required" example:"1"`
	ConsolidationID int    `json:"consolidation_id" example:"1"`
	// FormatterID int    `json:"formatter_id" validate:"required" example:"1"`
}

type TrialBalanceFilter struct {
	Period      *string `query:"period" example:"2022-01-01" filter:"DATESTRING"`
	Versions    *int    `query:"versions" example:"1"`
	ArrVersions *[]int  `filter:"CUSTOM" example:"1"`
	// FormatterID *int    `query:"formatter_id" example:"1"`
	Status    *int    `query:"status" example:"1"`
	ArrStatus *[]int  `filter:"CUSTOM" example:"1"`
	Search    *string `query:"s" example:"Lutfi Ramadhan" filter:"CUSTOM"`
}

type TrialBalanceEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	TrialBalanceEntity

	// relations
	Company CompanyEntityModel `json:"company" gorm:"foreignKey:CompanyID"`
	// Formatter          FormatterEntityModel            `json:"formatter,omitempty" gorm:"-"`
	TrialBalanceDetail []TrialBalanceDetailFmtEntityModel `json:"trial_balance_detail" gorm:"-"`
	FormatterBridges   []FormatterBridgesEntityModel      `json:"-" gorm:"foreignKey:TrxRefID"`
	UserRelationModel

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type TrialBalanceFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	TrialBalanceFilter
	CompanyCustomFilter
}

func (TrialBalanceEntityModel) TableName() string {
	return "trial_balance"
}

func (m *TrialBalanceEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *TrialBalanceEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
