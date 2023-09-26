package model

import (
	"worker/internal/abstraction"
)

type MutasiFaDetailEntity struct {
	// MutasiFaID              int      `json:"mutasi_Fa_id" validate:"required"`
	FormatterBridgesID      int      `json:"formatter_bridges_id" validate:"required"`
	Code                    string   `json:"code" validate:"required"`
	Description             string   `json:"description" validate:"required"`
	BeginningBalance        *float64 `json:"starting_balance" validate:"required"`
	AcquisitionOfSubsidiary *float64 `json:"acquisition_of_subsidiary" validate:"required"`
	Additions               *float64 `json:"additions" validate:"required"`
	Deductions              *float64 `json:"deductions" validate:"required"`
	Reclassification        *float64 `json:"reclassification" validate:"required"`
	Revaluation             *float64 `json:"revaluation" validate:"required"`
	EndingBalance           *float64 `json:"ending_balance" validate:"required"`
	Control                 *float64 `json:"control" validate:"required"`
	SortId                  int      `json:"sort_id" validate:"required"`
}

type MutasiFaDetailFilter struct {
	MutasiFaID              *int     `query:"mutasi_fa_id" validate:"required" filter:"CUSTOM"`
	FormatterBridgesID      *int     `query:"formatter_bridges_id" validate:"required"`
	Code                    *string  `query:"code" validate:"required"`
	Description             *string  `query:"description" validate:"required"`
	BeginningBalance        *float64 `query:"starting_balance" validate:"required"`
	AcquisitionOfSubsidiary *float64 `query:"acquisition_of_subsidiary" validate:"required"`
	Additions               *float64 `query:"additions" validate:"required"`
	Deductions              *float64 `query:"deductions" validate:"required"`
	Reclassification        *float64 `query:"reclassification" validate:"required"`
	Revaluation             *float64 `query:"revaluation" validate:"required"`
	EndingBalance           *float64 `query:"ending_balance" validate:"required"`
	Control                 *float64 `query:"control" validate:"required"`
	SortId                  *int     `query:"sort_id" validate:"required"`
}

type MutasiFaDetailEntityModel struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	MutasiFaDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
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
