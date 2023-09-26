package model

import (
	"worker/internal/abstraction"
)

type MutasiRuaDetailEntity struct {
	FormatterBridgesID             int      `json:"formatter_bridges_id" validate:"required"`
	Code                    string   `json:"code" validate:"required"`
	Description             string   `json:"description" validate:"required"`
	BeginningBalance        *float64 `json:"starting_balance" validate:"required"`
	AcquisitionOfSubsidiary *float64 `json:"acquisition_of_subsidiary" validate:"required"`
	Additions               *float64 `json:"additions" validate:"required"`
	Deductions              *float64 `json:"deductions" validate:"required"`
	Reclassification        *float64 `json:"reclassification" validate:"required"`
	Remeasurement           *float64 `json:"remeasurement" validate:"required"`
	EndingBalance           *float64 `json:"ending_balance" validate:"required"`
	Control                 *float64 `json:"control" validate:"required"`
	SortId                  int      `json:"sort_id" validate:"required"`
}

type MutasiRuaDetailFilter struct {
	FormatterBridgesID             *int     `query:"formatter_bridges_id" validate:"required"`
	Code                    *string  `query:"code" validate:"required"`
	Description             *string  `query:"description" validate:"required"`
	BeginningBalance        *float64 `query:"starting_balance" validate:"required"`
	AcquisitionOfSubsidiary *float64 `query:"acquisition_of_subsidRuary" validate:"required"`
	Additions               *float64 `query:"additions" validate:"required"`
	Deductions              *float64 `query:"deductions" validate:"required"`
	Reclassification        *float64 `query:"reclassification" validate:"required"`
	Remeasurement           *float64 `query:"remeasurement" validate:"required"`
	EndingBalance           *float64 `query:"ending_balance" validate:"required"`
	Control                 *float64 `query:"control" validate:"required"`
	SortId                  *int     `query:"sort_id" validate:"required"`
}

type MutasiRuaDetailEntityModel struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	MutasiRuaDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	FormatterBridges FormatterBridgesEntityModel `json:"formatter_bridges" gorm:"foreignKey:FormatterBridgesID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type MutasiRuaDetailFilterModel struct {

	// filter
	MutasiRuaDetailFilter
}

func (MutasiRuaDetailEntityModel) TableName() string {
	return "mutasi_rua_detail"
}
