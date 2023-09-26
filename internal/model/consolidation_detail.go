package model

import (
	"worker-consol/internal/abstraction"
)

type ConsolidationDetailEntity struct {
	ConsolidationID         int      `json:"consolidation_id" validate:"required"`
	Code                    string   `json:"code" validate:"required"`
	WpReff                  string   `json:"wp_reff" validate:"required"`
	Description             string   `json:"description" validate:"required"`
	SortID                  float64  `json:"sort_id" validate:"required"`
	AmountBeforeJpm         *float64 `json:"amount_before_jpm" validate:"required"`
	AmountJpmDr             *float64 `json:"amount_jpm_dr" validate:"required"`
	AmountJpmCr             *float64 `json:"amount_jpm_cr" validate:"required"`
	AmountAfterJpm          *float64 `json:"amount_after_jpm" validate:"required"`
	AmountJcteDr            *float64 `json:"amount_jcte_dr" validate:"required"`
	AmountJcteCr            *float64 `json:"amount_jcte_cr" validate:"required"`
	AmountAfterJcte         *float64 `json:"amount_after_jcte" validate:"required"`
	AmountCombineSubsidiary *float64 `json:"amount_combine_subsidiary" validate:"required"`
	AmountJelimDr           *float64 `json:"amount_jelim_dr" validate:"required"`
	AmountJelimCr           *float64 `json:"amount_jelim_cr" validate:"required"`
	AmountConsole           *float64 `json:"amount_console" validate:"required"`
}

type ConsolidationDetailFilter struct {
	ConsolidationID *int    `query:"consolidation_id"`
	Code            *string `query:"code"`
}

type ConsolidationDetailEntityModel struct {
	// abstraction
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	ConsolidationDetailEntity

	// relations
	Consolidation ConsolidationEntityModel `json:"consolidation" gorm:"foreignKey:ConsolidationID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type ConsolidationDetailFilterModel struct {

	// filter
	ConsolidationDetailFilter
	CompanyCustomFilter
}

func (ConsolidationDetailEntityModel) TableName() string {
	return "consolidation_detail"
}
