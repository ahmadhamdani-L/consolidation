package model

import (
	"worker-consol/internal/abstraction"
)

type JpmDetailEntity struct {
	JpmID             int      `json:"jpm_id"`
	CoaCode           string   `json:"coa_code"`
	ReffNumber        *string  `json:"reff_number"`
	Description       *string  `json:"description"`
	BalanceSheetDr    *float64 `json:"balance_sheet_dr"`
	BalanceSheetCr    *float64 `json:"balance_sheet_cr"`
	IncomeStatementDr *float64 `json:"income_statement_dr"`
	IncomeStatementCr *float64 `json:"income_statement_cr"`
	Note              *string  `json:"note"`
}

type JpmDetailFilter struct {
	JpmID             *int     `json:"jpm_id"`
	CoaCode           *string  `json:"coa_code"`
	ReffNumber        *string  `json:"reff_number"`
	Description       *string  `json:"description"`
	BalanceSheetDr    *float64 `json:"balance_sheet_dr"`
	BalanceSheetCr    *float64 `json:"balance_sheet_cr"`
	IncomeStatementDr *float64 `json:"income_statement_dr"`
	IncomeStatementCr *float64 `json:"income_statement_cr"`
	Note              *string  `json:"note"`
}

type JpmDetailEntityModel struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	JpmDetailEntity

	// relations
	Jpm JpmEntityModel `json:"jpm" gorm:"-"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type JpmDetailFilterModel struct {

	// filter
	JpmDetailFilter
}

func (JpmDetailEntityModel) TableName() string {
	return "jpm_detail"
}
