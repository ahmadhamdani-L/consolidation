package model

import (
	"mcash-finance-console-core/internal/abstraction"
)

type ConsolidationBridgeEntity struct {
	ConsolidationID       int    `json:"consolidation_id" validate:"required" example:"1"`
	CompanyID             int    `json:"company_id" validate:"required" example:"1"`
	Versions              int    `json:"versions" validate:"required" example:"1"`
	ConsolidationVersions int    `json:"consolidation_versions" validate:"required" example:"1"`
	Period                string `json:"period" validate:"required" example:"2022-01-01"`
}

type ConsolidationBridgeFilter struct {
	Period      	*string `query:"period" example:"2022-01-01" filter:"DATESTRING"`
	Versions    	*int    `query:"versions" example:"1"`
	ArrVersions 	*[]int  `filter:"CUSTOM" example:"1"`
	Search          *string `query:"s" example:"Lutfi Ramadhan" filter:"CUSTOM"`
	ConsolidationID *int    `query:"id"`
	CodeConsole      *string `query:"code_console" json:"code_console"`
	Amount      *string `query:"amount" json:"amount"`
	ConsolidationVersions *int `query:"consolidation_versions" json:"consolidation_versions"`
}

type ConsolidationBridgeEntityModel struct {
	// abstraction
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`
	// entity
	ConsolidationBridgeEntity

	// relations
	Company CompanyEntityModel `json:"company" gorm:"foreignKey:CompanyID"`
	ConsolidationBridgeDetail []ConsolidationBridgeDetailEntityModel `json:"consolidation_bridge_detail" gorm:"foreignKey:ConsolidationBridgeID"`
	// Formatter          FormatterEntityModel            `json:"formatter,omitempty" gorm:"foreignKey:FormatterID"`
	// TrialBalanceDetail []TrialBalanceDetailEntityModel `json:"trial_balance_detail" gorm:"-"`
	// FormatterBridges   []FormatterBridgesEntityModel   `json:"-" gorm:"foreignKey:TrxRefId"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type ConsolidationBridgeFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	ConsolidationBridgeFilter
	CompanyCustomFilter
}

func (ConsolidationBridgeEntityModel) TableName() string {
	return "consolidation_bridge"
}
