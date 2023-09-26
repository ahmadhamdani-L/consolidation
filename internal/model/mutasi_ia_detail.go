package model

import (
	"worker-validation/internal/abstraction"
)

type MutasiIaDetailEntity struct {
	// MutasiIaID              int      `json:"mutasi_ia_id" validate:"required"`
	FormatterBridgesID      int      `json:"formatter_bridges_id" validate:"required" example:"1"`
	Code                    string   `json:"code" validate:"required" example:"Software"`
	Description             string   `json:"description" validate:"required" example:"Software"`
	BeginningBalance        *float64 `json:"starting_balance" validate:"required" example:"10000.00"`
	AcquisitionOfSubsidiary *float64 `json:"acquisition_of_subsidiary" validate:"required" example:"10000.00"`
	Additions               *float64 `json:"additions" validate:"required" example:"10000.00"`
	Deductions              *float64 `json:"deductions" validate:"required" example:"10000.00"`
	Reclassification        *float64 `json:"reclassification" validate:"required" example:"10000.00"`
	Revaluation             *float64 `json:"revaluation" validate:"required" example:"10000.00"`
	EndingBalance           *float64 `json:"ending_balance" validate:"required" example:"10000.00"`
	Control                 *float64 `json:"control" validate:"required" example:"10000.00"`
	SortID                  int      `json:"sort_id" validate:"required" example:"1"`
}

type MutasiIaDetailFilter struct {
	FormatterBridgesID      *int     `query:"formatter_bridges_id" example:"1"`
	MutasiIaID              *int     `query:"mutasi_ia_id" validate:"required" example:"1" filter:"CUSTOM"`
	Code                    *string  `query:"code" example:"Software"`
	Description             *string  `query:"description" example:"Software"`
	BeginningBalance        *float64 `query:"starting_balance" example:"10000.00"`
	AcquisitionOfSubsidiary *float64 `query:"acquisition_of_subsidiary" example:"10000.00"`
	Additions               *float64 `query:"additions" example:"10000.00"`
	Deductions              *float64 `query:"deductions" example:"10000.00"`
	Reclassification        *float64 `query:"reclassification" example:"10000.00"`
	Revaluation             *float64 `query:"revaluation" example:"10000.00"`
	EndingBalance           *float64 `query:"ending_balance" example:"10000.00"`
	Control                 *float64 `query:"control" example:"10000.00"`
	SortID                  *int     `query:"sort_id" example:"1"`
}

type MutasiIaDetailEntityModel struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	MutasiIaDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	// MutasiIa MutasiIaEntityModel `json:"mutasi_ia" gorm:"foreignKey:MutasiIaID"`
	FormatterBridges FormatterBridgesEntityModel `json:"-" gorm:"foreignKey:FormatterBridgesID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type MutasiIaDetailFmtEntityModel struct {
	MutasiIaDetailEntityModel
	AutoSummary    *bool   `json:"auto_summary"`
	IsTotal        *bool   `json:"is_total"`
	IsControl      *bool   `json:"is_control"`
	IsLabel        *bool   `json:"is_label"`
	ControlFormula *string `json:"control_formula"`
}

type MutasiIaDetailFilterModel struct {

	// filter
	MutasiIaDetailFilter
}

func (MutasiIaDetailEntityModel) TableName() string {
	return "mutasi_ia_detail"
}
