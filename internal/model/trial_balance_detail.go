package model

import (
	"mcash-finance-console-core/internal/abstraction"
)

type TrialBalanceDetailEntity struct {
	Code               string   `json:"code" validate:"required" example:"110101001"`
	AmountBeforeAje    *float64 `json:"amount_before_aje" validate:"required" example:"10000.00"`
	AmountAjeDr        *float64 `json:"amount_aje_dr" validate:"required" example:"10000.00"`
	AmountAjeCr        *float64 `json:"amount_aje_cr" validate:"required" example:"10000.00"`
	AmountAfterAje     *float64 `json:"amount_after_aje" validate:"required" example:"10000.00"`
	ReffAjeDr          *string  `json:"reff_aje_dr" validate:"required" example:"reff"`
	ReffAjeCr          *string  `json:"reff_aje_cr" validate:"required" example:"reff"`
	Description        *string  `json:"description" validate:"required" example:"Kas Kecil"`
	FormatterBridgesID int      `json:"formatter_bridges_id" validate:"required" example:"1"`
	SortID             float64  `json:"sort_id" validate:"required" example:"1"`
	// TrialBalanceID  int      `json:"trial_balance_id" validate:"required" example:"1"`
}

type TrialBalanceDetailFilter struct {
	Code               *string  `query:"code" filter:"ILIKE" example:"110101001"`
	AmountBeforeAje    *float64 `query:"amount_before_aje" example:"10000.00"`
	AmountAjeDr        *float64 `query:"amount_aje_dr" example:"10000.00"`
	AmountAjeCr        *float64 `query:"amount_aje_cr" example:"10000.00"`
	AmountAfterAje     *float64 `query:"amount_after_aje" example:"10000.00"`
	ReffAjeDr          *string  `query:"reff_aje_dr" filter:"ILIKE" example:"reff"`
	ReffAjeCr          *string  `query:"reff_aje_cr" filter:"ILIKE" example:"reff"`
	Description        *string  `query:"description" filter:"ILIKE" example:"Kas Kecil"`
	FormatterBridgesID *int     `query:"formatter_bridges_id" example:"1"`
	TrialBalanceID     *int     `query:"trial_balance_id" validate:"required" example:"1" filter:"CUSTOM"`
}

type TrialBalanceDetailEntityModel struct {
	// abstraction
	// abstraction.Entity
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	TrialBalanceDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	// TrialBalance     TrialBalanceEntityModel     `json:"trial_balance" gorm:"foreignKey:TrialBalanceID"`
	FormatterBridges FormatterBridgesEntityModel `json:"-" gorm:"foreignKey:FormatterBridgesID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type TrialBalanceDetailFmtEntityModel struct {
	TrialBalanceDetailEntityModel
	FormatterDetailID          int                                `json:"formatter_detail_id"`
	ParentID                   int                                `json:"parent_id"`
	AutoSummary                bool                               `json:"auto_summary"`
	IsTotal                    bool                               `json:"is_total"`
	IsControl                  bool                               `json:"is_control"`
	IsLabel                    bool                               `json:"is_label"`
	IsParent                   bool                               `json:"is_parent"`
	IsCoa                      bool                               `json:"-"`
	FormatterDetailCode        string                             `json:"-"`
	FormatterDetailDescription string                             `json:"-"`
	TemporaryParentID          int                                `json:"-"`
	TemporaryID                int                                `json:"-"`
	ControlFormula             string                             `json:"control_formula"`
	ShowGroupCoa               *bool                              `json:"show_group_coa"`
	Children                   []TrialBalanceDetailFmtEntityModel `json:"children" gorm:"-"`
}

type TrialBalanceDetailFilterModel struct {
	// abstraction
	// abstraction.Filter

	// filter
	TrialBalanceDetailFilter
}

func (TrialBalanceDetailEntityModel) TableName() string {
	return "trial_balance_detail"
}

// func (m *TrialBalanceDetailEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
// 	m.CreatedAt = *date.DateTodayLocal()
// 	m.CreatedBy = m.Context.Auth.ID
// 	return
// }

// func (m *TrialBalanceDetailEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
// 	m.ModifiedAt = date.DateTodayLocal()
// 	m.ModifiedBy = &m.Context.Auth.ID
// 	return
// }
