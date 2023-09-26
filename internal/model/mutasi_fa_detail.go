package model

import (
	"mcash-finance-console-core/internal/abstraction"
)

type MutasiFaDetailEntity struct {
	// MutasiFaID              int      `json:"mutasi_Fa_id" validate:"required" example:"1"`
	FormatterBridgesID      int      `json:"formatter_bridges_id" validate:"required" example:"1"`
	Code                    string   `json:"code" validate:"required" example:"Tanah"`
	Description             string   `json:"description" validate:"required" example:"Tanah"`
	BeginningBalance        *float64 `json:"beginning_balance" validate:"required" example:"10000.00"`
	AcquisitionOfSubsidiary *float64 `json:"acquisition_of_subsidiary" validate:"required" example:"10000.00"`
	Additions               *float64 `json:"additions" validate:"required" example:"10000.00"`
	Deductions              *float64 `json:"deductions" validate:"required" example:"10000.00"`
	Reclassification        *float64 `json:"reclassification" validate:"required" example:"10000.00"`
	Revaluation             *float64 `json:"revaluation" validate:"required" example:"10000.00"`
	EndingBalance           *float64 `json:"ending_balance" validate:"required" example:"10000.00"`
	Control                 *float64 `json:"control" validate:"required" example:"10000.00"`
	SortID                  int      `json:"sort_id" validate:"required" example:"1"`
}

type MutasiFaDetailFilter struct {
	FormatterBridgesID      *int     `query:"formatter_bridges_id" example:"1"`
	MutasiFaID              *int     `query:"mutasi_fa_id" validate:"required" example:"1" filter:"CUSTOM"`
	Code                    *string  `query:"code" example:"Tanah"`
	Description             *string  `query:"description" example:"Tanah"`
	BeginningBalance        *float64 `query:"beginning_balance" example:"10000.00"`
	AcquisitionOfSubsidiary *float64 `query:"acquisition_of_subsidiary" example:"10000.00"`
	Additions               *float64 `query:"additions" example:"10000.00"`
	Deductions              *float64 `query:"deductions" example:"10000.00"`
	Reclassification        *float64 `query:"reclassification" example:"10000.00"`
	Revaluation             *float64 `query:"revaluation" example:"10000.00"`
	EndingBalance           *float64 `query:"ending_balance" example:"10000.00"`
	Control                 *float64 `query:"control" example:"10000.00"`
	SortID                  *int     `query:"sort_id" example:"1"`
}

type MutasiFaDetailEntityModel struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	MutasiFaDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	// MutasiFa MutasiFaEntityModel `json:"mutasi_fa" gorm:"foreignKey:MutasiFaID"`
	FormatterBridges FormatterBridgesEntityModel `json:"-" gorm:"foreignKey:FormatterBridgesID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type MutasiFaDetailFmtEntityModel struct {
	MutasiFaDetailEntityModel
	AutoSummary    *bool   `json:"auto_summary"`
	IsTotal        *bool   `json:"is_total"`
	IsControl      *bool   `json:"is_control"`
	IsLabel        *bool   `json:"is_label"`
	ControlFormula *string `json:"control_formula"`
}

type MutasiFaDetailFilterModel struct {

	// filter
	MutasiFaDetailFilter
}

func (MutasiFaDetailEntityModel) TableName() string {
	return "mutasi_fa_detail"
}
