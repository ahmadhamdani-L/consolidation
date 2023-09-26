package model

import (
	"mcash-finance-console-core/internal/abstraction"
)

type MutasiRuaDetailEntity struct {
	// MutasiRuaID             int      `json:"mutasi_rua_id" validate:"required" example:"1"`
	FormatterBridgesID      int      `json:"formatter_bridges_id" validate:"required" example:"1"`
	Code                    string   `json:"code" validate:"required" example:"Bangunan"`
	Description             string   `json:"description" validate:"required" example:"Bangunan"`
	BeginningBalance        *float64 `json:"beginning_balance" validate:"required" example:"10000.00"`
	AcquisitionOfSubsidiary *float64 `json:"acquisition_of_subsidiary" validate:"required" example:"10000.00"`
	Additions               *float64 `json:"additions" validate:"required" example:"10000.00"`
	Deductions              *float64 `json:"deductions" validate:"required" example:"10000.00"`
	Reclassification        *float64 `json:"reclassification" validate:"required" example:"10000.00"`
	Remeasurement           *float64 `json:"remeasurement" validate:"required" example:"10000.00"`
	EndingBalance           *float64 `json:"ending_balance" validate:"required" example:"10000.00"`
	Control                 *float64 `json:"control" validate:"required" example:"10000.00"`
	SortID                  int      `json:"sort_id" validate:"required" example:"1"`
	IsTotal                 *bool    `json:"is_total" gorm:"-"`
	IsControl               *bool    `json:"is_control" gorm:"-"`
	ControlFormula          *string  `json:"control_formula" gorm:"-"`
}

type MutasiRuaDetailFilter struct {
	MutasiRuaID             *int     `query:"mutasi_rua_id" example:"1" filter:"CUSTOM"`
	FormatterBridgesID      *int     `query:"formatter_bridges_id" example:"1"`
	Code                    *string  `query:"code" example:"Bangunan"`
	Description             *string  `query:"description" example:"Bangunan"`
	BeginningBalance        *float64 `query:"beginning_balance" example:"10000.00"`
	AcquisitionOfSubsidiary *float64 `query:"acquisition_of_subsidRuary" example:"10000.00"`
	Additions               *float64 `query:"additions" example:"10000.00"`
	Deductions              *float64 `query:"deductions" example:"10000.00"`
	Reclassification        *float64 `query:"reclassification" example:"10000.00"`
	Remeasurement           *float64 `query:"remeasurement" example:"10000.00"`
	EndingBalance           *float64 `query:"ending_balance" example:"10000.00"`
	Control                 *float64 `query:"control" example:"10000.00"`
	SortID                  *int     `query:"sort_id" example:"1"`
}

type MutasiRuaDetailEntityModel struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	MutasiRuaDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	// MutasiRua MutasiRuaEntityModel `json:"mutasi_rua" gorm:"foreignKey:MutasiRuaID"`
	FormatterBridges FormatterBridgesEntityModel `json:"-" gorm:"foreignKey:FormatterBridgesID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type MutasiRuaDetailFmtEntityModel struct {
	MutasiRuaDetailEntityModel
	AutoSummary    *bool   `json:"auto_summary"`
	IsTotal        *bool   `json:"is_total"`
	IsControl      *bool   `json:"is_control"`
	IsLabel        *bool   `json:"is_label"`
	ControlFormula *string `json:"control_formula"`
}

type MutasiRuaDetailFilterModel struct {

	// filter
	MutasiRuaDetailFilter
}

func (MutasiRuaDetailEntityModel) TableName() string {
	return "mutasi_rua_detail"
}
