package model

import (
	"mcash-finance-console-core/internal/abstraction"
)

type MutasiPersediaanDetailEntity struct {
	// MutasiPersediaanID int      `json:"mutasi_persediaan_id" validate:"required" example:"1"`
	FormatterBridgesID int      `json:"formatter_bridges_id" validate:"required" example:"1"`
	Code               string   `json:"code" validate:"required" example:"Saldo_Awal"`
	Description        string   `json:"description" validate:"required" example:"Saldo Awal"`
	Amount             *float64 `json:"amount" validate:"required" example:"10000.00"`
	SortID             int      `json:"sort_id" validate:"required" example:"1"`
}

type MutasiPersediaanDetailFilter struct {
	FormatterBridgesID *int     `query:"formatter_bridges_id" example:"1"`
	MutasiPersediaanID *int     `query:"mutasi_persediaan_id" validate:"required" example:"1" filter:"CUSTOM"`
	Code               *string  `query:"code" filter:"ILIKE" example:"Saldo_Awal"`
	Description        *string  `query:"description" filter:"ILIKE" example:"Saldo Awal"`
	Amount             *float64 `query:"amount" example:"10000.00"`
	SortID             *int     `query:"sort_id" example:"1"`
}

type MutasiPersediaanDetailEntityModel struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	MutasiPersediaanDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	// MutasiPersediaan MutasiPersediaanEntityModel `json:"mutasi_persediaan" gorm:"foreignKey:MutasiPersediaanID"`
	FormatterBridges FormatterBridgesEntityModel `json:"-" gorm:"foreignKey:FormatterBridgesID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type MutasiPersediaanDetailFmtEntityModel struct {
	MutasiPersediaanDetailEntityModel
	AutoSummary    *bool   `json:"auto_summary"`
	IsTotal        *bool   `json:"is_total"`
	IsControl      *bool   `json:"is_control"`
	IsLabel        *bool   `json:"is_label"`
	ControlFormula *string `json:"control_formula"`
}

type MutasiPersediaanDetailFilterModel struct {

	// filter
	MutasiPersediaanDetailFilter
}

func (MutasiPersediaanDetailEntityModel) TableName() string {
	return "mutasi_persediaan_detail"
}
