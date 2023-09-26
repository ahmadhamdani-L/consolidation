package model

import (
	"mcash-finance-console-core/internal/abstraction"
)

type ConsolidationDetailEntity struct {
	ConsolidationID				int 	`json:"consolidation_id" validate:"required"`
	Code						string 	`json:"code" validate:"required"`
	WpReff						string 	`json:"wp_reff" validate:"required"`
	Description					string 	`json:"description" validate:"required"`
	SortID						float64 `json:"sort_id" validate:"required"`
	AmountBeforeJpm				*float64 `json:"amount_before_jpm" validate:"required"`
	AmountJpmCr					*float64 `json:"amount_jpm_cr" validate:"required"`
	AmountJpmDr					*float64 `json:"amount_jpm_dr" validate:"required"`
	AmountAfterJpm				*float64 `json:"amount_after_jpm" validate:"required"`
	AmountJcteCr				*float64 `json:"amount_jcte_cr" validate:"required"`
	AmountJcteDr				*float64 `json:"amount_jcte_dr" validate:"required"`
	AmountAfterJcte				*float64 `json:"amount_after_jcte" validate:"required"`
	AmountCombineSubsidiary		*float64 `json:"amount_combine_subsidiary" validate:"required"`
	AmountJelimCr				*float64 `json:"amount_jelim_cr" validate:"required"`
	AmountJelimDr				*float64 `json:"amount_jelim_dr" validate:"required"`
	AmountConsole				*float64 `json:"amount_console" validate:"required"`
	Is_parent					bool 	 `json:"is_parent" validate:"required"`
}

type ConsolidationDetailFilter struct {
	ConsolidationID         *int      `query:"consolidation_id" json:"consolidation_id" validate:"required" filter:"CUSTOM"`
	Code                    *string   `query:"code" filter:"ILIKE" `
	WpReff                  *string    `json:"wp_reff" `
	Description             *string   `query:"description" json:"description" filter:"ILIKE" `
	AmountBeforeJpm         *float64  `query:"amount_before_jpm" json:"amount_before_jpm" `
	AmountJpmDr             *float64  `query:"amount_jpm_dr" json:"amount_jpm_cr" `
	AmountJpmCr             *float64  `query:"amount_jpm_cr" json:"amount_jpm_dr" `
	AmountAfterJpm          *float64  `query:"amount_after_jpm" json:"amount_after_jpm" `
	AmountJcteCr            *float64   `json:"amount_jcte_cr" `
	AmountJcteDr            *float64   `json:"amount_jcte_dr" `
	AmountAfterJcte         *float64   `json:"amount_after_jcte" `
	AmountCombineSubsidiary *float64   `json:"amount_combine_subsidiary" `
	AmountJelimCr           *float64   `json:"amount_jelim_cr" `
	AmountJelimDr           *float64   `json:"amount_jelim_dr" `
	AmountConsole           *float64   `json:"amount_console" `
	Is_parent				*bool 	   `json:"is_parent" `
	Parent                  *string   	`query:"parent"  example:"ASSET"`
	Source                  *[]string `query:"source"  example:"ASSET" gorm:"-"`
	Sourcea                 *string   `query:"source"  example:"ASSET"`
	ParentID                *int      `query:"parent_id"  `
	Amount                  *int      `json:"amount"  `
	CompanyID               *int      `json:"company_id"  `
	VersionConsolidation		*int   `query:"version_consolidation" json:"version_consolidation"`
	ConsolidationBridgeID		*int   `query:"consolidation_bridge_id" json:"consolidation_bridge_id"`
}

type ConsolidationDetailEntityModel struct {
	// abstraction
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	ConsolidationDetailEntity

	// relations
	Consolidation ConsolidationEntityModel `json:"-" gorm:"foreignKey:ConsolidationID"`

	ConsolidationBridge []ConsolidationBridgeEntityModel `json:"consolidation_bridge" gorm:"-"`
	

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}
type ConsolidationDetailFmtEntityModel struct {
	ConsolidationDetailEntityModel
	FormatterDetailID int                                `json:"formatter_detail_id"`
	ParentID          int                                `json:"parent_id"`
	AutoSummary       bool                               `json:"auto_summary"`
	IsTotal           bool                               `json:"is_total"`
	IsControl         bool                               `json:"is_control"`
	IsLabel           bool                               `json:"is_label"`
	ControlFormula    string                             `json:"control_formula"`
	Children          []ConsolidationDetailFmtEntityModel `json:"children" gorm:"-"`
	ConsolidationBridge   []ConsolidationBridgeEntityModel `json:"consolidation_bridge" gorm:"-"`
	ShowGroupCoa      *bool   `json:"show_group_coa"`
}
type ConsolidationDetailFilterModel struct {
	
	// filter
	ConsolidationDetailFilter
}

func (ConsolidationDetailEntityModel) TableName() string {
	return "consolidation_detail"
}
