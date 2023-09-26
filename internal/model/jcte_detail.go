package model

import (
	"worker/internal/abstraction"
)

type JcteDetailEntity struct {
	JcteID            int      `json:"jcte_id"`
	CoaCode           string   `json:"coa_code"`
	ReffNumber        *string  `json:"reff_number"`
	Description       *string  `json:"description"`
	BalanceSheetDr    *float64 `json:"balance_sheet_dr"`
	BalanceSheetCr    *float64 `json:"balance_sheet_cr"`
	IncomeStatementDr *float64 `json:"income_statement_dr"`
	IncomeStatementCr *float64 `json:"income_statement_cr"`
	Note              *string  `json:"note"`
}

type JcteDetailFilter struct {
	JcteID            *int     `query:"jcte_id"`
	CoaCode           *string  `query:"coa_code"`
	ReffNumber        *string  `query:"reff_number"`
	Description       *string  `query:"description"`
	BalanceSheetDr    *float64 `query:"balance_sheet_dr"`
	BalanceSheetCr    *float64 `query:"balance_sheet_cr"`
	IncomeStatementDr *float64 `query:"income_statement_dr"`
	IncomeStatementCr *float64 `query:"income_statement_cr"`
	Note              *string  `query:"note"`
}

type JcteDetailEntityModel struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	JcteDetailEntity

	// relations
	Jcte JcteEntityModel `json:"jcte" gorm:"foreignKey:JcteID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type JcteDetailFilterModel struct {

	// filter
	JcteDetailFilter
}

func (JcteDetailEntityModel) TableName() string {
	return "jcte_detail"
}
