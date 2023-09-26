package model

import (
	"worker/internal/abstraction"
)

type MutasiIaDetailEntity struct {
	FormatterBridgesID              int      `json:"formatter_bridges_id" validate:"required"`
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

type MutasiIaDetailFilter struct {
	FormatterBridgesID              *int     `query:"formatter_bridges_id" validate:"required"`
	Code                    *string  `query:"code"`
	Description             *string  `query:"description"`
	BeginningBalance        *float64 `query:"starting_balance"`
	AcquisitionOfSubsidiary *float64 `query:"acquisition_of_subsidiary"`
	Additions               *float64 `query:"additions"`
	Deductions              *float64 `query:"deductions"`
	Reclassification        *float64 `query:"reclassification"`
	Revaluation             *float64 `query:"revaluation"`
	EndingBalance           *float64 `query:"ending_balance"`
	Control                 *float64 `query:"control"`
	SortId                  *int     `query:"sort_id"`
}

type MutasiIaDetailEntityModel struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	MutasiIaDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	FormatterBridges FormatterBridgesEntityModel `json:"formatter_bridges" gorm:"foreignKey:FormatterBridgesID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type MutasiIaDetailFilterModel struct {

	// filter
	MutasiIaDetailFilter
}

func (MutasiIaDetailEntityModel) TableName() string {
	return "mutasi_ia_detail"
}
