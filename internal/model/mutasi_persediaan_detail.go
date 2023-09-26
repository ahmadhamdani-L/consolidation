package model

import (
	"worker/internal/abstraction"
)

type MutasiPersediaanDetailEntity struct {
	// MutasiPersediaanID int      `json:"mutasi_persediaan_id" validate:"required"`
	FormatterBridgesID int      `json:"formatter_bridges_id" validate:"required"`
	Code               string   `json:"code" validate:"required"`
	Description        string   `json:"description" validate:"required"`
	Amount             *float64 `json:"amount" validate:"required"`
	SortID             int      `json:"sort_id" validate:"required"`
}

type MutasiPersediaanDetailFilter struct {
	MutasiPersediaanID *int     `query:"mutasi_persediaan_id" validate:"required" filter:"CUSTOM"`
	FormatterBridgesID *int     `query:"formatter_bridges_id" validate:"required"`
	Code               *string  `query:"code" filter:"ILIKE"`
	Description        *string  `query:"description" filter:"ILIKE"`
	Amount             *float64 `query:"amount"`
	SortID             *int     `query:"sort_id"`
}

type MutasiPersediaanDetailEntityModel struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	MutasiPersediaanDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
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
