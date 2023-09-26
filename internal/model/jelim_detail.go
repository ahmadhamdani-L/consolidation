package model

import (
	"worker/internal/abstraction"
)

type JelimDetailEntity struct {
	JelimID           int      `json:"jelim_id"`
	CoaCode           string   `json:"coa_code"`
	ReffNumber        *string  `json:"reff_number"`
	Description       *string  `json:"description"`
	BalanceSheetDr    *float64 `json:"balance_sheet_dr"`
	BalanceSheetCr    *float64 `json:"balance_sheet_cr"`
	IncomeStatementDr *float64 `json:"income_statement_dr"`
	IncomeStatementCr *float64 `json:"income_statement_cr"`
	Note              *string  `json:"note"`
}

type JelimDetailFilter struct {
	JelimID           *int     `query:"jelim_id"`
	CoaCode           *string  `query:"coa_code"`
	ReffNumber        *string  `query:"reff_number"`
	Description       *string  `query:"description"`
	BalanceSheetDr    *float64 `query:"balance_sheet_dr"`
	BalanceSheetCr    *float64 `query:"balance_sheet_cr"`
	IncomeStatementDr *float64 `query:"income_statement_dr"`
	IncomeStatementCr *float64 `query:"income_statement_cr"`
	Note              *string  `query:"note"`
}

type JelimDetailEntityModel struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	JelimDetailEntity

	// relations
	Jelim JelimEntityModel `json:"jelim" gorm:"foreignKey:JelimID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type JelimDetailFilterModel struct {
	// filter
	JelimDetailFilter
}

func (JelimDetailEntityModel) TableName() string {
	return "jelim_detail"
}
